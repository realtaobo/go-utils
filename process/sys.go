package process

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// ProcessInfo 进程详情
type ProcessInfo struct {
	Path    string
	Args    []string
	WorkDir string
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

// 启动指定的单个进程
func StartProcess(a *ProcessInfo) error {
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

// ListProcess 列出当前系统正在运行的进程, key为二进制文件的路径
// 只适用于Linux操作系统
func ListProcess() map[string]bool {
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

// CheckRestartProcess 如果指定的进程不存在则重启进程
// 暂时只适用于linux操作系统
func CheckRestartProcess(procs []*ProcessInfo) error {
	pList := ListProcess()
	var retErr error
	for _, proc := range procs {
		if pList[proc.Path] {
			continue
		}
		if err := StartProcess(proc); err != nil {
			retErr = err
		}
	}
	return retErr
}
