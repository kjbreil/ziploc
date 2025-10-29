package dss

import (
	"testing"
)

func TestRead(t *testing.T) {
	filename := "./dss.sql"

	_, err := Read(filename)
	if err != nil {
		t.Fatal(err)
	}
}
