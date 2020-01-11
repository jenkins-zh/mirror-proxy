package pkg_test

import (
	"fmt"
	server "github.com/jenkins-zh/mirror-proxy/pkg"
	"testing"
	"time"
)

func TestWorkPool(t *testing.T) {
	pool := &server.WorkPool{}
	pool.InitPool(5)
	pool.AddTask(server.Task{Data: "echo data", TaskFunc: func(data interface{}) {
		fmt.Println(data)
	}})
	time.Sleep(5 * time.Second)
}
