package tropo


import (
	"time"
)


type SessionWrapper struct {
	session Session `json:"session"`
}

type Session struct {
	ID string `json:"id"`
	AccountID string `json:"accountId"`
	Timestamp time.Time `json:"timestamp"`
	UserType string `json:"userType"`
	Initialtext interface{} `json:"initialText"`
	CallID string `json:"callId"`
	To struct {
		ID string `json:"id"`
		Name interface{} `json:"name"`
		Channel string `json:"channel"`
		Network string `json:"network"`
	} `json:"to"`
	From struct {
		ID string `json:"id"`
		Name interface{} `json:"name"`
		Channel string `json:"channel"`
		Network string `json:"network"`
	} `json:"from"`
	Headers struct {
		MaxForwards string `json:"Max-Forwards"`
		Xsid string `json:"x-sid"`
		RecordRoute string `json:"Record-Route"`
		ContentLength string `json:"Content-Length"`
		Contact string `json:"Contact"`
		To string `json:"To"`
		Cseq string `json:"CSeq"`
		UserAgent string `json:"User-Agent"`
		Via string `json:"Via"`
		CallID string `json:"Call-ID"`
		ContentType string `json:"Content-Type"`
		From string `json:"From"`
	} `json:"headers"`
}
