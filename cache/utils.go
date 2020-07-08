package cache

import (
	"github.com/shirou/gopsutil/mem"
	"log"
	"runtime"
)

type MemStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
	Self uint64 `json:"self"`
}

func GetMemStatus() MemStatus {
	//自身占用
	memStat := new(runtime.MemStats)
	runtime.ReadMemStats(memStat)
	memStatus := MemStatus{}
	memStatus.Self = memStat.Alloc / 1024

	vm, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal("Could not retrieve RAM details.", err)
	} else {
		memStatus.All = vm.Total / 1024
		memStatus.Used = (vm.Total - vm.Available) / 1024
		memStatus.Free = vm.Available / 1024
	}

	return memStatus
}
