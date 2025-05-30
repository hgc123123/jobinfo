package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	ages := make(map[int]int)

	arr := [8]int{1884728,1884729,1884730,1792116,1792215,1792342,1792501,1875599}

	for _, i := range arr {
		cmd := exec.Command("sh", "-c", fmt.Sprintf("nvidia-smi pmon -s u -c 1 | grep %d | awk '{print $1}'", i))

		output, err := cmd.Output()
		if err != nil {
			log.Fatalf("Failed to execute command: %s", err)
		}

		outputStr := string(output)
		lines := strings.Split(outputStr, "\n")
		for _, line := range lines {
			if len(strings.TrimSpace(line)) > 0 {
				gpuNumber, _ := strconv.Atoi(strings.TrimSpace(line))
				cmd1 := exec.Command("sh", "-c", fmt.Sprintf("nvidia-smi pmon -s u -c 1 | grep %d | awk '{print $4}'", i))

				output1, err := cmd1.Output()
				if err != nil {
					log.Fatalf("Failed to execute command: %s", err)
				}

				outputStr1 := string(output1)
				lines1 := strings.Split(outputStr1, "\n")
				for _, line1 := range lines1 {
					if len(strings.TrimSpace(line1)) > 0 {
						gpuUtil, _ := strconv.Atoi(strings.TrimSpace(line1))
						ages[gpuNumber] = gpuUtil
						fmt.Printf("%d is %d years old\n", gpuNumber, gpuUtil)
					}
				}
			}
		}
	}

	if len(ages) == 0 {
		fmt.Println("No data available")
	} else {
		for gpuNumber, gpuUtil := range ages {
			fmt.Printf("%d is %d years old\n", gpuNumber, gpuUtil)
		}
	}
}
