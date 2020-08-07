package logging

import (
	"context"
	"encoding/json"
	"math"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

const (
	kinesisMaxRetries = 8

	// See
	// https://docs.aws.amazon.com/sdk-for-go/api/service/kinesis/#Kinesis.PutRecords
	// for documentation on limits.
	kinesisMaxRecordsInBatch = 500
	kinesisMaxSizeOfRecord   = 1000 * 1000     // 1,000 KB
	kinesisMaxSizeOfBatch    = 5 * 1000 * 1000 // 5 MB
)

type kinesisLogWriter struct {
	client kinesisiface.KinesisAPI
	stream string
	logger log.Logger
	rand   *rand.Rand
}

func NewKinesisLogWriter(region, id, secret, stsAssumeRoleArn, stream string, logger log.Logger) (*kinesisLogWriter, error) {
	conf := &aws.Config{
		Region: &region,
	}

	// Only provide static credentials if we have them
	// otherwise use the default credentials provider chain
	if id != "" && secret != "" {
		conf.Credentials = credentials.NewStaticCredentials(id, secret, "")
	}

	sess, err := session.NewSession(conf)
	if err != nil {
		return nil, errors.Wrap(err, "create Kinesis client")
	}

	if stsAssumeRoleArn != "" {
		creds := stscreds.NewCredentials(sess, stsAssumeRoleArn)
		conf.Credentials = creds

		sess, err = session.NewSession(conf)

		if err != nil {
			return nil, errors.Wrap(err, "create Kinesis client")
		}
	}
	client := kinesis.New(sess)

	// This will be used to generate random partition keys to balance
	// records across Kinesis shards.
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	f := &kinesisLogWriter{
		client: client,
		stream: stream,
		logger: logger,
		rand:   rand,
	}
	if err := f.validateStream(); err != nil {
		return nil, errors.Wrap(err, "create Kinesis writer")
	}
	return f, nil
}

func (f *kinesisLogWriter) validateStream() error {
	out, err := f.client.DescribeStream(
		&kinesis.DescribeStreamInput{
			StreamName: &f.stream,
		},
	)
	if err != nil {
		return errors.Wrapf(err, "describe stream %s", f.stream)
	}

	if (*(*out.StreamDescription).StreamStatus) != kinesis.StreamStatusActive {
		return errors.Errorf("stream %s not active", f.stream)
	}

	return nil
}

func (f *kinesisLogWriter) Write(ctx context.Context, logs []json.RawMessage) error {
	var records []*kinesis.PutRecordsRequestEntry
	totalBytes := 0
	for _, log := range logs {
		// We don't really have a good option for what to do with logs
		// that are too big for Kinesis. This behavior is consistent
		// with osquery's behavior in the Kinesis logger plugin, and
		// the beginning bytes of the log should help the Fleet admin
		// diagnose the query generating huge results.
		if len(log) > kinesisMaxSizeOfRecord {
			level.Info(f.logger).Log(
				"msg", "dropping log over 1MB Kinesis limit",
				"size", len(log),
				"log", string(log[:100])+"...",
			)
			continue
		}

		partitionKey := string(f.rand.Intn(256))

		// If adding this log will exceed the limit on number of
		// records in the batch, or the limit on total size of the
		// records in the batch, we need to push this batch before
		// adding any more.
		if len(records) >= kinesisMaxRecordsInBatch ||
			totalBytes+len(log)+len(partitionKey) > kinesisMaxSizeOfBatch {
			if err := f.putRecords(0, records); err != nil {
				return errors.Wrap(err, "put records")
			}
			totalBytes = 0
			records = nil
		}

		records = append(records, &kinesis.PutRecordsRequestEntry{Data: []byte(log), PartitionKey: aws.String(partitionKey)})
		totalBytes += len(log) + len(partitionKey)
	}

	// Push the final batch
	if len(records) > 0 {
		if err := f.putRecords(0, records); err != nil {
			return errors.Wrap(err, "put records")
		}
	}

	return nil
}

func (f *kinesisLogWriter) putRecords(try int, records []*kinesis.PutRecordsRequestEntry) error {
	if try > 0 {
		time.Sleep(100 * time.Millisecond * time.Duration(math.Pow(2.0, float64(try))))
	}
	input := &kinesis.PutRecordsInput{
		StreamName: &f.stream,
		Records:    records,
	}

	output, err := f.client.PutRecords(input)
	if err != nil {
		if try < kinesisMaxRetries {
			// Retry with backoff
			return f.putRecords(try+1, records)
		}

		// Not retryable or retries expired
		return err
	}

	// Check errors on individual records
	if output.FailedRecordCount != nil && *output.FailedRecordCount > 0 {
		if try >= kinesisMaxRetries {
			// Retrieve first error message to provide to user.
			// There could be up to kinesisMaxRecordsInBatch
			// errors here and we don't want to flood that.
			var errMsg string
			for _, record := range output.Records {
				if record.ErrorCode != nil && record.ErrorMessage != nil {
					errMsg = *record.ErrorMessage
					break
				}
			}

			return errors.Errorf(
				"failed to put %d records, retries exhausted. First error: %s",
				output.FailedRecordCount, errMsg,
			)
		}

		var failedRecords []*kinesis.PutRecordsRequestEntry
		// Collect failed records for retry
		for i, record := range output.Records {
			if record.ErrorCode != nil {
				failedRecords = append(failedRecords, records[i])
			}
		}

		return f.putRecords(try+1, failedRecords)
	}

	return nil
}
