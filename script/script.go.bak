package script

import (
	"io/ioutil"
)


// GetUsageInBytes returns the current memory in use by the cgroup in bytes
func GetContentOfScript(name string) (string,error) {
	data,err := ioutil.ReadFile(name)
	if err != nil{
		return "0",err
	}
	return string(data),nil
}
