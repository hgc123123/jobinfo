package cgroups

import (
	"strconv"
	"strings"
	"os/exec"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type devices string

func (c devices) GetUsagePerGPU() (map[int]int, error) {
        usage := make(map[int]int)
	data, err := readFile(string(c), "cgroup.procs")
	if err != nil {
		return usage, err
	}
	for _, procStr := range strings.Split(strings.TrimSpace(data), "\n")[1:] {
		procInt, err := strconv.Atoi(strings.TrimSpace(procStr))
		if err != nil {
			log.Errorf("unable to convert per-gpu usage to integer: %v", err)
			return usage, err
		}
                cmd := exec.Command("sh", "-c", fmt.Sprintf("nvidia-smi pmon -s u -c 1 | grep %d | awk '{print $1}'", procInt))

		output, err := cmd.Output()
		if err != nil {
			log.Fatalf("Failed to execute command: %s", err)
		}

		outputStr := string(output)
		lines := strings.Split(outputStr, "\n")
		for _, line := range lines {
			if len(strings.TrimSpace(line)) > 0 {
				gpuNumber, _ := strconv.Atoi(strings.TrimSpace(line))
				cmd1 := exec.Command("sh", "-c", fmt.Sprintf("nvidia-smi pmon -s u -c 1 | grep %d | awk '{print $4}'", procInt))

				output1, err := cmd1.Output()
				if err != nil {
					log.Fatalf("Failed to execute command: %s", err)
				}

				outputStr1 := string(output1)
				lines1 := strings.Split(outputStr1, "\n")
				for _, line1 := range lines1 {
					if len(strings.TrimSpace(line1)) > 0 {
						gpuUtil, _ := strconv.Atoi(strings.TrimSpace(line1))
						usage[gpuNumber] = gpuUtil
						//fmt.Printf("%d is %d years old\n", gpuNumber, gpuUtil)
					}
				}
			}
		}
	}
	return usage, nil 
}

func (c devices) GetVRAMUsagePerGPU() (map[int]int, error) {
        usage := make(map[int]int)
        data, err := readFile(string(c), "cgroup.procs")
        if err != nil {
                return usage, err
        }
        for _, procStr := range strings.Split(strings.TrimSpace(data), "\n")[1:] {
                procInt, err := strconv.Atoi(strings.TrimSpace(procStr))
                if err != nil {
                        log.Errorf("unable to convert per-gpu usage to integer: %v", err)
                        return usage, err
                }
                //usage = append(usage, usageInt)
                cmd := exec.Command("sh", "-c", fmt.Sprintf("nvidia-smi pmon -s u -c 1 | grep %d | awk '{print $1}'", procInt))

                output, err := cmd.Output()
                if err != nil {
                        log.Fatalf("Failed to execute command: %s", err)
                }

                outputStr := string(output)
                lines := strings.Split(outputStr, "\n")
                for _, line := range lines {
                        if len(strings.TrimSpace(line)) > 0 {
                                gpuNumber, _ := strconv.Atoi(strings.TrimSpace(line))
                                cmd1 := exec.Command("sh", "-c", fmt.Sprintf("nvidia-smi --query-compute-apps=pid,used_memory --format=csv | grep %d | awk '{print $2}'", procInt))

                                output1, err := cmd1.Output()
                                if err != nil {
                                        log.Fatalf("Failed to execute command: %s", err)
                                }

                                outputStr1 := string(output1)
                                lines1 := strings.Split(outputStr1, "\n")
                                for _, line1 := range lines1 {
                                        if len(strings.TrimSpace(line1)) > 0 {
                                                gpuUtil, _ := strconv.Atoi(strings.TrimSpace(line1))
                                                usage[gpuNumber] = gpuUtil
                                                //fmt.Printf("%d is %d year old\n", gpuNumber, gpuUtil)
                                        }
                                }
                        }
                }
        }
        return usage, nil
}
