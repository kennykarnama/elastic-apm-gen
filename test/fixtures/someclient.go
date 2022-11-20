package fixtures

import "context"

type SomeClient struct{}

func (sc *SomeClient) DoSomeTask(ctx context.Context) error {
	// apm:startSpan DoSomeTask request
	return nil
}
