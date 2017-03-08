package tasks

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/golang-devops/release-co-pilot/execution"
	"github.com/golang-devops/release-co-pilot/logging"
	"github.com/golang-devops/release-co-pilot/util"
)

func NewPSCPUploadDir(localDir string, remoteHost, remoteUser string, remotePort int, remoteParentDir string) execution.Task {
	return &pscpUploadDir{
		localDir:        localDir,
		remoteHost:      remoteHost,
		remoteUser:      remoteUser,
		remotePort:      remotePort,
		remoteParentDir: remoteParentDir,
	}
}

type pscpUploadDir struct {
	localDir        string
	remoteHost      string
	remoteUser      string
	remotePort      int
	remoteParentDir string
}

func (p *pscpUploadDir) Execute(logger logging.Logger) error {
	logger = logger.WithFields(map[string]interface{}{
		"obj-type":          fmt.Sprintf("%T", p),
		"local-dir":         p.localDir,
		"remote-host":       p.remoteHost,
		"remote-user":       p.remoteUser,
		"remote-port":       p.remotePort,
		"remote-parent-dir": p.remoteParentDir,
	})

	pscpArgs := []string{
		"-r",
		"-C",
		"-P",
		fmt.Sprintf("%d", p.remotePort),
		"-agent",
		p.localDir,
		fmt.Sprintf("%s@%s:%s", p.remoteUser, p.remoteHost, p.remoteParentDir),
	}
	cmd := exec.Command("pscp", pscpArgs...)
	if err := util.ExecCommand(logger, cmd); err != nil {
		return errors.New("PSCP upload failed, details written to logger")
	}

	return nil
}
