package util

import (
	"errors"
	"os/exec"

	"github.com/golang-devops/release-co-pilot/logging"
)

func ExecCommand(logger logging.Logger, cmd *exec.Cmd) error {
	out, err := cmd.CombinedOutput()
	if err != nil {
		outStr := ""
		if out != nil {
			outStr = string(out)
		}
		logger.WithError(err).WithFields(map[string]interface{}{
			"output": outStr,
		}).Error("Command failed")
		return errors.New("Command failed")
	}
	return nil
}
