package filewatch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"time"

	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v2"
)

var (
	unMarshaller  = yaml.Unmarshal
	watchInterval = time.Second * 5
)

// 设置默认的unMarshaller结构
func SetUnmarshal(name string) {
	if name == "yaml" {
		unMarshaller = yaml.Unmarshal
	} else if name == "json" {
		unMarshaller = json.Unmarshal
	}
}

// WatchFile ...
func WatchFile(initVal interface{}, path string, callback func(interface{}, bool)) error {
	// copy value first
	cloneInitVal := clonePtrValue(initVal)

	// read config file and callback first time
	cloneVal := clonePtrValue(cloneInitVal)
	lastVal, err := readFileToValue(path, cloneVal)
	if err != nil {
		return fmt.Errorf("read file error %v", err)
	}
	callback(cloneVal, false)

	// loop update file
	if watchInterval > 0 {
		go func() {
			for range time.NewTicker(watchInterval).C {
				cloneVal := clonePtrValue(cloneInitVal)
				currVal, err := readFileToValue(path, cloneVal)
				if err != nil {
					// log.Println("Update value error: " + err.Error())
					continue
				}
				if lastVal == currVal {
					continue
				}
				// log.Println("Update Config: " + currVal)
				lastVal = currVal
				callback(cloneVal, true)
			}
		}()
	}
	return nil
}

func readFileToValue(path string, val interface{}) (string, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bs), unMarshaller(bs, val)
}

func clonePtrValue(src interface{}) interface{} {
	v := reflect.New(reflect.ValueOf(src).Elem().Type()).Interface()
	copier.Copy(v, src)
	return v
}
