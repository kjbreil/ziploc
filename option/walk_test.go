package option

import "testing"

func Test_checkIncluded(t *testing.T) {
	included := map[string][]string{
		"every": []string{},
	}
	type args struct {
		include map[string][]string
	}
	tests := []struct {
		name    string
		current string
		want    bool
		want1   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := checkIncluded(tt.current, included)
			if got != tt.want {
				t.Errorf("checkIncluded() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("checkIncluded() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
