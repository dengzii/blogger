package repo

import (
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		url         string
		accessToken string
		gitDir      string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "\\",
			args: args{
				url:         "https://github.com/dengzii/RespberryPi",
				accessToken: "",
				gitDir:      "../source",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := New(tt.args.url, tt.args.accessToken, tt.args.gitDir)
			rep.Remove()
			rep.Clone()
		})
	}
}
