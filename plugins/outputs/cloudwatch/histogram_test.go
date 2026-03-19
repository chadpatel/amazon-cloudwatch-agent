// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package cloudwatch

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func newHistogramDP(bounds []float64, counts []uint64) pmetric.HistogramDataPoint {
	dp := pmetric.NewHistogramDataPoint()
	dp.ExplicitBounds().FromRaw(bounds)
	dp.BucketCounts().FromRaw(counts)
	return dp
}

// --- CheckValidity tests ---

func TestCheckValidity_Valid(t *testing.T) {
	dp := newHistogramDP([]float64{5, 10}, []uint64{600, 10, 387})
	dp.SetMin(1)
	dp.SetMax(20)
	dp.SetSum(5000)
	dp.SetCount(997)
	assert.NoError(t, CheckValidity(dp))
}

func TestCheckValidity_MismatchedBucketBoundsLength(t *testing.T) {
	dp := newHistogramDP([]float64{5, 10}, []uint64{600, 10})
	assert.ErrorContains(t, CheckValidity(dp), "bucket counts length 2 != bounds length 2 + 1")
}

func TestCheckValidity_NaNBound(t *testing.T) {
	dp := newHistogramDP([]float64{math.NaN(), 10}, []uint64{1, 2, 3})
	assert.ErrorContains(t, CheckValidity(dp), "NaN or Inf")
}

func TestCheckValidity_NonAscendingBounds(t *testing.T) {
	dp := newHistogramDP([]float64{10, 5}, []uint64{1, 2, 3})
	assert.ErrorContains(t, CheckValidity(dp), "ascending order")
}

func TestCheckValidity_MinGreaterThanMax(t *testing.T) {
	dp := newHistogramDP([]float64{5}, []uint64{1, 2})
	dp.SetMin(20)
	dp.SetMax(10)
	assert.ErrorContains(t, CheckValidity(dp), "min > max")
}

func TestCheckValidity_NaNMin(t *testing.T) {
	dp := newHistogramDP([]float64{5}, []uint64{1, 2})
	dp.SetMin(math.NaN())
	assert.ErrorContains(t, CheckValidity(dp), "min is NaN")
}

func TestCheckValidity_InfSum(t *testing.T) {
	dp := newHistogramDP([]float64{5}, []uint64{1, 2})
	dp.SetSum(math.Inf(1))
	assert.ErrorContains(t, CheckValidity(dp), "sum is NaN or Inf")
}

// --- ConvertOTelToCloudWatch tests ---

func TestConvertOTelToCloudWatch_SingleBucketNoBounds(t *testing.T) {
	dp := newHistogramDP(nil, []uint64{42})
	dp.SetMin(3)
	dp.SetMax(7)
	dp.SetSum(210)
	dp.SetCount(42)

	h := ConvertOTelToCloudWatch(dp)
	assert.Equal(t, float64(42), h.SampleCount())
	assert.Equal(t, float64(210), h.Sum())
	assert.Equal(t, float64(3), h.Minimum())
	assert.Equal(t, float64(7), h.Maximum())

	values, counts := h.ValuesAndCounts()
	require.Len(t, values, 1)
	require.Len(t, counts, 1)
	assert.Equal(t, 5.0, values[0]) // (3+7)/2
	assert.Equal(t, float64(42), counts[0])
}

func TestConvertOTelToCloudWatch_Standard3Bucket(t *testing.T) {
	dp := newHistogramDP([]float64{5, 10}, []uint64{600, 10, 387})
	dp.SetMin(0)
	dp.SetMax(100)
	dp.SetSum(50000)
	dp.SetCount(997)

	h := ConvertOTelToCloudWatch(dp)
	assert.Equal(t, float64(997), h.SampleCount())
	assert.Equal(t, float64(50000), h.Sum())
	assert.Equal(t, float64(0), h.Minimum())
	assert.Equal(t, float64(100), h.Maximum())

	values, counts := h.ValuesAndCounts()
	require.NotEmpty(t, values)
	require.Equal(t, len(values), len(counts))

	// Total counts across all inner buckets should equal total sample count.
	totalCount := 0.0
	for _, c := range counts {
		totalCount += c
	}
	assert.Equal(t, float64(997), totalCount)

	// Last value should be snapped to maximum.
	assert.Equal(t, float64(100), values[len(values)-1])
}

