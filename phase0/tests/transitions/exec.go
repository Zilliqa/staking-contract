package transitions

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
)

func ExecZli(arg ...string) (error, string) {
	cmd := exec.Command("zli", arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err.Error()
	}
	defer stdout.Close()
	fmt.Println(cmd.String())
	if err := cmd.Start(); err != nil {
		return errors.New(fmt.Sprintf("exec error on command: \"zli %s\" %s", arg[0], err.Error())), ""
	}
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err.Error()
	}
	if err := cmd.Wait(); err != nil {
		return errors.New(fmt.Sprintf("exec error on wait: \"zli %s\" %s", arg[0], err.Error())), ""
	}
	return nil, string(opBytes)
}
