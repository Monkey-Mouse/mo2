package controller

import "time"

// Controller example
type Controller struct {
}

// NewController example
func NewController() *Controller {
	return &Controller{}
}

// Message example
type Message struct {
	Message string `json:"message" example:"message"`
}

type ResponseError struct {
	time   time.Time `json:"time"`
	reason string    `json:"reason"`
}

func setResponseError(err error) (r ResponseError) {
	r.reason = err.Error()
	r.time = time.Now()
	return
}
