package cloudwatchlogs

import (
	"testing"

	"github.com/aws/amazon-cloudwatch-agent/plugins/outputs/cloudwatchlogs/internal/pusher"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/require"

	"github.com/aws/amazon-cloudwatch-agent/logs"
	"github.com/aws/amazon-cloudwatch-agent/tool/util"
)

func TestCreateDestination(t *testing.T) {
	testCases := map[string]struct {
		cfgLogGroup               string
		cfgLogStream              string
		cfgLogRetention           int
		cfgLogClass               string
		cfgTailerSrc              logs.LogSrc
		expectedLogGroup          string
		expectedLogStream         string
		expectedLogGroupRetention int
		expectedLogClass          string
		expectedTailerSrc         logs.LogSrc
	}{
		"WithTomlGroupStream": {
			cfgLogGroup:               "",
			cfgLogStream:              "",
			cfgLogRetention:           -1,
			cfgLogClass:               "",
			cfgTailerSrc:              nil,
			expectedLogGroup:          "G1",
			expectedLogStream:         "S1",
			expectedLogGroupRetention: -1,
			expectedTailerSrc:         nil,
		},
		"WithOverrideGroupStreamStandardLogGroup": {
			cfgLogGroup:               "",
			cfgLogStream:              "",
			cfgLogRetention:           -1,
			cfgLogClass:               util.StandardLogGroupClass,
			cfgTailerSrc:              nil,
			expectedLogGroup:          "G1",
			expectedLogStream:         "S1",
			expectedLogGroupRetention: -1,
			expectedLogClass:          util.StandardLogGroupClass,
			expectedTailerSrc:         nil,
		},
		"WithOverrideGroupStreamInfrequentLogGroup": {
			cfgLogGroup:               "Group5",
			cfgLogStream:              "Stream5",
			cfgLogRetention:           -1,
			cfgLogClass:               util.InfrequentAccessLogGroupClass,
			cfgTailerSrc:              nil,
			expectedLogGroup:          "Group5",
			expectedLogStream:         "Stream5",
			expectedLogGroupRetention: -1,
			expectedLogClass:          util.InfrequentAccessLogGroupClass,
			expectedTailerSrc:         nil,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			cache, err := lru.New[pusher.Target, *cwDest](1000)
			require.NoError(t, err)

			c := &CloudWatchLogs{
				Log:            testutil.Logger{Name: "test"},
				LogGroupName:   "G1",
				LogStreamName:  "S1",
				AccessKey:      "access_key",
				SecretKey:      "secret_key",
				pusherStopChan: make(chan struct{}),
				cwDests:        cache,
			}
			dest := c.CreateDest(testCase.cfgLogGroup, testCase.cfgLogStream, testCase.cfgLogRetention, testCase.cfgLogClass, testCase.cfgTailerSrc).(*cwDest)
			require.Equal(t, testCase.expectedLogGroup, dest.pusher.Group)
			require.Equal(t, testCase.expectedLogStream, dest.pusher.Stream)
			require.Equal(t, testCase.expectedLogGroupRetention, dest.pusher.Retention)
			require.Equal(t, testCase.expectedLogClass, dest.pusher.Class)
			require.Equal(t, testCase.expectedTailerSrc, dest.pusher.EntityProvider)
		})
	}
}

func TestDuplicateDestination(t *testing.T) {
	cache, err := lru.New[pusher.Target, *cwDest](1000)
	require.NoError(t, err)

	c := &CloudWatchLogs{
		Log:            testutil.Logger{Name: "test"},
		AccessKey:      "access_key",
		SecretKey:      "secret_key",
		cwDests:        cache,
		pusherStopChan: make(chan struct{}),
	}
	// Given the same log group, log stream, same retention, and logClass
	d1 := c.CreateDest("FILENAME", "", -1, util.InfrequentAccessLogGroupClass, nil)
	d2 := c.CreateDest("FILENAME", "", -1, util.InfrequentAccessLogGroupClass, nil)

	// Then the destination for cloudwatchlogs endpoint would be the same
	require.Equal(t, d1, d2)
}

func TestLRUEviction(t *testing.T) {
	cache, err := lru.New[pusher.Target, *cwDest](2) // Small size for testing
	require.NoError(t, err)

	c := &CloudWatchLogs{
		Log:            testutil.Logger{Name: "test"},
		AccessKey:      "access_key",
		SecretKey:      "secret_key",
		cwDests:        cache,
		pusherStopChan: make(chan struct{}),
	}

	// Create 3 destinations
	c.CreateDest("group1", "stream1", -1, "", nil)
	c.CreateDest("group2", "stream2", -1, "", nil)
	c.CreateDest("group3", "stream3", -1, "", nil)

	// Verify that the cache maintains its maximum size of 2
	require.Equal(t, 2, c.cwDests.Len(), "Cache should maintain max size of 2")
}
