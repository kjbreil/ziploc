package extract

import (
	"github.com/kjbreil/ziploc/option"
	"testing"
)

func TestReadZip(t *testing.T) {
	_, err := ReadZip("./samples/Heinens_Lane_099.zip", "out", true, &option.Option{
		Name:     "HEINENS_LANE_COMMON",
		Priority: 30,
		IsOption: false,
		Include: map[string][]string{
			"every": []string{
				".*",
			},
		},
		BaseFolder: "out/COMMON",
		TempDir:    "build",
		ZipOut:     "zips",
	})
	if err != nil {
		t.Fatal(err)
	}

}
