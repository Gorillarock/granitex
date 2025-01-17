package model

import "errors"

var (
	ERROR_UNKNOWN = errors.New("Unknown Service Error")
)

type GetRxParams struct {
	ID     string `query:"i" json:"id" bson:"_id"`
	Verify string `query:"v" json:"verify" bson:"verify"`
	Answer string `query:"a" json:"answer" bson:"answer"`
}

type ResponsePayloadRxHandler struct {
	Emsg    string `json:"emsg,omitempty"`
	Status  int    `json:"-"`
	Error   error  `json:"error,omitempty"`
	Deleted bool   `json:"-"`
}
