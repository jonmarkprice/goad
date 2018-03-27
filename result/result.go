package result

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/goadapp/goad/api"
	"github.com/goadapp/goad/goad/util"
)

// AggData type
type AggData struct {
	TotalReqs            int
	TotalTimedOut        int
	TotalConnectionError int
	AveTimeToFirst       int64
	TotBytesRead         int
	Statuses             map[string]int
	AveTimeForReq        int64
	AveReqPerSec         float64
	TimeDelta            time.Duration
	AveKBytesPerSec      float64
	Slowest              int64
	Fastest              int64
	Region               string
	FatalError           string
	Finished             bool
	StartTime            time.Time
	EndTime              time.Time

	// new
	SumReqTime     int64
	SumReqSq       int64
	// ReqTimesBinned map[int64]int
}

// LambdaResults type
type LambdaResults struct {
	Lambdas []AggData
}

// Regions the LambdaResults were collected from
func (r *LambdaResults) Regions() []string {
	regions := make([]string, 0)
	for _, lambda := range r.Lambdas {
		if lambda.Region != "" {
			regions = append(regions, lambda.Region)
		}
	}
	regions = util.RemoveDuplicates(regions)
	sort.Strings(regions)
	return regions
}

// RegionsData aggregates the individual lambda functions results per region
func (r *LambdaResults) RegionsData() map[string]AggData {
	regionsMap := make(map[string]AggData)
	for _, region := range r.Regions() {
		regionsMap[region] = sumAggData(r.ResultsForRegion(region))
	}
	return regionsMap
}

//SumAllLambdas aggregates results of all Lambda functions
func (r *LambdaResults) SumAllLambdas() AggData {
	return sumAggData(r.Lambdas)
}

//ResultsForRegion return the sum of results for a given regions
func (r *LambdaResults) ResultsForRegion(region string) []AggData {
	lambdasOfRegion := make([]AggData, 0)
	for _, lambda := range r.Lambdas {
		if lambda.Region == region {
			lambdasOfRegion = append(lambdasOfRegion, lambda)
		}
	}
	return lambdasOfRegion
}

func SetupRegionsAggData(lambdaCount int) *LambdaResults {
	lambdaResults := &LambdaResults{
		Lambdas: make([]AggData, lambdaCount),
	}
	for i := 0; i < lambdaCount; i++ {
		lambdaResults.Lambdas[i].Statuses = make(map[string]int)
		// lambdaResults.Lambdas[i].ReqTimesBinned = make(map[int64]int)
	}
	return lambdaResults
}

func sumAggData(dataArray []AggData) AggData {
	sum := AggData{
		Fastest:        math.MaxInt64,
		Statuses:       make(map[string]int),
		// ReqTimesBinned: make(map[int64]int),
		Finished:       true,
	}
	for _, lambda := range dataArray {
		sum.AveKBytesPerSec += lambda.AveKBytesPerSec
		sum.AveReqPerSec += lambda.AveReqPerSec
		sum.AveTimeForReq += lambda.AveTimeForReq
		sum.AveTimeToFirst += lambda.AveTimeToFirst
		if lambda.Fastest < sum.Fastest {
			sum.Fastest = lambda.Fastest
		}
		sum.FatalError += lambda.FatalError
		if !lambda.Finished {
			sum.Finished = false
		}
		sum.Region = lambda.Region
		if lambda.Slowest > sum.Slowest {
			sum.Slowest = lambda.Slowest
		}
		for key := range lambda.Statuses {
			sum.Statuses[key] += lambda.Statuses[key]
		}
		//for key := range lambda.ReqTimesBinned {
		//	sum.ReqTimesBinned[key] += lambda.ReqTimesBinned[key]
		//}
		if sum.StartTime.IsZero() || lambda.StartTime.Before(sum.StartTime) {
			sum.StartTime = lambda.StartTime
		}
		if lambda.EndTime.After(sum.EndTime) {
			sum.EndTime = lambda.EndTime
		}

		sum.SumReqTime += lambda.SumReqTime
		sum.SumReqSq += lambda.SumReqSq
		sum.TimeDelta += lambda.TimeDelta
		sum.TotalConnectionError += lambda.TotalConnectionError
		sum.TotalReqs += lambda.TotalReqs
		sum.TotalTimedOut += lambda.TotalTimedOut
		sum.TotBytesRead += lambda.TotBytesRead
	}

	// Debugging
	// In case there was no data
	if len(dataArray) > 0 {
		sum.AveTimeForReq = sum.AveTimeForReq / int64(len(dataArray))
		sum.AveTimeToFirst = sum.AveTimeToFirst / int64(len(dataArray))
	} else {
		fmt.Println("No data for region ", sum.Region)
	}
	return sum
}

func (r *LambdaResults) AllLambdasFinished() bool {
	for _, lambda := range r.Lambdas {
		if !lambda.Finished {
			return false
		}
	}
	return true
}

func AddResult(data *AggData, result *api.RunnerResult) {
	initCountOk := int64(data.TotalReqs - data.TotalTimedOut - data.TotalConnectionError)
	addCountOk := int64(result.RequestCount - result.TimedOut - result.ConnectionErrors)
	totalCountOk := initCountOk + addCountOk

	data.TotalReqs += result.RequestCount
	data.TotalTimedOut += result.TimedOut
	data.TotalConnectionError += result.ConnectionErrors
	data.TotBytesRead += result.BytesRead
	data.TimeDelta += result.TimeDelta

	if totalCountOk > 0 {
		data.AveTimeToFirst = addToTotalAverage(data.AveTimeToFirst, initCountOk, result.AveTimeToFirst, addCountOk)
		data.AveTimeForReq = addToTotalAverage(data.AveTimeForReq, initCountOk, result.AveTimeForReq, addCountOk)
		data.AveKBytesPerSec = float64(data.TotBytesRead) / float64(data.TimeDelta.Seconds())
		data.AveReqPerSec = float64(data.TotalReqs) / float64(data.TimeDelta.Seconds())
	}

	if data.StartTime.IsZero() || result.StartTime.Before(data.StartTime) {
		data.StartTime = result.StartTime
	}

	if result.EndTime.After(data.EndTime) {
		data.EndTime = result.EndTime
	}

	// Aggregate maps // TODO this could be a fn.
	for key, value := range result.Statuses {
		data.Statuses[key] += value
	}

	//for key, value := range result.ReqTimesBinned {
	//	data.ReqTimesBinned[key] += value
	//}

	data.SumReqTime += result.SumReqTime
	data.SumReqSq += result.SumReqSq

	if result.Slowest > data.Slowest {
		data.Slowest = result.Slowest
	}

	if result.Fastest > 0 && (data.Fastest == 0 || result.Fastest < data.Fastest) {
		data.Fastest = result.Fastest
	}
	data.Finished = result.Finished
	data.Region = result.Region
}

func addToTotalAverage(currentAvg, currentCount, addAvg, addCount int64) int64 {
	return ((currentAvg * currentCount) + (addAvg * addCount)) / (currentCount + addCount)
}

func addToTotalAverageFloat(currentAvg, currentCount, addAvg, addCount float64) float64 {
	return ((currentAvg * currentCount) + (addAvg * addCount)) / (currentCount + addCount)
}
