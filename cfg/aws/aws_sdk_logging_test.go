// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestSetSDKLogLevel(t *testing.T) {
	cases := []struct {
		sdkLogLevelString string
		expectedVal       aws.ClientLogMode
	}{
		// sdkLogLevelString does not match
		{"FOO", 0},
		// Wrong case.
		{"logDEBUG", 0},
		// Extra char.
		{"LogDebug1", 0},
		// Single match.
		{"LogRequest", aws.LogRequest},
		{"LogRequestWithBody", aws.LogRequestWithBody},
		{"LogResponse", aws.LogResponse},
		{"LogResponseWithBody", aws.LogResponseWithBody},
		{"LogRetries", aws.LogRetries},
		{"LogSigning", aws.LogSigning},
		{"LogDeprecatedUsage", aws.LogDeprecatedUsage},
		{"LogRequestEventMessage", aws.LogRequestEventMessage},
		{"LogResponseEventMessage", aws.LogResponseEventMessage},
		{"LogDebug", LogDebug},
		{"LogDebugWithSigning", LogDebugWithSigning},
		{"LogDebugWithHTTPBody", LogDebugWithHTTPBody},
		{"LogDebugRequestRetries", LogDebugRequestRetries},
		{"LogDebugWithEventStreamBody", LogDebugWithEventStreamBody},
		// Extra space around is allowed.
		{"   LogDebug  ", LogDebug},
		// Multiple matches.
		{"LogDebugWithEventStreamBody|LogDebugWithHTTPBody",
			LogDebugWithEventStreamBody | LogDebugWithHTTPBody},
		{"  LogDebugWithHTTPBody  |  LogDebugWithEventStreamBody  ",
			LogDebugWithEventStreamBody | LogDebugWithHTTPBody},
		{"LogDebugRequestRetries|LogDebugWithEventStreamBody",
			LogDebugWithEventStreamBody | LogDebugRequestRetries},
		//{"LogDebugRequestRetries|LogDebugWithRequestErrors",
		//	LogDebugRequestRetries | aws.LogDebugWithRequestErrors},
		//{"LogDebugRequestRetries|LogDebugWithRequestErrors|LogDebugWithEventStreamBody",
		//	LogDebugRequestRetries | aws.LogDebugWithRequestErrors | LogDebugWithEventStreamBody},
	}

	for _, tc := range cases {
		SetSDKLogLevel(tc.sdkLogLevelString)
		// check the internal var
		if sdkLogLevel != tc.expectedVal {
			t.Errorf("input: %v, actual: %v", tc, sdkLogLevel)
		}
	}
}