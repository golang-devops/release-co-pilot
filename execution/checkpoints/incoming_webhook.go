package checkpoints

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/golang-devops/release-co-pilot/execution"
	"github.com/golang-devops/release-co-pilot/logging"
	"github.com/golang-devops/release-co-pilot/util"
)

func NewIncomingWebhook(listenHTTPAddress string, randomTokenLength int) execution.Task {
	return &incomingWebhook{
		listenHTTPAddress: listenHTTPAddress,
		randomTokenLength: randomTokenLength,
	}
}

type incomingWebhook struct {
	listenHTTPAddress string
	randomTokenLength int
}

func (i *incomingWebhook) runServer(logger logging.Logger, randomToken, abortKeyword string) error {
	var wg sync.WaitGroup
	wg.Add(1)
	var asyncError error

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.EqualFold(r.Method, "POST") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Only POST method currently allowed"))
			return
		}

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.WithError(err).WithFields(map[string]interface{}{
				"remote-address":  r.RemoteAddr,
				"request-headers": r.Header,
			}).Warn(fmt.Sprintf("Failed to read POST body"))
			return
		}

		bodyStr := strings.TrimSpace(string(body))
		if bodyStr == randomToken {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Thank you, token is accepted"))
			wg.Done()
		}

		if bodyStr == abortKeyword {
			asyncError = errors.New("User aborted with keyword " + abortKeyword)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Abort keyword received"))
			wg.Done()
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unexpected input '%s' received, try again.", bodyStr)))
	})

	go func() {
		defer wg.Done()

		logger.Info(fmt.Sprintf("Now starting incoming webhook on '%s' to wait for token %s to continue or '%s' to abort", i.listenHTTPAddress, randomToken, abortKeyword))
		if err := http.ListenAndServe(i.listenHTTPAddress, nil); err != nil {
			logger.WithError(err).Error("Unable to start HTTP server")
			asyncError = errors.New("Unable to start HTTP server, see logs")
			return
		}
	}()

	wg.Wait()

	if asyncError != nil {
		return asyncError
	}

	return nil
}

func (i *incomingWebhook) Execute(logger logging.Logger) error {
	logger = logger.WithFields(map[string]interface{}{
		"obj-type":            fmt.Sprintf("%T", i),
		"listen-http-address": i.listenHTTPAddress,
	})

	randomToken := util.RandomAlphaNumericString(i.randomTokenLength)
	abortKeyword := "abort"

	if err := i.runServer(logger, randomToken, abortKeyword); err != nil {
		return err
	}

	return nil
}
