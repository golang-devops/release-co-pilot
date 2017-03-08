package tasks

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/golang-devops/release-co-pilot/execution"
	"github.com/golang-devops/release-co-pilot/logging"
	"github.com/golang-devops/release-co-pilot/util"
)

func NewPlinkRemoteExec(outDest *[]byte, remoteHost, remoteUser string, remotePort int, remoteCommandArgs ...string) execution.Task {
	if len(remoteCommandArgs) == 0 {
		panic("Arg remoteCommandArgs should have at least one element when calling NewPlinkRemoteExec")
	}
	return &plinkRemoteExec{
		outDest:           outDest,
		remoteHost:        remoteHost,
		remoteUser:        remoteUser,
		remotePort:        remotePort,
		remoteCommandArgs: remoteCommandArgs,
	}
}

type plinkRemoteExec struct {
	outDest           *[]byte
	remoteHost        string
	remoteUser        string
	remotePort        int
	remoteCommandArgs []string
}

func (p *plinkRemoteExec) Execute(logger logging.Logger) error {
	logger = logger.WithFields(map[string]interface{}{
		"obj-type":            fmt.Sprintf("%T", p),
		"remote-host":         p.remoteHost,
		"remote-user":         p.remoteUser,
		"remote-port":         p.remotePort,
		"remote-command-args": p.remoteCommandArgs,
	})

	plinkArgs := []string{
		"-C",
		"-P",
		fmt.Sprintf("%d", p.remotePort),
		"-agent",
		"-ssh",
		fmt.Sprintf("%s@%s", p.remoteUser, p.remoteHost),
	}
	plinkArgs = append(plinkArgs, p.remoteCommandArgs...)

	chmodCmd := exec.Command("plink", plinkArgs...)
	if err := util.ExecCommand(logger, chmodCmd, p.outDest); err != nil {
		return errors.New("Plink chmod failed, details written to logger")
	}

	return nil
}
