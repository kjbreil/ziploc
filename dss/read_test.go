package dss

import (
	"testing"
)

func TestRead(t *testing.T) {
	filename := "./dss.sql"

	Read(filename)
	t.Fail()
}
