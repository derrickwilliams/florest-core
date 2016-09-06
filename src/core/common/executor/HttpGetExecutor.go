package executor

import (
	"fmt"
	"github.com/jabong/florest-core/src/common/constants"
	"github.com/jabong/florest-core/src/common/logger"
	"github.com/jabong/florest-core/src/common/profiler"
	"github.com/jabong/florest-core/src/common/utils/http"
	workflow "github.com/jabong/florest-core/src/core/common/orchestrator"
	"time"
)

type HttpGetExecutor struct {
	id        string
	ioDatakey HttpGetKey
}

type HttpGetKey struct {
	keyWithUrl      string
	keyWithHeaders  string
	keyWithResponse string
	keyWithTimeout  string
}

func (n *HttpGetExecutor) SetID(id string) {
	n.id = id
}

func (n HttpGetExecutor) GetID() (id string, err error) {
	return n.id, nil
}

func (n HttpGetExecutor) Name() string {
	return "HttpGetExecutor"
}

func (n *HttpGetExecutor) SetKey(key HttpGetKey) {
	n.ioDatakey = key
}

func (n HttpGetExecutor) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	prof := profiler.NewProfiler()
	profiler.StartProfile(prof, "HttpGetExecutor")
	defer func() {
		profiler.EndProfile(prof, "HttpGetExecutor")
	}()
	rc, _ := io.ExecContext.Get(constants.REQUEST_CONTEXT)

	//get url
	url, err := n.getUrl(io)
	if err != nil {
		return io, err
	}
	//get headers
	headers := n.getHeaders(io)

	timeout, err := n.getTimeout(io)
	if err != nil {
		return io, err
	}

	//log transaction id
	tid, _ := headers[constants.TRANSACTION_ID]
	logger.Info(fmt.Sprintf("tid =%+v", tid), rc)

	response, rerr := http.HttpGet(url, headers, time.Duration(timeout*time.Millisecond))
	if rerr != nil {
		return io, rerr
	}
	//set response
	io.ExecContext.Set(n.ioDatakey.keyWithResponse, response)
	return io, nil
}

func (n *HttpGetExecutor) getUrl(io workflow.WorkFlowData) (string, *constants.AppError) {
	urlObj, errUrl := io.ExecContext.Get(n.ioDatakey.keyWithUrl)
	if errUrl != nil {
		return "", &constants.AppError{Code: constants.IncorrectDataErrorCode, Message: errUrl.Error()}
	}
	url, isType := urlObj.(string)
	if !isType {
		return "", &constants.AppError{Code: constants.IncorrectDataErrorCode, Message: "Url is not a string"}
	}
	return url, nil
}

func (n *HttpGetExecutor) getHeaders(io workflow.WorkFlowData) map[string]string {
	//get headers
	rcObj, _ := io.ExecContext.Get(constants.REQUEST_CONTEXT)
	rc, _ := rcObj.(http.RequestContext)
	reqHeaders := http.GetHTTPHeaders(&rc)

	//get headers passed to executor
	headers, _ := io.ExecContext.Get(n.ioDatakey.keyWithHeaders)
	headerMap, _ := headers.(map[string]string)
	for k, v := range headerMap {
		reqHeaders[k] = v
	}
	return reqHeaders
}

func (n *HttpGetExecutor) getTimeout(io workflow.WorkFlowData) (time.Duration, *constants.AppError) {
	timeoutObj, errUrl := io.ExecContext.Get(n.ioDatakey.keyWithTimeout)
	if errUrl != nil {
		return time.Duration(0), &constants.AppError{Code: constants.IncorrectDataErrorCode, Message: errUrl.Error()}
	}
	timeout, isInt64 := timeoutObj.(int)
	if !isInt64 {
		return time.Duration(0), &constants.AppError{Code: constants.IncorrectDataErrorCode, Message: fmt.Sprintf("timeout is not an int64. value=%+v", timeoutObj)}
	}

	//zero is not valid timeout value
	if timeout == 0 {
		return time.Duration(0), &constants.AppError{Code: constants.IncorrectDataErrorCode, Message: fmt.Sprintf("invalid timeout of 0 ")}
	}
	return time.Duration(timeout), nil
}
