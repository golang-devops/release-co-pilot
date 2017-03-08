package execution

import (
	"github.com/golang-devops/release-co-pilot/logging"
)

type Step interface {
	Execute(logger logging.Logger) error
}
