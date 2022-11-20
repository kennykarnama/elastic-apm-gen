package service

import (
	"context"

	"github.com/kennykarnama/elastic-apm-gen/fakeclient"
)

type Service struct {
	fakeClient fakeclient.Client
}

func NewService(fakeClient fakeclient.Client) *Service {
	return &Service{
		fakeClient: fakeClient,
	}
}

func (s *Service) Process(ctx context.Context) (string, error) {
	// apm:startSpan process request
	data, err := s.fakeClient.GetData(ctx)
	if err != nil {
		return "", err
	}
	return data.ID, nil
}
