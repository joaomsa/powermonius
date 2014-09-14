package main

import (
	"fmt"
	"github.com/kballard/go-shellquote"
	"os/exec"
)

var cmd_start = "transmission-gtk -m"
var cmd_stop = "pkill -f transmission-gtk"
var cmd_status = "pgrep -f transmission-gtk"

func main() {
	//argv, split_err := shellquote.Split(cmd_start)
	//argv, split_err := shellquote.Split(cmd_stop)
	argv, split_err := shellquote.Split(cmd_status)

	if split_err != nil {
		fmt.Println(split_err.Error())
		return
	}

	cmd := exec.Command(argv[0], argv[1:]...)
	//err := cmd.Start()
	err := cmd.Run()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
