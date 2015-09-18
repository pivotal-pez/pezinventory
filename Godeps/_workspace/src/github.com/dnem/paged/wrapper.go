package paged

//ResponseMessage structures output into a standard format.
type ResponseWrapper struct {
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

func SuccessWrapper(data interface{}) (rsp ResponseWrapper) {
	rsp = ResponseWrapper{
		Status: SuccessStatus,
		Data:   data,
	}
	return
}

func CollectionWrapper(data interface{}, count int, params *RequestParams) (rsp ResponseWrapper) {
	rsp = ResponseWrapper{
		Status: SuccessStatus,
		Data:   data,
		Count:  count,
	}
	return
}

func ErrorWrapper(message string) (rsp ResponseWrapper) {
	rsp = ResponseWrapper{
		Status:  ErrorStatus,
		Message: message,
	}
	return
}
