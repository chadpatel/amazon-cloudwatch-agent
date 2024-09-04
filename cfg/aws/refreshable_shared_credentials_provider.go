package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type RefreshableSharedCredentialsProvider struct {
	aws.CredentialsProvider
	ExpiryWindow time.Duration
	ExpirationProvider ExpirationProvider
}

type ExpirationProvider interface {
	SetExpiration(time.Time, time.Duration)
}

func (p *RefreshableSharedCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}

	p.ExpirationProvider.SetExpiration(time.Now().Add(p.ExpiryWindow), 0)
	return creds, nil
}