package wrapper

//ResponseWrapper provides a standard structure for API responses.
type ResponseWrapper struct {
	//Status indicates the result of a request as "success" or "error"
	Status string `json:"status"`
	//Data holds the payload of the response
	Data interface{} `json:"data,omitempty"`
	//Message contains the nature of an error
	Message string `json:"message,omitempty"`
	//Count contains the number of records in the result set
	Count int `json:"count,omitempty"`
}

//SuccessWrapper is for successful requests that yield a single
//result value.
func SuccessWrapper(data interface{}) (rsp *ResponseWrapper) {
	rsp = &ResponseWrapper{
		Status: successStatus,
		Data:   data,
	}
	return
}

//CollectionWrapper is for successful reuqests that have the potential
//to yield multiple results.
func CollectionWrapper(data interface{}, count int) (rsp *ResponseWrapper) {
	rsp = &ResponseWrapper{
		Status: successStatus,
		Data:   data,
		Count:  count,
	}
	return
}

//ErrorWrapper is for requests that yield no results due to an error.
func ErrorWrapper(message string) (rsp *ResponseWrapper) {
	rsp = &ResponseWrapper{
		Status:  errorStatus,
		Message: message,
	}
	return
}
