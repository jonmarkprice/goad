package cli

import (
	"testing"
	// "errors"
	"fmt"
	"math"

	"github.com/goadapp/goad/result"
)

func TestEmptyCalcStat(t *testing.T) {
	// TODO set up a scenario
	agg := result.AggData{}

	_, err := CalculateFinalStatistics(agg)
	if err == nil { // errors.New("Got N <= 0.") {
		fmt.Println("Actual: ", err);
		t.Error("Should fail if 0 results.")
	}

	// empty := analyzedResults{}

	//if stats != empty {
	//	t.Error("Results should be empty on failure")
	//}
}

func TestNormalVar(t *testing.T) {
	// sumReqSq
	// sumReqTime
	data := result.AggData{
		SumReqSq: 1012757,
		SumReqTime: 2957,
		TotalReqs: 9,
		TotalTimedOut: 0,
		TotalConnectionError: 0,
	}

	stats, err := CalculateFinalStatistics(data)
	if err != nil {
		t.Error("Unexpected error")
	}
	if !within(0.1, stats.Variance, 5152.3) {
		t.Error("Variance")
	}

	if !within(0.1, stats.Mean, 328.56) {
		t.Error("Mean")
	}

	if !within(0.1, stats.StandardDeviation, 71.779) {
		t.Error("SD")
	}

	// Suppose we have request times
	/*
	reqTimes := []int64{423, 194, 321, 344, 409, 328, 345, 347, 246};
	squares := make([]int64, len(reqTimes))
	sum := 2957
	n := 9
	sumSq := 1012757
	*/
}

func within(delta float64, x float64, y float64) bool {
	if math.Abs(x - y) < delta {
		return true
	}
	return false
}
