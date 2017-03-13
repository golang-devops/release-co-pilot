package tasks

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/golang-devops/release-co-pilot/execution"
	"github.com/golang-devops/release-co-pilot/logging"
	"github.com/golang-devops/release-co-pilot/util"
)

func NewPSCPUploadDir(outDest *[]byte, localDir string, remoteHost, remoteUser string, remotePort int, remoteParentDir string) execution.Task {
	return &pscpUpload{
		outDest:    outDest,
		flags:      []string{"-r", "-C"},
		localDir:   localDir,
		remoteHost: remoteHost,
		remoteUser: remoteUser,
		remotePort: remotePort,
		remotePath: remoteParentDir,
	}
}

func NewPSCPUploadFile(outDest *[]byte, localDir string, remoteHost, remoteUser string, remotePort int, remoteParentDir string) execution.Task {
	return &pscpUpload{
		outDest:    outDest,
		flags:      []string{"-C"},
		localDir:   localDir,
		remoteHost: remoteHost,
		remoteUser: remoteUser,
		remotePort: remotePort,
		remotePath: remoteParentDir,
	}
}

type pscpUpload struct {
	flags      []string
	outDest    *[]byte
	localDir   string
	remoteHost string
	remoteUser string
	remotePort int
	remotePath string
}

func (p *pscpUpload) Execute(logger logging.Logger) error {
	logger = logger.WithFields(map[string]interface{}{
		"obj-type":    fmt.Sprintf("%T", p),
		"local-dir":   p.localDir,
		"remote-host": p.remoteHost,
		"remote-user": p.remoteUser,
		"remote-port": p.remotePort,
		"remote-path": p.remotePath,
	})

	pscpArgs := []string{}
	pscpArgs = append(pscpArgs, p.flags...)
	pscpArgs = append(pscpArgs, []string{
		"-P",
		fmt.Sprintf("%d", p.remotePort),
		"-agent",
		p.localDir,
		fmt.Sprintf("%s@%s:%s", p.remoteUser, p.remoteHost, p.remotePath),
	}...)
	cmd := exec.Command("pscp", pscpArgs...)
	if err := util.ExecCommand(logger, cmd, p.outDest); err != nil {
		return errors.New("PSCP upload failed, details written to logger")
	}

	return nil
}
