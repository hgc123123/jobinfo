package cgroups

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type cpuacct string

// GetUsagePerCPU returns the per-nanosecond CPU usage of each CPU indexed from
// 0
func (c cpuacct) GetUsagePerCPU() ([]int, error) {
	var usage []int
	data, err := readFile(string(c), "cpuacct.usage_percpu")
	/*
	The content of "string(c)" is:
	/sys/fs/cgroup/cpuacct/slurm/uid_42xxxx/job_242184/step_0/task_0
	*/
	if err != nil {
		return usage, err
	}
	for _, usageStr := range strings.Split(strings.TrimSpace(data), " ") {
		usageInt, err := strconv.Atoi(strings.TrimSpace(usageStr))
		if err != nil {
			log.Errorf("unable to convert per-cpu usage to integer: %v", err)
			return usage, err
		}
		usage = append(usage, usageInt)
	}
	return usage, nil
}
