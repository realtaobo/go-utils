package filewatch

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Yaml配置解析
// 用于解析服务配置文件中的TCCM配置
type TCCM struct {
	ReportIP    string `yaml:"reportIP"`
	AppID       int64  `yaml:"appId"`
	Namespace   string `yaml:"namespace"`
	Type        string `yaml:"type"`
	MeasureMent string `yaml:"measurement"`
}

// 用于解析服务配置文件中的日志文件，agent配置文件，监听端口等消息
type Config struct {
	Port       string `yaml:"port"`
	ConfigFile string `yaml:"configFile"`
	LogFile    string `yaml:"logFile"`
	CMonitor   TCCM   `yaml:"TCCM"`
}

// 读取yaml文件并返回对应的结构体
func ReadYaml(path string) Config {
	setting := Config{}
	config, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(config, &setting); err != nil {
		panic(err)
	}
	return setting
}

// 返回v的md5值
func RetMd5(v interface{}) string {
	x, _ := yaml.Marshal(v)
	return fmt.Sprintf("%x", md5.Sum(x))
}

//	用于解析自身的agent配置结构
type Params struct {
	MD5     string `json:"md5"`
	URL     string `json:"url"`
	Version string `json:"version"`
	Extra   string `json:"extra"` // 指明当前版本是否为默认版本
}

// agent配置对应结构体
type Agent struct {
	Cls     []Params `json:"cls"`
	Jaeger  []Params `json:"jaeger"`
	Monitor []Params `json:"monitor"`
}

// 读取json文件
func ReadJson(path string) *Agent {
	setting := new(Agent)
	config, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(config, &setting); err != nil {
		panic(err)
	}
	return setting
}
