package utils

import (
	"testing"
	"time"
)

func TestChangeFileTimeAttr(t *testing.T) {
	type args struct {
		path  string
		cTime *time.Time
		aTime *time.Time
		mTime *time.Time
	}
	n := time.Now()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				path:  "E:\\Go\\out\\static\\app.css",
				cTime: &n,
				aTime: nil,
				mTime: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ChangeFileTimeAttr(tt.args.path, tt.args.cTime, tt.args.aTime, tt.args.mTime); (err != nil) != tt.wantErr {
				t.Errorf("ChangeFileTimeAttr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
