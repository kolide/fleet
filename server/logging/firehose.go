package logging

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/firehose/firehoseiface"
	"github.com/pkg/errors"
)

const (
	firehoseMaxRetries = 3

	// See
	// https://docs.aws.amazon.com/sdk-for-go/api/service/firehose/#Firehose.PutRecordBatch
	// for documentation on limits.
	firehoseMaxRecordsInBatch = 500
	firehoseMaxSizeOfRecord   = 1000 * 1000     // 1,000 KB
	firehoseMaxSizeOfBatch    = 4 * 1000 * 1000 // 4 MB
)

type firehoseLogWriter struct {
	client firehoseiface.FirehoseAPI
	stream string
}

func NewFirehoseLogWriter(region, id, secret, stream string) (*firehoseLogWriter, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(id, secret, ""),
		Region:      &region,
	})
	if err != nil {
		return nil, err
	}
	client := firehose.New(sess)
	return &firehoseLogWriter{client, stream}, nil
}

func (f *firehoseLogWriter) Write(logs []json.RawMessage) error {
	var records []*firehose.Record
	totalBytes := 0
	for _, log := range logs {
		if len(log) > firehoseMaxSizeOfRecord {
			// TODO zwass: better logging
			fmt.Println("too big!")
			continue
		}

		if len(records) >= firehoseMaxRecordsInBatch ||
			totalBytes+len(log) >= firehoseMaxSizeOfBatch {
			if err := f.putRecordBatch(0, records); err != nil {
				return errors.Wrap(err, "put records")
			}
			totalBytes = 0
			records = nil
		}
		records = append(records, &firehose.Record{Data: []byte(log)})
		totalBytes += len(log)
	}
	if len(records) > 0 {
		if err := f.putRecordBatch(0, records); err != nil {
			return errors.Wrap(err, "put records")
		}
	}

	return nil
}

func (f *firehoseLogWriter) putRecordBatch(try int, records []*firehose.Record) error {
	if try > 0 {
		time.Sleep(100 * time.Millisecond * time.Duration(math.Pow(2.0, float64(try))))
	}
	input := &firehose.PutRecordBatchInput{
		DeliveryStreamName: &f.stream,
		Records:            records,
	}

	output, err := f.client.PutRecordBatch(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == firehose.ErrCodeServiceUnavailableException && try < firehoseMaxRetries {
				// Retry with backoff
				return f.putRecordBatch(try+1, records)
			}
		}

		// Not retryable or retries expired
		return err
	}

	// Check errors on individual records
	if output.FailedPutCount != nil && *output.FailedPutCount > 0 {
		if try >= firehoseMaxRetries {
			return errors.Errorf(
				"failed to put %d records, retries exhausted",
				output.FailedPutCount,
			)
		}

		var failedRecords []*firehose.Record
		// Collect failed records for retry
		for i, record := range output.RequestResponses {
			if record.ErrorCode != nil {
				failedRecords = append(failedRecords, records[i])
			}
		}

		return f.putRecordBatch(try+1, failedRecords)
	}

	return nil
}
