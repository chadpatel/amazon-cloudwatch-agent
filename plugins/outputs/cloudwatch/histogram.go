// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package cloudwatch

import (
	"errors"
	"fmt"
	"math"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

// CloudWatchHistogram holds converted histogram data ready for PutMetricData.
type CloudWatchHistogram struct {
	maximum     float64
	minimum     float64
	sampleCount float64
	sum         float64
	values      []float64
	counts      []float64
}

func (h *CloudWatchHistogram) ValuesAndCounts() ([]float64, []float64) {
	return h.values, h.counts
}

func (h *CloudWatchHistogram) Maximum() float64     { return h.maximum }
func (h *CloudWatchHistogram) Minimum() float64     { return h.minimum }
func (h *CloudWatchHistogram) Sum() float64         { return h.sum }
func (h *CloudWatchHistogram) SampleCount() float64 { return h.sampleCount }

// CheckValidity validates an OTel histogram data point.
func CheckValidity(dp pmetric.HistogramDataPoint) error {
	bounds := dp.ExplicitBounds()
	bucketCounts := dp.BucketCounts()
	if bucketCounts.Len() != bounds.Len()+1 {
		return fmt.Errorf("bucket counts length %d != bounds length %d + 1", bucketCounts.Len(), bounds.Len())
	}
	if dp.HasMin() && dp.HasMax() && dp.Min() > dp.Max() {
		return errors.New("min > max")
	}
	for i := 0; i < bounds.Len(); i++ {
		v := bounds.At(i)
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return fmt.Errorf("bound at index %d is NaN or Inf", i)
		}
		if i > 0 && v <= bounds.At(i-1) {
			return fmt.Errorf("bounds not in ascending order at index %d", i)
		}
	}
	if dp.HasMin() && (math.IsNaN(dp.Min()) || math.IsInf(dp.Min(), 0)) {
		return errors.New("min is NaN or Inf")
	}
	if dp.HasMax() && (math.IsNaN(dp.Max()) || math.IsInf(dp.Max(), 0)) {
		return errors.New("max is NaN or Inf")
	}
	if dp.HasSum() && (math.IsNaN(dp.Sum()) || math.IsInf(dp.Sum(), 0)) {
		return errors.New("sum is NaN or Inf")
	}
	return nil
}

// AddHistograms merges dp into state. If boundaries differ, state is replaced by dp.
// If boundaries match, bucket counts are added element-wise and min/max/sum/count are merged.
func AddHistograms(state, dp pmetric.HistogramDataPoint) {
	sb := state.ExplicitBounds()
	db := dp.ExplicitBounds()
	if sb.Len() != db.Len() {
		dp.CopyTo(state)
		return
	}
	for i := 0; i < sb.Len(); i++ {
		if sb.At(i) != db.At(i) {
			dp.CopyTo(state)
			return
		}
	}
	// Boundaries match — merge.
	sc := state.BucketCounts()
	dc := dp.BucketCounts()
	merged := make([]uint64, sc.Len())
	for i := 0; i < sc.Len(); i++ {
		merged[i] = sc.At(i) + dc.At(i)
	}
	state.BucketCounts().FromRaw(merged)
	state.SetCount(state.Count() + dp.Count())
	if dp.HasSum() {
		state.SetSum(state.Sum() + dp.Sum())
	}
	if dp.HasMin() && (!state.HasMin() || dp.Min() < state.Min()) {
		state.SetMin(dp.Min())
	}
	if dp.HasMax() && (!state.HasMax() || dp.Max() > state.Max()) {
		state.SetMax(dp.Max())
	}
}

