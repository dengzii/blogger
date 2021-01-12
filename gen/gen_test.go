package gen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFrom(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "\\",
			args: args{
				dir: "..\\sample_repo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := From(tt.args.dir, &RenderConfig{
				OutputDir:   "..\\out",
				TemplateDir: "..\\template",
			})

			assert.Nil(t, b)
		})
	}
}
