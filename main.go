package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/inu1255/gev"
	"github.com/inu1255/green/service"
)

func main() {
	// run()
	// 新建项目时使用该函数重命名 github.com/inu1255/green
	rename()
}

func run() {
	app := gev.New()
	maker := gev.NewRouteMaker()
	maker.AddParamManager(service.UserManager)

	maker.AddRoute(service.NewAddressServ())
	maker.AddRoute(service.NewUserServ())
	maker.AddRoute(service.NewVerifyServ())
	maker.AddRoute(service.NewFileServ())

	maker.RouteTo(app)
	app.Swagger("/api")
	app.Run(":8080")
}

func rename() {
	old := []byte("github.com/inu1255/green")
	curdir, _ := os.Getwd()
	gopath, _ := getGopathGoroot()
	var new []byte
	if strings.HasPrefix(curdir, gopath) {
		new = []byte(curdir[len(gopath)+len("/src/"):])
	} else {
		fmt.Println("当前目录不在GOPATH下")
		return
	}
	if bytes.Compare(old, new) == 0 {
		fmt.Println("当前在 github.com/inu1255/green 目录下,不用转换")
		return
	}
	fmt.Println(string(old), "---->", string(new))
	filepath.Walk("./", func(path string, f os.FileInfo, err error) error {
		if f.IsDir() && (path == "upload" || path == "api") {
			return filepath.SkipDir
		}
		if strings.HasSuffix(f.Name(), ".go") {
			body, _ := ioutil.ReadFile(path)
			body = bytes.Replace(body, old, new, -1)
			ioutil.WriteFile(path, body, 0644)
		}
		return nil
	})
}

func getGopathGoroot() (gopath, goroot string) {
	output, _ := exec.Command("go", "env").Output()
	s := string(output)
	lines := strings.Split(s, "\n")
	for _, s := range lines {
		ss := strings.Split(s, "=")
		if len(ss) > 1 {
			if ss[0] == "GOPATH" {
				gopath = strings.Trim(ss[1], "\"")
			} else if ss[0] == "GOROOT" {
				goroot = strings.Trim(ss[1], "\"")
			}
		}
	}
	return
}
