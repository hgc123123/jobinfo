package cgroups

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	procCgroupIdxSubsystems = 1
	procCgroupIdxPath       = 2
)

type Cgroups struct {
	Cpuset  cpuset
	Cpuacct cpuacct
	Memory  memory
	Devices devices
}

func readFile(root string, filename string) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join(root, filename))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func LoadCgroups(specPath string, cgroupsRootPath string) (Cgroups, error) {
	var cgroups Cgroups
	cgroupsPath := filepath.Clean(specPath)
	cgroupsFile, err := os.Open(cgroupsPath)
	if err != nil {
		return cgroups, err
	}
	defer cgroupsFile.Close()
	csvReader := csv.NewReader(cgroupsFile)
	csvReader.Comma = ':'
	csvLines, err := csvReader.ReadAll()
	if err != nil {
		return cgroups, err
	}
	for _, csvLine := range csvLines {
		subsystems := strings.Split(csvLine[procCgroupIdxSubsystems], ",")
		for _, subsystem := range subsystems {
			if len(subsystem) < 1 {
				log.Debug("skipping empty subsystem")
				continue
			}
			cgroupAbsolutePath := filepath.Join(cgroupsRootPath, strings.TrimPrefix(subsystem, "name="), csvLine[procCgroupIdxPath])
			if _, err := os.Stat(cgroupAbsolutePath); os.IsNotExist(err) {
				return cgroups, fmt.Errorf("cgroup path doesn't exist: %s", cgroupAbsolutePath)
			}
			switch subsystem {
			case "cpuset":
				cgroups.Cpuset = cpuset(cgroupAbsolutePath)
			case "cpuacct":
				cgroups.Cpuacct = cpuacct(cgroupAbsolutePath)
			case "memory":
				cgroups.Memory = memory(cgroupAbsolutePath)
			case "devices":
                                cgroups.Devices = devices(cgroupAbsolutePath)
			default:
				log.Debugf("skipping unimplemented subsystem: %v", subsystem)
			}
		}
	}
	return cgroups, nil
}

func LoadProcessCgroups(pid int, cgroupsRootPath string) (Cgroups, error) {
	cgroupsPath := filepath.Join("/proc", strconv.Itoa(pid), "cgroup")
	return LoadCgroups(cgroupsPath, cgroupsRootPath)
}
