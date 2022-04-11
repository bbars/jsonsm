package jsonsm

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"time"
)


type Handler interface {
	HandleJson(*JsonRequest) (interface{}, error)
}


type funcHandler struct {
	methodPathRegexp *regexp.Regexp
	fn func(*JsonRequest) (interface{}, error)
}

func (this funcHandler) HandleJson(req *JsonRequest) (interface{}, error) {
	return this.fn(req)
}

type JsonServeMux struct {
	*http.ServeMux
	handlers []funcHandler
	logger *log.Logger
}


func NewJsonServeMux(serveMux *http.ServeMux, logger *log.Logger) *JsonServeMux {
	this := &JsonServeMux{
		ServeMux: serveMux,
		handlers: make([]funcHandler, 0),
		logger: logger,
	}
	serveMux.Handle("/", this)
	return this
}

func (this *JsonServeMux) log(v ...interface{}) {
	if this.logger == nil {
		return
	}
	this.logger.Println(v...)
}

func (this *JsonServeMux) HandleFunc(methodPathRegexp *regexp.Regexp, fn func(*JsonRequest) (interface{}, error)) {
	this.handlers = append(this.handlers, funcHandler{
		methodPathRegexp: methodPathRegexp,
		fn: fn,
	})
}

func (this *JsonServeMux) executeHandler(jsonRequest *JsonRequest, handler Handler) (int, interface{}) {
	responseData, err := handler.HandleJson(jsonRequest)
	
	if err == nil {
		return 200, responseData
	}
	
	var err2 *Error
	if err3, ok := err.(*Error); ok {
		err2 = err3
	} else {
		err2 = WrapError(err, nil)
	}
	return err2.HttpCode(), err2
}

func (this *JsonServeMux) ServeHTTP(resposeWriter http.ResponseWriter, req *http.Request) {
	var matches []string
	reqMethodPath := req.Method + " " + req.URL.Path
	this.log("<", reqMethodPath)
	for _, handler := range this.handlers {
		matches = handler.methodPathRegexp.FindStringSubmatch(reqMethodPath)
		if len(matches) < 1 {
			continue
		}
		jsonRequest := &JsonRequest{
			Request: req,
			MatchedPattern: handler.methodPathRegexp,
			Matches: matches,
		}
		
		timeStart := time.Now()
		httpCode, responseData := this.executeHandler(jsonRequest, handler)
		responseJson, err := json.Marshal(responseData)
		if err != nil {
			this.log("!", "Unable to encode JSON", err)
		}
		this.log(">", time.Now().Sub(timeStart), httpCode, string(responseJson))
		
		resposeWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		resposeWriter.WriteHeader(httpCode)
		resposeWriter.Write(responseJson)
		
		return
	}
	this.log("!", "Unable to find handler")
}
