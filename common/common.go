package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

// 读取文件内容，返回类型为string。
func ReadFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// 将 json 串转换为指定的结构体。
func JsonToInterface(jsonStr string, info interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), info)
	if err != nil {
		return err
	}
	return nil
}

// 将 yaml 串转换为指定的结构体。
func YamlToInterface(yamlStr string, info interface{}) error {
	err := yaml.Unmarshal([]byte(yamlStr), info)
	if err != nil {
		return err
	}
	return nil
}

// 对yyyymm格式的日期按给定的年，月数字返回处理之后的对应值。
// 如 CalcMonth(202101, -1, -2)返回201911。
func CalcMonth(yyyymm, years, months int) int {
	year := yyyymm / 100
	month := yyyymm % 100
	lastFetch := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	nextFetch := lastFetch.AddDate(years, months, 0)
	y, m, _ := nextFetch.Date()
	res, _ := strconv.Atoi(fmt.Sprintf("%d%02d", y, m))
	return res
}

// 计算日期相差多少月。start与end的格式为yyyymm，如202101。
// 计算值包括左右两个值，所以结果比直接计算大1。
func CalcSubMonth(start, end int) int {
	year1 := start / 100
	month1 := start % 100
	year2 := end / 100
	month2 := end % 100
	res := (year2-year1)*12 + (month2 - month1) + 1
	return res
}
