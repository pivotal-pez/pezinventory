package pezinventory

import (
	"net/url"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

const (
	successStatus = "success"
	errorStatus   = "error"

	limitDefault  = 10
	limitKeyword  = "limit"
	scopeKeyword  = "scope"
	offsetKeyword = "offset"
)

//ResponseMessage structures output into a standard format.
type ResponseMessage struct {
	//Status indicates the result of a request as "success" or "error"
	Status string `json:"status"`
	//Data holds the payload of the response
	Data interface{} `json:"data,omitempty"`
	//Message contains the nature of an error
	Message string `json:"message,omitempty"`
	//Count contains the number of records in the result set
	Count int `json:"count,omitempty"`
	//Prev provides the URL to the previous result set
	Prev string `json:"prev_url,omitempty"`
	//Next provides the URL to the next result set
	Next string `json:"next_url,omitempty"`
}

func successMessage(data interface{}) (rsp ResponseMessage) {
	rsp = ResponseMessage{
		Status: successStatus,
		Data:   data,
	}
	return
}

func collectionMessage(data interface{}, count int, params *RequestParams) (rsp ResponseMessage) {
	rsp = ResponseMessage{
		Status: successStatus,
		Data:   data,
		Count:  count,
	}
	return
}

func errorMessage(message string) (rsp ResponseMessage) {
	rsp = ResponseMessage{
		Status:  errorStatus,
		Message: message,
	}
	return
}

//RequestParams holds state parsed from a given HTTP request
type RequestParams struct {
	//RawQuery contains the raw query string Values object
	RawQuery url.Values `json:"raw_query"`
	//Selector holds the query parameters specified in the request.
	//Defaults to bson.M{}.
	Selector bson.M `json:"selector"`
	//Scope specifies the fields to be included in the result set.
	//Defaults to bson.M{}.  A nil scope will return the entire dataset.
	Scope bson.M `json:"scope"`
	//Limit specifies the maximum number of records to be retrieved
	//for a given request.  Limit defaults to 10.
	Limit int `json:"limit"`
	//Offset specifies the number of records to skip in the result set.
	//This is useful for paging through large result sets.
	//Offset defaults to 0.
	Offset int `json:"offset"`
}

func ExtractRequestParams(query url.Values) (p *RequestParams) {
	p = newRequestParams(query)
	p.parseSelector()
	p.parseLimit()
	p.parseOffset()
	p.parseScope()
	return
}

func newRequestParams(raw url.Values) (p *RequestParams) {
	p = new(RequestParams)
	p.RawQuery = raw
	p.Selector = bson.M{}
	p.Scope = bson.M{}
	p.Limit = limitDefault
	return
}

func (p *RequestParams) parseSelector() {
	for k, v := range p.RawQuery {
		if k == scopeKeyword || k == limitKeyword || k == offsetKeyword {
			continue
		} else {
			p.Selector[k] = v[0]
		}
	}
	return
}

func (p *RequestParams) parseScope() {
	s := p.RawQuery.Get(scopeKeyword)
	if len(s) > 0 {
		s1 := strings.Split(s, ",")
		for _, v := range s1 {
			p.Scope[v] = 1
		}
	}
	return
}

func (p *RequestParams) parseLimit() {
	s := p.RawQuery.Get(limitKeyword)
	if len(s) > 0 {
		l, err := strconv.Atoi(s)
		if err == nil {
			p.Limit = l
		}
	}
	return
}

func (p *RequestParams) parseOffset() {
	s := p.RawQuery.Get(offsetKeyword)
	if len(s) > 0 {
		o, err := strconv.Atoi(s)
		if err == nil {
			p.Offset = o
		}
	}
	return
}
