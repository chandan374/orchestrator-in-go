package worker

import (
	"log"
	"github.com/c9s/goprocinfo/linux"
)

type Stats struct {
	MemStats *linux.MemInfo
	DiskStats *linux.Disk
	CpuStats *linux.CPUStat 
	LoadStats *linux.LoadAvg
	TaskCount int
}

func (s *Stats) MemTotalKb() uint64 {
	return s.MemStats.MemTotal
}

func (s *Stats) MemAvailableKb() uint64 {
	return s.MemStats.MemFree
}

func (s *Stats) MemUsedPercent() float64 {
	return float64(s.MemStats.MemTotal-s.MemStats.MemFree) / float64(s.MemStats.MemTotal) * 100
}

func (s *Stats) DiskTotal() uint64 {
	return s.DiskStats.All
}

func (s *Stats) DiskFree() uint64 {
	return s.DiskStats.Free
}

func (s *Stats) DiskUsed() uint64 {
	return s.DiskStats.Used
}

func (s *Stats) CpuUsage() float64 {
	idle := s.CpuStats.Idle + s.CpuStats.IOWait
	nonIdle := s.CpuStats.User + s.CpuStats.Nice + s.CpuStats.System + s.CpuStats.IRQ + s.CpuStats.SoftIRQ + s.CpuStats.Steal
	total := idle + nonIdle
	return (float64(nonIdle) / float64(total)) * 100
}

func GetStats() *Stats {
	return &Stats{
		MemStats: GetMemoryInfo(),
		DiskStats: GetDiskInfo(),
		CpuStats: GetCPUInfo(),
		LoadStats: GetLoadAvg(),
	}
}

func GetMemoryInfo() *linux.MemInfo {
	meminfo, err := linux.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Printf("Failed to read memory info: %v", err)
		return &linux.MemInfo{}
	}
	return meminfo
}

func GetDiskInfo() *linux.Disk {
	diskinfo, err := linux.ReadDisk("/")
	if err != nil {
		log.Printf("Failed to read disk info: %v", err)
		return &linux.Disk{}
	}
	return diskinfo
}

func GetCPUInfo() *linux.CPUStat {
	stat, err := linux.ReadStat("/proc/stat")
	if err != nil {
		log.Printf("Failed to read CPU info: %v", err)
		return &linux.CPUStat{}
	}
	return &stat.CPUStatAll
}

func GetLoadAvg() *linux.LoadAvg {
	loadavg, err := linux.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		log.Printf("Failed to read load average: %v", err)
		return &linux.LoadAvg{}
	}
	return loadavg
}
