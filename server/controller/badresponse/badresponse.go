package badresponse

import "time"

type ResponseError struct {
	Time   time.Time `json:"time"`
	Reason string    `json:"reason"`
}

func SetResponseError(err error) (r ResponseError) {
	r.Reason = err.Error()
	r.Time = time.Now()
	return
}
func SetResponseReason(err string) (r ResponseError) {
	r.Reason = err
	r.Time = time.Now()
	return
}
