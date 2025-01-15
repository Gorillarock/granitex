package model

import (
	"fmt"
)

type DocumentEntry struct {
	ID     string `json:"id"`
	Verify string `json:"verify"`
	Answer string `json:"answer"`
	Emsg   string `json:"emsg"`
	Try    int    `json:"try"`
}

func (d DocumentEntry) EmsgResponse() string {
	return fmt.Sprintf("{\"eMsg\": \"%s\"", d.Emsg)
}

/*
verify that the ID and Verify parameters match for the document.
Ensures that the attempt at the document is allowed.
*/
func (d DocumentEntry) VerifyInput(input GetRxParams) bool {
	if d.ID != input.ID {
		return false
	}
	if d.Verify != input.Verify {
		return false
	}
	return true

}

// Returns match == true if the answer is correct and mustDelete == true if the document must be deleted.  Document should be deleted if the answer is correct or if the number of tries exceeds 3.
func (d DocumentEntry) CheckAnswer(answer string) (match bool, mustDelete bool) {
	// Increment try
	d.Try++
	if d.Try >= 3 {
		// Delete document from DB
		// TODO: Implement
		mustDelete = true
		return
	}

	match = d.Answer == answer
	if match {
		// Delete document from DB
		mustDelete = true
	}
	return
}
