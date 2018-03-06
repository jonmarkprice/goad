package printing

import (
	"strconv"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/goadapp/goad/result"
)

const (
	nano           = 1000000000
)

func totErrors(data result.AggData) int {
	var okReqs int
	for statusStr, value := range data.Statuses {
		status, _ := strconv.Atoi(statusStr)
		if status < 400 {
			okReqs += value
		}
	}
	return data.TotalReqs - okReqs
}

func boldPrintln(msg string) {
	fmt.Printf("\033[1m%s\033[0m\n", msg)
}

func PrintData(data result.AggData) {
	boldPrintln("   TotReqs   TotBytes    AvgTime    AvgReq/s  (post)unzip")
	fmt.Printf("%10d %10s   %7.3fs  %10.2f %10s/s\n", data.TotalReqs, humanize.Bytes(uint64(data.TotBytesRead)), float64(data.AveTimeForReq)/nano, data.AveReqPerSec, humanize.Bytes(uint64(data.AveKBytesPerSec)))
	boldPrintln("   Slowest    Fastest   Timeouts  TotErrors")
	fmt.Printf("  %7.3fs   %7.3fs %10d %10d", float64(data.Slowest)/nano, float64(data.Fastest)/nano, data.TotalTimedOut, totErrors(data))
	fmt.Println("")
}


