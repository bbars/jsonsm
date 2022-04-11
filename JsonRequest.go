package jsonsm

import (
	"encoding/json"
	"net/http"
	"regexp"
)


type JsonRequest struct {
	*http.Request
	MatchedPattern *regexp.Regexp
	Matches []string
}


func (this *JsonRequest) Payload(payload interface{}) error {
	// TODO: check Content-Type header
	return json.NewDecoder(this.Request.Body).Decode(payload)
}
