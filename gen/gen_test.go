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

			assert.NotNil(t, b)
			assert.NotEmpty(t, b.Description)
			assert.NotNil(t, b.Info)
			assert.NotEmpty(t, b.Friends)
			assert.NotEmpty(t, b.Category)
			assert.Len(t, b.Category, 3)
			assert.NotEmpty(t, b.Category)

			for _, s := range b.Category {
				//t.Log(s.Name)
				//for _, article := range s.Articles {
				//	t.Log(article.String())
				//}
				assert.NotNil(t, s)
			}
		})
	}
}
