package model

import (
	"fmt"
)

type PostTxBody struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Emsg     string `json:"emsg"`
}

func (p PostTxBody) ToDocumentEntry() DocumentEntry {
	return DocumentEntry{
		Answer: p.Answer,
		Emsg:   p.Emsg,
	}
}

type PostTxResponse struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Verify   string `json:"verify"`
}

func (p PostTxResponse) toURL() string {
	return fmt.Sprintf("/v1/rx?i=%s&q=%s&v=%s", p.ID, p.Question, p.Verify)
}

func (p PostTxResponse) Response() string {
	return fmt.Sprintf("{\"path\": \"%s\"}", p.toURL())
}
