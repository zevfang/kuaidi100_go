package system

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

type Com struct {
	AfuCode   string `json:"afu_code"`
	Kd100Code string `json:"kd100_code"`
	Name      string `json:"name"`
}

var coms *[]Com

func LoadComs(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(b), &coms)
}

func GetComArray() *[]Com {
	return coms
}