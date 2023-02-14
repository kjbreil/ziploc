package option

import "testing"

func Test_statBase(t *testing.T) {
	type args struct {
		baseWithRoot string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Network Share",
			args: args{
				baseWithRoot: "\\\\sco\\Lane\\Office",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := statBase(tt.args.baseWithRoot); (err != nil) != tt.wantErr {
				t.Errorf("statBase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
