package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type ChainCredentialsProvider struct {
	providers []aws.CredentialsProvider
}

func NewChainCredentialsProvider(providers ...aws.CredentialsProvider) *ChainCredentialsProvider {
	return &ChainCredentialsProvider{
		providers: providers,
	}
}

func (c *ChainCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	var err error
	for _, provider := range c.providers {
		creds, err := provider.Retrieve(ctx)
		if err == nil {
			return creds, nil
		}
	}
	return aws.Credentials{}, err
}

func (c *ChainCredentialsProvider) IsExpired(ctx context.Context) bool {
	for _, provider := range c.providers {
		if creds, err := provider.Retrieve(ctx); err == nil {
			return creds.Expired()
		}
	}
	return true
}
