package checkpoints

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-devops/release-co-pilot/execution"
	"github.com/golang-devops/release-co-pilot/logging"
	"github.com/golang-devops/release-co-pilot/util"
)

func NewStdin(description string, randomTokenLength int) execution.Task {
	return &stdin{
		description:       description,
		randomTokenLength: randomTokenLength,
	}
}

type stdin struct {
	description       string
	randomTokenLength int
}

func (s *stdin) Execute(logger logging.Logger) error {
	logger = logger.WithFields(map[string]interface{}{
		"obj-type":               fmt.Sprintf("%T", s),
		"checkpoint-description": s.description,
	})

	randomToken := util.RandomAlphaNumericString(s.randomTokenLength)
	abortKeyword := "abort"

	if _, err := os.Stdout.WriteString(fmt.Sprintf("Checkpoint %s. Please type token %s to continue or '%s' to abort.\n> ", s.description, randomToken, abortKeyword)); err != nil {
		logger.WithError(err).Error("Unable to write message to Stdout")
		return errors.New("Unable to write message to Stdout, see logs")
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			logger.WithError(err).Error("Unable to read Stdin")
			return fmt.Errorf("Failed to read Stdin, see logs")
		}
		text = strings.TrimSpace(text)

		if text == randomToken {
			return nil
		}

		if text == abortKeyword {
			logger.Error("User aborted with keyword " + abortKeyword)
			return fmt.Errorf("User aborted with keyword " + abortKeyword)
		}

		if _, err := os.Stdout.WriteString(fmt.Sprintf("Unexpected input '%s' received, try again.\n> ", text)); err != nil {
			logger.WithError(err).Error("Unable to write to Stdout (2)")
			return errors.New("Unable to write to Stdout (2), see logs")
		}
	}
}
