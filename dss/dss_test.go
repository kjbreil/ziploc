package dss

import (
	"testing"
)

func Test_Destination(t *testing.T) {
	type args struct {
		fp string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				fp: "/storeman/office/system.ini",
			},
			want: "OFFICE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Destination(tt.args.fp); got != tt.want {
				t.Errorf("destination() = %v, want %v", got, tt.want)
			}
		})
	}
}
