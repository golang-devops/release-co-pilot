package checkpoints

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

	resultChan              chan<- *incomingWebhookResult
	randomToken             string
	abortKeyword            string
	handlePostRequestLogger logging.Logger
}

type incomingWebhookResult struct {
	Error        error
	ServerClosed bool
}

func (i *incomingWebhook) handlePostRequest(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, "POST") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only POST method currently allowed"))
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		i.handlePostRequestLogger.WithError(err).WithFields(map[string]interface{}{
			"remote-address":  r.RemoteAddr,
			"request-headers": r.Header,
		}).Warn(fmt.Sprintf("Failed to read POST body"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Body is empty, try again.")))
		return
	}

	bodyStr := strings.TrimSpace(string(body))
	if bodyStr == i.randomToken {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Thank you, token is accepted"))
		i.resultChan <- &incomingWebhookResult{}
		return
	}

	if bodyStr == i.abortKeyword {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Abort keyword received"))
		i.resultChan <- &incomingWebhookResult{
			Error: errors.New("User aborted with keyword " + i.abortKeyword),
		}
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("Unexpected input '%s' received, try again.", bodyStr)))
}

func (i *incomingWebhook) runServer(logger logging.Logger, server *http.Server) {
	logger.Info(fmt.Sprintf("Now starting incoming webhook on '%s' to wait for token %s to continue or '%s' to abort", i.listenHTTPAddress, i.randomToken, i.abortKeyword))
	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			logger.WithError(err).Error("Unable to start HTTP server")
			i.resultChan <- &incomingWebhookResult{
				Error:        errors.New("Server listen error: " + err.Error()),
				ServerClosed: true,
			}
		}
		logger.Debug("HTTP server connection closed")
		return
	}
}

func (i *incomingWebhook) runAndWait(logger logging.Logger) error {
	resultChan := make(chan *incomingWebhookResult)
	i.resultChan = resultChan

	mux := http.NewServeMux()
	i.handlePostRequestLogger = logger
	mux.HandleFunc("/", i.handlePostRequest)

	server := &http.Server{
		Addr:    i.listenHTTPAddress,
		Handler: mux,
	}

	go i.runServer(logger, server)

	result := <-resultChan

	if !result.ServerClosed {
		if err := server.Close(); err != nil {
			logger.WithError(err).Error("Unable to shutdown webhook server")
		}
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (i *incomingWebhook) Execute(logger logging.Logger) error {
	logger = logger.WithFields(map[string]interface{}{
		"obj-type":            fmt.Sprintf("%T", i),
		"listen-http-address": i.listenHTTPAddress,
	})

	i.randomToken = util.RandomAlphaNumericString(i.randomTokenLength)
	i.abortKeyword = "abort"

	if err := i.runAndWait(logger); err != nil {
		return err
	}

	return nil
}
