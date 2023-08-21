package prof

import (
	"fmt"
	"runtime"
	"strings"
)

// PrintMemUsage 打印内存使用情况
//  https://golang.org/pkg/runtime/#MemStats
func PrintMemUsage(tips ...string) {
	const MB = 1 << 20

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	tip := strings.Join(tips, " ")
	if len(tip) > 0 {
		fmt.Println(tip)
	}

	fmt.Printf("Alloc:%dMB", m.Alloc/MB)
	fmt.Printf("  TotalAlloc:%dMB", m.TotalAlloc/MB)
	fmt.Printf("  Sys:%dMB", m.Sys/MB)
	fmt.Printf("  NumGC:%d\n", m.NumGC)
}
