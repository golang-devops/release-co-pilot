package logging

import (
	"os"

	apex "github.com/francoishill/log"

	"github.com/go-zero-boilerplate/extended-apex-logger/logging"
	"github.com/go-zero-boilerplate/extended-apex-logger/logging/text_handler"
)

var logger Logger

type Logger interface {
	logging.Logger
}

func DefaultLogger() Logger {
	return logger
}

func init() {
	level := apex.DebugLevel
	loggerFields := apex.Fields{}
	// loggerFields["SourceVersion"] = constants.SourceVersion
	//TODO: Use an Environment variable to define Git sha1 (see VERSION label in Dockerfile)
	/*if len(strings.TrimSpace(GitSha1)) > 0 {
		//cater for scenario where git sha is not available
		loggerFields["git_sha1"] = GitSha1[:8]
	}*/
	apexEntry := apex.WithFields(loggerFields)

	logHandler := text_handler.New(os.Stdout, os.Stderr, text_handler.DefaultTimeStampFormat, text_handler.DefaultMessageWidth)
	exitOnEmergency := true
	logger = logging.NewApexLogger(level, logHandler, apexEntry, exitOnEmergency)
}
