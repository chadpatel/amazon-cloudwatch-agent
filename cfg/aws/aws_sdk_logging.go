// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package aws

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

const (
	LogDebug = aws.LogRetries | aws.LogRequest
	LogDebugWithSigning = LogDebug | aws.LogSigning
	LogDebugWithHTTPBody = LogDebug | aws.LogRequestWithBody | aws.LogResponseWithBody
	LogDebugRequestRetries = LogDebug | aws.LogRetries
	LogDebugWithEventStreamBody = LogDebug | aws.LogRequestEventMessage | aws.LogResponseEventMessage
)

// Hard coded strings that match actual variable names in the AWS SDK v2.
// Update this map if/when AWS SDK adds more levels (unlikely).
var stringToLevelMap = map[string]aws.ClientLogMode{
	"LogRequest":             aws.LogRequest,
	"LogRequestWithBody":     aws.LogRequestWithBody,
	"LogResponse":            aws.LogResponse,
	"LogResponseWithBody":    aws.LogResponseWithBody,
	"LogRetries":             aws.LogRetries,
	"LogSigning":             aws.LogSigning,
	"LogDeprecatedUsage":     aws.LogDeprecatedUsage,
	"LogRequestEventMessage": aws.LogRequestEventMessage,
	"LogResponseEventMessage": aws.LogResponseEventMessage,
	// backwards compatibility with legacy AWS Go SDK v1
	"LogDebug": LogDebug,
	"LogDebugWithSigning": LogDebugWithSigning,
	"LogDebugWithHTTPBody": LogDebugWithHTTPBody,
	"LogDebugRequestRetries": LogDebugRequestRetries,
	"LogDebugWithEventStreamBody": LogDebugWithEventStreamBody,
}

var sdkLogLevel = aws.ClientLogMode(0) // Default to no logging

func setupAWSConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithClientLogMode(sdkLogLevel),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return cfg
}

// SetSDKLogLevel sets the global log level which will be used in all AWS SDK calls.
// The levels are a bit field that is OR'd together.
// So the user can specify multiple levels and we OR them together.
// Example: "aws_sdk_log_level": "LogDebugWithSigning | LogDebugWithRequestErrors".
// JSON string value must contain the levels seperated by "|" and optionally whitespace.
func SetSDKLogLevel(sdkLogLevelString string) {
	var temp aws.ClientLogMode = 0 // all flags disabled

	levels := strings.Split(sdkLogLevelString, "|")
	for _, v := range levels {
		trimmed := strings.TrimSpace(v)
		// If v not in map, then OR with 0 is harmless.
		temp |= stringToLevelMap[trimmed]
	}

	sdkLogLevel = temp
}

// SDKLogLevel returns the single global value so it can be used in all
// AWS SDK calls scattered throughout the Agent.
func SDKLogLevel() *aws.ClientLogMode {
	return &sdkLogLevel
}

// SDKLogger implements the aws.Logger interface.
type SDKLogger struct {
}

// Log is the only method in the aws.Logger interface.
func (SDKLogger) Log(args ...interface{}) {
	// Always use info logging level.
	tempSlice := append([]interface{}{"I!"}, args...)
	log.Println(tempSlice...)
}
