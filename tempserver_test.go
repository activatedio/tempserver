package tempserver

import (
	"testing"
	"text/template"
)

func TestStart(t *testing.T) {

	cases := map[string]struct {
		path      string
		template  string
		arguments func(configPath string) []string
		config    interface{}
		waitFor   string
	}{
		"tail with config file": {
			path:     "tail",
			template: "{{ .A }}, {{ .B }}\n",
			arguments: func(configPath string) []string {
				return []string{"-f", configPath}
			},
			config: struct {
				A string
				B string
			}{
				A: "a",
				B: "b",
			},
			waitFor: "a, b",
		},
		"sleep with no config file": {
			path: "sleep",
			arguments: func(configPath string) []string {
				return []string{"60"}
			},
		},
	}

	for k, v := range cases {

		t.Run(k, func(t *testing.T) {

			var tmpl *template.Template

			if v.template != "" {

				var err error
				tmpl, err = template.New("test").Parse("{{ .A }}, {{ .B }}\n")

				if err != nil {
					panic(err)
				}
			}

			c := &Config{
				Path:           v.path,
				Arguments:      v.arguments,
				Config:         v.config,
				ConfigTemplate: tmpl,
				WaitFor:        v.waitFor,
			}

			s, err := Start(c)

			if err != nil {
				t.Fatalf("Unexpected error %s", err)
			}

			s.Term()
		})
	}

}