func TestConvertOTelToCloudWatch_EmptyHistogram(t *testing.T) {
	dp := newHistogramDP([]float64{5, 10}, []uint64{0, 0, 0})
	dp.SetMin(0)
	dp.SetMax(10)
	dp.SetSum(0)
	dp.SetCount(0)

	h := ConvertOTelToCloudWatch(dp)
	assert.Equal(t, float64(0), h.SampleCount())
	values, counts := h.ValuesAndCounts()
	assert.Empty(t, values)
	assert.Empty(t, counts)
}

func TestConvertOTelToCloudWatch_WithoutMinMax(t *testing.T) {
	dp := newHistogramDP([]float64{5, 10}, []uint64{600, 10, 387})
	dp.SetSum(50000)
	dp.SetCount(997)
	// No SetMin/SetMax — min/max should be inferred.

	h := ConvertOTelToCloudWatch(dp)
	assert.Equal(t, float64(997), h.SampleCount())

	values, counts := h.ValuesAndCounts()
	require.NotEmpty(t, values)

	totalCount := 0.0
	for _, c := range counts {
		totalCount += c
	}
	assert.Equal(t, float64(997), totalCount)
}

// --- AddHistograms tests ---

func TestAddHistograms_MatchingBoundariesMerge(t *testing.T) {
	state := newHistogramDP([]float64{5, 10}, []uint64{100, 20, 80})
	state.SetMin(1)
	state.SetMax(15)
	state.SetSum(1000)
	state.SetCount(200)

	dp := newHistogramDP([]float64{5, 10}, []uint64{50, 10, 40})
	dp.SetMin(2)
	dp.SetMax(20)
	dp.SetSum(500)
	dp.SetCount(100)

	AddHistograms(state, dp)

	// Bucket counts merged element-wise.
	assert.Equal(t, uint64(150), state.BucketCounts().At(0))
	assert.Equal(t, uint64(30), state.BucketCounts().At(1))
	assert.Equal(t, uint64(120), state.BucketCounts().At(2))
	// Count and sum accumulated.
	assert.Equal(t, uint64(300), state.Count())
	assert.Equal(t, float64(1500), state.Sum())
	// Min takes lower, max takes higher.
	assert.Equal(t, float64(1), state.Min())
	assert.Equal(t, float64(20), state.Max())
}

func TestAddHistograms_DifferentBoundariesReset(t *testing.T) {
	state := newHistogramDP([]float64{5, 10}, []uint64{100, 20, 80})
	state.SetCount(200)
	state.SetSum(1000)

	dp := newHistogramDP([]float64{3, 7, 15}, []uint64{10, 20, 30, 40})
	dp.SetCount(100)
	dp.SetSum(500)

	AddHistograms(state, dp)

	// State replaced by dp since boundaries differ.
	assert.Equal(t, 3, state.ExplicitBounds().Len())
	assert.Equal(t, uint64(100), state.Count())
	assert.Equal(t, float64(500), state.Sum())
}

func TestAddHistograms_SumCountAccumulation(t *testing.T) {
	state := newHistogramDP([]float64{10}, []uint64{50, 50})
	state.SetSum(500)
	state.SetCount(100)

	dp := newHistogramDP([]float64{10}, []uint64{25, 75})
	dp.SetSum(300)
	dp.SetCount(100)

	AddHistograms(state, dp)

	assert.Equal(t, uint64(200), state.Count())
	assert.Equal(t, float64(800), state.Sum())
	assert.Equal(t, uint64(75), state.BucketCounts().At(0))
	assert.Equal(t, uint64(125), state.BucketCounts().At(1))
}
