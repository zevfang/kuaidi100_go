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

func GetComByCodeName(afu_code, name string) string {
	var kd100_code string
	for _, v := range *coms {
		if v.AfuCode == afu_code && v.Name == name {
			kd100_code = v.Kd100Code
			break
		}
	}
	return kd100_code
}

func GetAfuCodeByCodeName(kd_code, name string) string {
	var afu_code string
	for _, v := range *coms {
		if v.AfuCode == kd_code && v.Name == name {
			afu_code = v.AfuCode
			break
		}
	}
	return afu_code
}
