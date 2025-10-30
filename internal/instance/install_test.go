package instance

import "testing"

func TestInstallLocation(t *testing.T) {
	type args struct {
		filename       string
		instanceFolder string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Office",
			args: args{
				filename:       "./SmsCode/Office/a_file.sqi",
				instanceFolder: "//Volumes/Storeman/",
			},
			want:    "/Volumes/Storeman/Office/a_file.sqi",
			wantErr: false,
		},
		{
			name: "SSM",
			args: args{
				filename:       "./SmsCode/SSM/a_file.sqi",
				instanceFolder: "//Volumes/Storeman/",
			},
			want:    "/Volumes/Storeman/Office/SSM/a_file.sqi",
			wantErr: false,
		},
		{
			name: "XCHDEV",
			args: args{
				filename:       "./SmsCode/XCHDEV/a_file.sqi",
				instanceFolder: "//Volumes/Storeman/",
			},
			want:    "/Volumes/Storeman/XCHDEV/a_file.sqi",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InstallLocation(tt.args.filename, tt.args.instanceFolder)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstanceFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("InstanceFolder() got = %v, want %v", got, tt.want)
			}
		})
	}
}
