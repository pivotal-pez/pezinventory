package pezinventory

const (
	successStatus = "success"
	errorStatus   = "error"
	failStatus    = "fail"
)

//ResponseMessage structures output into a standard format.
type ResponseMessage struct {
	//Status returns a string indicating [success|error|fail]
	Status string `json:"status"`
	//Data holds the payload of the response
	Data interface{} `json:"data,omitempty"`
	//Message contains the nature of an error
	Message string `json:"message,omitempty"`
	//Meta contains information about the data and the current request
	Meta map[string]interface{} `json:"_metaData,omitempty"`
	//Links contains [prev|next] links for paginated responses
	Links map[string]interface{} `json:"_links,omitempty"`
}

func successMessage(data interface{}) (rsp ResponseMessage) {
	rsp = ResponseMessage{
		Status: successStatus,
		Data:   data,
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

func failureMessage(data interface{}) (rsp ResponseMessage) {
	rsp = ResponseMessage{
		Status: failStatus,
		Data:   data,
	}
	return
}
