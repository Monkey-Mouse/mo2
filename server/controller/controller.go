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
	Time   time.Time `json:"time"`
	Reason string    `json:"reason"`
}

func setResponseError(err error) (r ResponseError) {
	r.Reason = err.Error()
	r.Time = time.Now()
	return
}
func setResponseReason(err string) (r ResponseError) {
	r.Reason = err
	r.Time = time.Now()
	return
}
