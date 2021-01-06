package gen

import (
	"testing"
)

func TestIndexTemplate_execute(t *testing.T) {
	type fields struct {
		Template *Template
	}
	type args struct {
		variables  interface{}
		outputPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "\\",
			fields: fields{
				Template: &Template{
					Path:      "..\\template\\template_index.html",
					Variables: nil,
				},
			},
			args: args{
				variables:  From("..\\sample_repo"),
				outputPath: "..\\out\\index.html",
			},
			wantErr: false,
		},
		{
			name: "OutputNotExist",
			fields: fields{
				Template: &Template{
					Path:      "..\\template\\template_index.html",
					Variables: nil,
				},
			},
			args: args{
				variables:  From("..\\sample_repo"),
				outputPath: "..\\out\\index.html",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			that := IndexTemplate{
				Template: tt.fields.Template,
			}
			if err := that.execute(tt.args.variables, tt.args.outputPath); (err != nil) != tt.wantErr {
				t.Errorf("execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConvertConfig_validate(t *testing.T) {
	type fields struct {
		OutputDir   string
		TemplateDir string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "\\",
			fields: fields{
				OutputDir:   "..\\out",
				TemplateDir: "..\\template",
			},
			wantErr: false,
		},
		{
			name: "DirNotExist",
			fields: fields{
				OutputDir:   "not_exist_dir",
				TemplateDir: "not_exist_dir",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			that := &ConvertConfig{
				OutputDir:   tt.fields.OutputDir,
				TemplateDir: tt.fields.TemplateDir,
			}
			if err := that.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
