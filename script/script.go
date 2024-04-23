package script

import (
	"strconv"
	"strings"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type script string

// GetUsageInBytes returns the current memory in use by the cgroup in bytes
func (c script) GetContentOfScript() (int, error) {
	data, err := readFile(string(c))
	if err != nil {
		return 0, err
	}
	fmt.Println(string(data))
}
