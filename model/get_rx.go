package model

type GetRxParams struct {
	ID     string `query:"i" json:"id"`
	Verify string `query:"v" json:"verify"`
	Answer string `query:"a" json:"answer"`
}
