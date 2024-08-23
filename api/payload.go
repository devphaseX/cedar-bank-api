package api

type FailedResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

func errorResponse(err error) FailedResponse {
	return FailedResponse{Status: false, Error: err.Error()}
}

type SuccessResponse struct {
	Status  bool   `json:"status"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func sucessResponse(data any, message ...string) SuccessResponse {
	var msg string
	if len(message) != 0 {
		msg = message[0]
	}

	return SuccessResponse{
		Status:  true,
		Data:    data,
		Message: msg,
	}
}
