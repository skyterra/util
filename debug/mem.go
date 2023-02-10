package debug

import (
	"fmt"
	"runtime"
)

const (
	Byte = 1
	KB   = 1024 * Byte
	MB   = 1024 * KB
)

// PrintMemUsage 打印内存使用情况
//  https://golang.org/pkg/runtime/#MemStats
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc:%dMB", ByteToMB(m.Alloc))
	fmt.Printf("  TotalAlloc:%dMB", ByteToMB(m.TotalAlloc))
	fmt.Printf("  Sys:%dMB", ByteToMB(m.Sys))
	fmt.Printf("  NumGC:%d\n", m.NumGC)
}

func ByteToMB(size uint64) uint64 {
	return size / MB
}
