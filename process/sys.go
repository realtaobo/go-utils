package process

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// InitFunc
type InitFunc func() error

// ProcessInfo 进程详情
type ProcessInfo struct {
	Path     string
	Args     []string
	InitFunc InitFunc
	WorkDir  string
}

// TplWrite 通过模板变量values重写src文件
// 替换src文件中的key变量,并生成新文件到dst
func TplWrite(src, dst string, values map[string]string) error {
	tplBs, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("readFile error %v", err)
	}
	tpl := string(tplBs)
	for k, v := range values {
		tpl = strings.Replace(tpl, k, v, -1)
	}
	return ioutil.WriteFile(dst, []byte(tpl), 0644)
}

// 回收进程
// 待寻找平台通用包

func listProcess() map[string]bool {
	ret := make(map[string]bool)
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		return ret
	}

	for _, f := range files {
		fpath := "/proc/" + f.Name() + "/exe"
		rpath, _ := os.Readlink(fpath)
		if rpath != "" {
			ret[rpath] = true
		}
	}
	return ret
}

func startAgent(a *ProcessInfo) error {
	cmd := exec.Command(a.Path, a.Args...)
	if a.WorkDir != "" {
		cmd.Dir = a.WorkDir
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start process error %v", err)
	} else {
		log.Printf("start process successfully %s %v", a.Path, a.Args)
	}
	return nil
}

// CheckRestartAgent 如果agent进程不在则重启
func CheckRestartAgent(aa []*ProcessInfo) error {
	pList := listProcess()
	var retErr error
	for _, a := range aa {
		if pList[a.Path] {
			continue
		}
		if err := startAgent(a); err != nil {
			retErr = err
		}
	}
	return retErr
}

type DataSource struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Msg  string `json:"message"`
}

// 创建数据源
func CreateDataSource() error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("reportMetrics panic, %v", err)
		}
	}()
	url := "http://admin:admin@localhost:3000/api/datasources"
	fmt.Println("url:>", url)
	//json序列化
	post := `
	  {
		"name":"Prometheus",
		"type":"prometheus",
		"url":"http://localhost:9090",
		"access":"proxy",
		"basicAuth":true
      }
	`

	fmt.Println(url, "post", post)
	var jsonStr = []byte(post)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("status", resp.Status)
	// fmt.Println("response:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	return nil
}
