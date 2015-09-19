package params

import (
	"net/url"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

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

//ExtractRequestParams initializes the RequestParams object.
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
	p.L = limitDefault
	return
}

//Selector returns a mongodb bson.M object containing the query parameters
//supplied in the HTTP request.
func (p *RequestParams) Selector() bson.M {
	return p.Q
}

//Scope returns a mongodb bson.M object containing the set of fields to
//be returned from mongodb.  Empty scope objects will return all fields
//from the database.
func (p *RequestParams) Scope() bson.M {
	return p.S
}

//Limit returns an integer value indicating the size of the result set
//returned from the database.
func (p *RequestParams) Limit() int {
	return p.L
}

//Offset returns an integer value indicating the number of records to skip
//when returning a result set.
func (p *RequestParams) Offset() int {
	return p.F
}

func (p *RequestParams) parseSelector() {
	for k, v := range p.RawQuery {
		if k == scopeKeyword || k == limitKeyword || k == offsetKeyword {
			continue
		} else {
			p.Q[k] = v[0]
		}
	}
	return
}

func (p *RequestParams) parseScope() {
	s := p.RawQuery.Get(scopeKeyword)
	if len(s) > 0 {
		s1 := strings.Split(s, ",")
		for _, v := range s1 {
			p.S[v] = 1
		}
	}
	return
}

func (p *RequestParams) parseLimit() {
	s := p.RawQuery.Get(limitKeyword)
	if len(s) > 0 {
		l, err := strconv.Atoi(s)
		if err == nil {
			p.L = l
		}
	}
	return
}

func (p *RequestParams) parseOffset() {
	s := p.RawQuery.Get(offsetKeyword)
	if len(s) > 0 {
		o, err := strconv.Atoi(s)
		if err == nil {
			p.F = o
		}
	}
	return
}
