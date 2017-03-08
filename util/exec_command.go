package util

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/golang-devops/release-co-pilot/logging"
)

func ExecCommand(logger logging.Logger, cmd *exec.Cmd, output *[]byte) error {
	out, err := cmd.CombinedOutput()
	if out != nil && output != nil {
		*output = out
	}
	if err != nil {
		outStr := ""
		if out != nil {
			outStr = string(out)
		}
		logger.WithError(err).WithFields(map[string]interface{}{
			"output": strings.Replace(strings.Replace(outStr, "\n", "\\n", -1), "\r", "", -1),
		}).Error("Command failed")
		return errors.New("Command failed")
	}
	return nil
}
