package process

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"golang.org/x/sys/unix"
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
func ReapProcess() {
	cc := make(chan os.Signal, 1)
	signal.Notify(cc, unix.SIGCHLD)
	for range cc {
		func() {
			for {
				var status unix.WaitStatus
				pid, err := unix.Wait4(-1, &status, unix.WNOHANG, nil)
				switch err {
				case nil:
					if pid > 0 {
						continue
					}
					return
				case unix.ECHILD:
					return
				case unix.EINTR:
					continue
				default:
					return
				}
			}
		}()
	}
}

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
