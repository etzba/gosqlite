package server

import (
	"fmt"
	"io"
	"math/bits"
	"net/http"
	"strconv"
	"strings"
)

var hist hdr // allocated once, shared by all requests

// calculatePercentileInteger accept integer in request body, calculate p90, p95, p99 and median and reply with the current p50, p90, p95, p99
// after every request.
func (s *Server) calculatePercentileInteger() func(w http.ResponseWriter, r *http.Request) { //nolint:unused
	return func(w http.ResponseWriter, r *http.Request) {
		// percentils map would look like p50=123, p90=321, p99=333
		percentiles := make(map[string]int)
		body, err := io.ReadAll(io.LimitReader(r.Body, 64))
		if err != nil {
			s.Logger.Error("failed to read request body", err)
			s.Respoder.SendBadRequest(w)
			return
		}
		defer r.Body.Close()

		value, err := strconv.ParseInt(strings.TrimSpace(string(body)), 10, 64)
		if err != nil || value < 0 {
			s.Logger.Error("request body must be a non-negative integer", err)
			s.Respoder.SendBadRequest(w)
			return
		}

		mu.Lock()
		defer mu.Unlock()
		hist.record(value)
		percentiles["p50"] = int(hist.valueAtPercentile(50))
		percentiles["p75"] = int(hist.valueAtPercentile(75))
		percentiles["p90"] = int(hist.valueAtPercentile(90))
		percentiles["p95"] = int(hist.valueAtPercentile(95))
		percentiles["p99"] = int(hist.valueAtPercentile(99))
		s.Logger.Info(fmt.Sprintf("percentiles calculation: \tp50 = %d\tp75 = %d\tp90 = %d\tp95 = %d\tp99 = %d", int(hist.valueAtPercentile(50)), int(hist.valueAtPercentile(75)), int(hist.valueAtPercentile(90)), int(hist.valueAtPercentile(95)), int(hist.valueAtPercentile(99))))
		s.Respoder.SendOK(w, percentiles)
		return
	}
}

const (
	subBucketBits  = 11                  // 2048 sub-buckets per magnitude => ~0.1% precision
	subBucketCount = 1 << subBucketBits  // 2048
	subBucketHalf  = subBucketCount >> 1 // 1024
	// Covers every non-negative int64: max shift = 63 - subBucketBits = 52.
	hdrBucketCount = (63 - subBucketBits + 2) * subBucketHalf // 55296 counters (~432 KB), fixed
)

// hdr is a fixed-size HDR-style histogram. Memory never grows with input.
type hdr struct {
	counts [hdrBucketCount]uint64
	total  uint64
}

// index maps a non-negative value to its counter.
func hdrIndex(v int64) int {
	shift := bits.Len64(uint64(v)) - subBucketBits
	if shift < 0 {
		shift = 0
	}
	return shift*subBucketHalf + int(v>>uint(shift))
}

// highestEquivalent returns the highest value that falls into the given counter.
func hdrValue(idx int) int64 {
	if idx < subBucketCount {
		return int64(idx) // exact range
	}
	shift := idx/subBucketHalf - 1
	sub := int64(idx - shift*subBucketHalf)
	return ((sub + 1) << uint(shift)) - 1
}

func (h *hdr) record(v int64) {
	h.counts[hdrIndex(v)]++
	h.total++
}

// valueAtPercentile returns the value at percentile p (0 < p <= 100).
func (h *hdr) valueAtPercentile(p uint64) int64 {
	if h.total == 0 {
		return 0
	}
	target := (h.total*p + 99) / 100 // ceil(total * p / 100)
	if target == 0 {
		target = 1
	}
	var cum uint64
	for i := range h.counts {
		cum += h.counts[i]
		if cum >= target {
			return hdrValue(i)
		}
	}
	return 0
}
