package mock

import (
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/firehose/firehoseiface"
)

var _ firehoseiface.FirehoseAPI = (*FirehoseMock)(nil)

type PutRecordBatchFunc func(*firehose.PutRecordBatchInput) (*firehose.PutRecordBatchOutput, error)

type FirehoseMock struct {
	firehoseiface.FirehoseAPI
	PutRecordBatchFunc        PutRecordBatchFunc
	PutRecordBatchFuncInvoked bool
}

func (f *FirehoseMock) PutRecordBatch(input *firehose.PutRecordBatchInput) (*firehose.PutRecordBatchOutput, error) {
	f.PutRecordBatchFuncInvoked = true
	return f.PutRecordBatchFunc(input)
}
