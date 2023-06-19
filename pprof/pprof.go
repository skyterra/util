package pprof

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
)

func ProfCPUStart() *os.File {
	err := os.MkdirAll("../bin", 0700)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("../bin/cpu.pprof")
	if err != nil {
		log.Fatal(err)
	}

	if err = pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}

	return f
}

func ProfCPUStop(f *os.File) {
	pprof.StopCPUProfile()
	f.Close()

	output, _ := exec.Command("go", "tool", "pprof", "-top", "../bin/cpu.pprof").CombinedOutput()
	r := bytes.NewReader(output)
	s := bufio.NewScanner(r)

	// 打印Top 20
	lineCount := 26
	fmt.Println("--- CPU Top20 ---")
	for lineCount > 0 && s.Scan() {
		fmt.Println(s.Text())
		lineCount--
	}

	fmt.Println()
}

func ProfMemory() {
	err := os.MkdirAll("../bin", 0700)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("../bin/memory.pprof")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	runtime.GC()
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	output, _ := exec.Command("go", "tool", "pprof", "-top", "../bin/memory.pprof").CombinedOutput()
	r := bytes.NewReader(output)
	s := bufio.NewScanner(r)

	// 打印Top 20
	lineCount := 26
	fmt.Println("--- Mem Top20 ---")
	for lineCount > 0 && s.Scan() {
		fmt.Println(s.Text())
		lineCount--
	}
	fmt.Println()
}
