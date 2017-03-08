package tasks

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/golang-devops/release-co-pilot/execution"
	"github.com/golang-devops/release-co-pilot/logging"
	"github.com/golang-devops/release-co-pilot/util"
)

func NewGitClone(localCloneDir, remoteURI string) execution.Task {
	return &gitClone{
		localCloneDir: localCloneDir,
		remoteURI:     remoteURI,
	}
}

type gitClone struct {
	localCloneDir string
	remoteURI     string
}

func (g *gitClone) Execute(logger logging.Logger) error {
	logger = logger.WithFields(map[string]interface{}{
		"obj-type":        fmt.Sprintf("%T", g),
		"remote-uri":      g.remoteURI,
		"local-clone-dir": g.localCloneDir,
	})

	cmd := exec.Command("git", "clone", g.remoteURI, g.localCloneDir)
	if err := util.ExecCommand(logger, cmd); err != nil {
		return errors.New("Git clone failed, details written to logger")
	}

	return nil
}
