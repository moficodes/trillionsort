package stat

import (
	"fmt"
	"runtime"
)

// MemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func MemUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	message := fmt.Sprintf("alloc = %s, totalAlloc = %s, sys = %s, numGC = %v",
		HumanReadableFilesize(int64(m.Alloc)),
		HumanReadableFilesize(int64(m.TotalAlloc)),
		HumanReadableFilesize(int64(m.Sys)),
		m.NumGC)
	return message
}
