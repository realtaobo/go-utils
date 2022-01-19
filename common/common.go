package common

import (
	"encoding/json"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// 读取文件内容，返回类型为string
func ReadFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// 将 json 串转换为指定的结构体
func JsonToInterface(jsonStr string, info interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), info)
	if err != nil {
		return err
	}
	return nil
}

// 将 yaml 串转换为指定的结构体
func YamlToInterface(yamlStr string, info interface{}) error {
	err := yaml.Unmarshal([]byte(yamlStr), info)
	if err != nil {
		return err
	}
	return nil
}
