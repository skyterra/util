package pprof

import (
	"fmt"
	"runtime"
)

const (
	Byte = 1
	KB   = 1024 * Byte
	MB   = 1024 * KB
	GB   = 1024 * MB
)

// PrintMemUsage 打印内存使用情况
//  https://golang.org/pkg/runtime/#MemStats
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc:%dMB", byteToMB(m.Alloc))
	fmt.Printf("\tTotalAlloc:%dMB", byteToMB(m.TotalAlloc))
	fmt.Printf("\tSys:%dMB", byteToMB(m.Sys))
	fmt.Printf("\tNumGC:%d\n", m.NumGC)
}

func byteToMB(size uint64) uint64 {
	return size / MB
}
