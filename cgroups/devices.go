package cgroups

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type devices string

// GetUsagePerCPU returns the per-nanosecond CPU usage of each CPU indexed from
// 0
func (c devices) GetUsagePerGPU() ([]int, error) {
        var usage []int
	data, err := readFile(string(c), "cgroup.procs")
	if err != nil {
		return usage, err
	}
	for _, usageStr := range strings.Split(strings.TrimSpace(data), "\n")[1:] {
		usageInt, err := strconv.Atoi(strings.TrimSpace(usageStr))
		if err != nil {
			log.Errorf("unable to convert per-gpu usage to integer: %v", err)
			return usage, err
		}
		usage = append(usage, usageInt)
	}
	return usage, nil 
}
