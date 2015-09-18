package paged

import (
	"net/url"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

//Pager types provide data necessary for mongodb queries to support
//filtering, scoping, and paging.
type Pager interface {
	Selector() bson.M
	Scope() bson.M
	Limit() int
	Offset() int
}

//RequestParams holds state parsed from a given HTTP request
type RequestParams struct {
	//RawQuery contains the raw query string Values object
	RawQuery url.Values `json:"raw_query"`
	//Q (selector) holds the query parameters specified in the request.
	//Defaults to bson.M{}.
	Q bson.M `json:"selector"`
	//S (scope) specifies the fields to be included in the result set.
	//Defaults to bson.M{}.  A nil scope will return the entire dataset.
	S bson.M `json:"scope"`
	//L (limit) specifies the maximum number of records to be retrieved
	//for a given request.  Limit defaults to 10.
	L int `json:"limit"`
	//F (offset) specifies the number of records to skip in the result set.
	//This is useful for paging through large result sets.
	//F defaults to 0.
	F int `json:"offset"`
}

//ExtractRequestParameters initializes the RequestParams object.
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
	p.Q = bson.M{}
	p.S = bson.M{}
	p.L = LimitDefault
	return
}

func (p *RequestParams) Selector() bson.M {
	return p.Q
}

func (p *RequestParams) Scope() bson.M {
	return p.S
}

func (p *RequestParams) Limit() int {
	return p.L
}

func (p *RequestParams) Offset() int {
	return p.F
}

func (p *RequestParams) parseSelector() {
	for k, v := range p.RawQuery {
		if k == ScopeKeyword || k == LimitKeyword || k == OffsetKeyword {
			continue
		} else {
			p.Q[k] = v[0]
		}
	}
	return
}

func (p *RequestParams) parseScope() {
	s := p.RawQuery.Get(ScopeKeyword)
	if len(s) > 0 {
		s1 := strings.Split(s, ",")
		for _, v := range s1 {
			p.S[v] = 1
		}
	}
	return
}

func (p *RequestParams) parseLimit() {
	s := p.RawQuery.Get(LimitKeyword)
	if len(s) > 0 {
		l, err := strconv.Atoi(s)
		if err == nil {
			p.L = l
		}
	}
	return
}

func (p *RequestParams) parseOffset() {
	s := p.RawQuery.Get(OffsetKeyword)
	if len(s) > 0 {
		o, err := strconv.Atoi(s)
		if err == nil {
			p.F = o
		}
	}
	return
}
