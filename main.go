package main

import (
	cmd "github.com/jenkins-zh/mirror-proxy/pkg"
)

func main() {
	//opt := &cmd.ServerOptions{}
	//opt.GetURL("")
	//opt.GetURL("2.204")
	//opt.GetURL("2.190.2")
	cmd.Execute()
}