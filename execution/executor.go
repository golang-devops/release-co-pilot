package execution

import (
	"github.com/golang-devops/release-co-pilot/logging"
)

type Executor interface {
	Execute(logger logging.Logger, steps []Step) error
}

func NewExecutor() Executor {
	return &executor{}
}

type executor struct{}

func (e *executor) Execute(logger logging.Logger, steps []Step) error {
	for index, step := range steps {
		if err := step.Execute(logger); err != nil {
			logger.WithError(err).WithFields(map[string]interface{}{
				"step-index": index,
			}).Error("Failed to execute step")
			return err
		}
	}
	return nil
}
