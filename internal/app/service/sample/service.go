package sample

import "context"

type SampleService struct {
	sampleServiceDynamodb SampleServiceDynamodbI
}

// SampleServiceParam 是用於建構 SampleService 的參數結構體
type SampleServiceParam struct {
	SampleServiceDynamodb SampleServiceDynamodbI
}

// NewSampleService 建立並回傳一個新的 SampleService 實例
// param: SampleService 的建構參數，包含一個 SampleServiceDynamodb 實例
func NewSampleService(_ context.Context, param SampleServiceParam) *SampleService {
	return &SampleService{
		sampleServiceDynamodb: param.SampleServiceDynamodb,
	}
}