// ConvertOTelToCloudWatch converts an OTel histogram data point to a CloudWatchHistogram
// using density-aware bucket subdivision. Ported from contrib fork commit 3b43f3b3b3.
func ConvertOTelToCloudWatch(dp pmetric.HistogramDataPoint) *CloudWatchHistogram {
	const maximumInnerBucketCount = 10

	h := &CloudWatchHistogram{
		maximum:     dp.Max(),
		minimum:     dp.Min(),
		sampleCount: float64(dp.Count()),
		sum:         dp.Sum(),
	}

	bounds := dp.ExplicitBounds()
	lenBounds := bounds.Len()
	bucketCounts := dp.BucketCounts()
	lenBucketCounts := bucketCounts.Len()

	// Special case: no boundaries implies a single bucket.
	if lenBounds == 0 {
		h.counts = append(h.counts, float64(bucketCounts.At(0)))
		switch {
		case dp.HasMax() && dp.HasMin():
			h.values = append(h.values, h.minimum/2.0+h.maximum/2.0)
		case dp.HasMax():
			h.values = append(h.values, h.maximum)
		case dp.HasMin():
			h.values = append(h.values, h.minimum)
		default:
			h.values = append(h.values, 0)
		}
		return h
	}

	// Infer min if not present.
	if !dp.HasMin() {
		minBucketIdx := 0
		for i := 0; i < lenBucketCounts; i++ {
			if bucketCounts.At(i) > 0 {
				minBucketIdx = i
				break
			}
		}
		if minBucketIdx != 0 {
			h.minimum = bounds.At(minBucketIdx - 1)
		} else {
			bucketWidth := 0.001
			if lenBounds > 1 {
				bucketWidth = bounds.At(1) - bounds.At(0)
			}
			h.minimum = bounds.At(0) - bucketWidth
			if bounds.At(0) >= 0 {
				h.minimum = max(h.minimum, 0.0)
			}
		}
	}

	// Infer max if not present.
	if !dp.HasMax() {
		maxBucketIdx := lenBounds - 1
		for i := lenBucketCounts - 1; i >= 0; i-- {
			if bucketCounts.At(i) > 0 {
				maxBucketIdx = i
				break
			}
		}
		if maxBucketIdx <= lenBounds-1 {
			h.maximum = bounds.At(maxBucketIdx)
		} else {
			bucketWidth := 0.01
			if lenBounds > 1 {
				bucketWidth = bounds.At(lenBounds-1) - bounds.At(lenBounds-2)
			}
			h.maximum = bounds.At(lenBounds-1) + bucketWidth
		}
	}

	// Pre-calculate total output size.
	totalOutputSize := 0
	for i := 0; i < lenBucketCounts; i++ {
		sc := bucketCounts.At(i)
		if sc > 0 {
			totalOutputSize += int(min(sc, maximumInnerBucketCount))
		}
	}
	if totalOutputSize == 0 {
		return h
	}

	h.values = make([]float64, 0, totalOutputSize)
	h.counts = make([]float64, 0, totalOutputSize)

	for i := 0; i < lenBucketCounts; i++ {
		sampleCount := int(bucketCounts.At(i))
		if sampleCount == 0 {
			continue
		}

		lowerBound := h.minimum
		if i > 0 {
			lowerBound = bounds.At(i - 1)
		}
		upperBound := h.maximum
		if i < lenBucketCounts-1 {
			upperBound = bounds.At(i)
		}

		// Compute density ratio to decide distribution shape.
		ratio := 0.0
		if i < lenBucketCounts-1 {
			nextSampleCount := bucketCounts.At(i + 1)
			if nextSampleCount == 0 {
				ratio = 0.0
			} else {
				var nextUpperBound float64
				if i+1 == lenBucketCounts-1 {
					nextUpperBound = h.maximum
				} else {
					nextUpperBound = bounds.At(i + 1)
				}
				denom := (nextUpperBound - upperBound) * float64(sampleCount)
				numerator := (upperBound - lowerBound) * float64(nextSampleCount)
				ratio = numerator / denom
			}
		}

		innerBucketCount := min(sampleCount, maximumInnerBucketCount)
		delta := (upperBound - lowerBound) / float64(innerBucketCount)

		switch {
		case ratio < 1.0/math.E: // Density decreasing: inverse quadratic.
			sigma := float64(sumOfSquares(innerBucketCount))
			epsilon := float64(sampleCount) / sigma
			entryStart := len(h.counts)

			runningSum := 0
			for j := 0; j < innerBucketCount; j++ {
				v := epsilon * float64((j-innerBucketCount)*(j-innerBucketCount))
				adj := int(math.Floor(v))
				if adj > 0 {
					runningSum += adj
					h.values = append(h.values, lowerBound+delta*float64(j+1))
					h.counts = append(h.counts, float64(adj))
				}
			}
			remainder := sampleCount - runningSum
			if len(h.counts) < entryStart+remainder {
				h.counts = append(h.counts, make([]float64, remainder)...)
				h.values = append(h.values, make([]float64, remainder)...)
			}
			for j := 0; j < remainder; j++ {
				h.counts[entryStart]++
				entryStart++
			}

		case ratio < math.E: // Density stable: uniform.
			baseCount := sampleCount / innerBucketCount
			remainder := sampleCount % innerBucketCount
			for j := 1; j <= innerBucketCount; j++ {
				count := baseCount
				if j <= remainder {
					count++
				}
				h.values = append(h.values, lowerBound+delta*float64(j))
				h.counts = append(h.counts, float64(count))
			}

		default: // Density increasing: quadratic.
			sigma := float64(sumOfSquares(innerBucketCount))
			epsilon := float64(sampleCount) / sigma
			entryStart := len(h.counts)

			runningSum := 0
			for j := 1; j <= innerBucketCount; j++ {
				v := epsilon * float64(j*j)
				adj := int(math.Floor(v))
				if adj > 0 {
					runningSum += adj
					h.values = append(h.values, lowerBound+delta*float64(j))
					h.counts = append(h.counts, float64(adj))
				}
			}
			remainder := sampleCount - runningSum
			if len(h.counts) < entryStart+remainder {
				h.counts = append(h.counts, make([]float64, remainder)...)
				h.values = append(h.values, make([]float64, remainder)...)
			}
			entryStart = len(h.counts) - 1
			for j := 0; j < remainder; j++ {
				h.counts[entryStart]++
				entryStart--
			}
		}
	}

	// Snap last non-zero entry to maximum.
	if dp.HasMax() && len(h.values) > 0 {
		lastIdx := len(h.values) - 1
		for i := lastIdx; i >= 0; i-- {
			if h.counts[i] > 0 {
				lastIdx = i
				break
			}
		}
		h.values[lastIdx] = h.maximum
		h.values = h.values[:lastIdx+1]
		h.counts = h.counts[:lastIdx+1]
	}

	return h
}

// sumOfSquares returns Σi² for i=1..n using the closed-form formula.
func sumOfSquares(n int) int {
	return n * (n + 1) * (2*n + 1) / 6
}
