package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func CreateOutput(servers []Server) table.Table {
	headerFmt := color.New(color.FgWhite, color.Bold, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgHiBlue).SprintfFunc()

	tbl := table.New("Server", "Name", "Index", "Free Memory", "Used Memory", "Total Memory")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, server := range servers {
		name := server.Name
		for _, gpu := range server.Devices {
			tbl.AddRow(name, gpu.Name, gpu.Index, AlignRight(gpu.FreeMemory), AlignRight(gpu.UsedMemory), AlignRight(gpu.TotalMemory))
			name = ""
		}
	}

	return tbl
}

func AlignRight(s string) string {
	return fmt.Sprintf("%10s", s)
}

// Returns a memory value in MiB.
func GetMemoryInMB(s string) (int, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "MiB") {
		return strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(s, "MiB")))
	}

	if strings.HasSuffix(s, "GiB") {
		gib, err := strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(s, "GiB")))
		return gib * 1024, err
	}

	return 0, fmt.Errorf("Unknown memory format: %s", s)
}

// Returns a list of remote servers containing free GPUs.
func FilterFreeGPUs(servers []Server) []Server {
	return filterGPUs(servers, CheckIsGPUFree)
}

func FilterUsedGPUs(servers []Server) []Server {
	return filterGPUs(servers, func(gpu GPU) bool {
		return !CheckIsGPUFree(gpu)
	})
}

func filterGPUs(servers []Server, f func(GPU) bool) []Server {
	var filtered []Server

	for _, server := range servers {
		var devices []GPU

		for _, gpu := range server.Devices {
			if f(gpu) {
				devices = append(devices, gpu)
			}
		}

		filtered = append(filtered, Server{Name: server.Name, Devices: devices})
	}

	return filtered
}
