package filewatch

import (
	"log"
	"sync"
	"testing"
)

func TestYaml(t *testing.T) {
	path := "./conf/server.yaml"
	nowVal := ReadYaml(path)
	v := &nowVal
	///////////////////////
	var wg sync.WaitGroup
	wg.Add(1)
	err := WatchFile(v, path, func(val interface{}, isUpdate bool) {
		vv := val.(*Config)
		if isUpdate {
			log.Println("change file\n", vv)
			wg.Done()
		} else {
			log.Println("no change")
		}
	})
	if err != nil {
		panic(err)
	}
	wg.Wait()
}

func TestJson(t *testing.T) {
	SetUnmarshal("json")
	path := "./conf/agent.json"
	nowVal := ReadJson(path)
	///////////////////////
	var wg sync.WaitGroup
	wg.Add(1)
	err := WatchFile(nowVal, path, func(val interface{}, isUpdate bool) {
		vv := val.(*Agent)
		if isUpdate {
			log.Println("change file\n", vv)
			wg.Done()
		} else {
			log.Println("no change")
		}
	})
	if err != nil {
		panic(err)
	}
	wg.Wait()
}