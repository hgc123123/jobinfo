package script

import (
	"io/ioutil"
)


func GetContentOfScript(name string) (string,error) {
	data,err := ioutil.ReadFile(name)
	if err != nil{
		return "0",err
	}
	return string(data),nil
}
