package ziploc

import (
	"testing"
)

func TestOption_excluded(t *testing.T) {
	o := &Option{
		Exclude: []string{
			"TEMP",
			"SAMPLES",
			"OPTIONS",
			".DS_Store",
		},
	}

	switch {
	case o.excluded("TEMP"):
	case o.excluded("SAMPLES"):
	case o.excluded(".DS_Store"):
	case !o.excluded("NOTEXCLUDED"):
	default:
		t.Fail()

	}

}
