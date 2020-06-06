package histogram

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/rakyll/spannerbench/internal/stats"
)

const barChar = "■"

type Bucket struct {
	Mark      int64
	Count     int
	Frequency float64
}

type Histogram struct {
	buckets []Bucket
}

const numBuckets = 5

func NewHistogram(dur []int64) *Histogram {
	// TODO(jbd): Handle when len(dur) < 5.
	buckets := make([]int64, numBuckets+1)
	counts := make([]int, numBuckets+1)
	dur = stats.SortInt64s(dur)

	if len(dur) < 5 {
		return nil
	}

	slowest := dur[len(dur)-1]
	fastest := dur[0]

	bs := float64(slowest-fastest) / numBuckets
	for i := 0; i < numBuckets; i++ {
		buckets[i] = fastest + int64(float64(i)*bs)
	}
	buckets[numBuckets] = slowest

	var bi int
	var max int
	for i := 0; i < len(dur); {
		if dur[i] <= buckets[bi] {
			i++
			counts[bi]++
			if max < counts[bi] {
				max = counts[bi]
			}
		} else if bi < len(buckets)-1 {
			bi++
		}
	}
	res := make([]Bucket, len(buckets))
	for i := 0; i < len(buckets); i++ {
		res[i] = Bucket{
			Mark:      buckets[i],
			Count:     counts[i],
			Frequency: float64(counts[i]) / float64(len(dur)),
		}
	}
	return &Histogram{
		buckets: res,
	}
}

func (h *Histogram) String() string {
	max := 0
	for _, b := range h.buckets {
		if v := b.Count; v > max {
			max = v
		}
	}
	res := new(bytes.Buffer)
	for i := 0; i < len(h.buckets); i++ {
		// Normalize bar lengths.
		bucket := h.buckets[i]
		var barLen int
		if max > 0 {
			barLen = (bucket.Count*20 + max/2) / max
		}
		if bucket.Count > 0 && barLen == 0 {
			barLen = 1 // bar length should be at least one if there are items.
		}

		dur := time.Duration(bucket.Mark)
		res.WriteString(fmt.Sprintf("  %-10v: %v ", dur, strings.Repeat(barChar, barLen)))
		if bucket.Count > 0 {
			res.WriteString(fmt.Sprintf("(%v)", bucket.Count))
		}
		res.WriteRune('\n')
	}
	return res.String()
}
