package tempserver

import (
	"fmt"
	"testing"
	"text/template"
)

func TestStart(t *testing.T) {

	tmpl, err := template.New("test").Parse("{{ .A }}, {{ .B }}\n")

	if err != nil {
		panic(err)
	}

	c := &Config{
		Path: "tail",
		Arguments: func(config string) []string {
			fmt.Println("Config: " + config)
			return []string{"-f", config}
		},
		Config: struct {
			A string
			B string
		}{
			A: "a",
			B: "b",
		},
		ConfigTemplate: tmpl,
		WaitFor:        "a, b",
	}

	s, err := Start(c)

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	s.Term()

}
