//nolint:forbidigo // enable printing function
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pagu-project/pagu/internal/engine/command"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//go:embed command.tmpl
var commandTemplate string

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: generator <path_to_yml_file> ...")
	}
	flag.Parse()
	args := flag.Args()

	for _, ymlPath := range args {
		data, err := os.ReadFile(ymlPath)
		if err != nil {
			fmt.Printf("Error reading YAML file %s: %v\n", ymlPath, err)
			os.Exit(1)
		}

		cmd := new(command.Command)
		if err := yaml.Unmarshal(data, &cmd); err != nil {
			fmt.Printf("Error unmarshalling YAML file %s: %v\n", ymlPath, err)
			os.Exit(1)
		}

		code, err := generateCode(cmd)
		if err != nil {
			fmt.Printf("Unable to generate the code: %v\n", err)
			os.Exit(1)
		}

		baseName := strings.TrimSuffix(filepath.Base(ymlPath), ".yml")
		outputFile := filepath.Join(filepath.Dir(ymlPath), baseName+".gen.go")
		if err := os.WriteFile(outputFile, []byte(code), 0o600); err != nil {
			fmt.Printf("Error writing Go file %s: %v\n", outputFile, err)
			os.Exit(1)
		}

		fmt.Printf("Generated code for %s command\n", baseName)
	}
}

func generateCode(cmd *command.Command) (string, error) {
	funcMap := template.FuncMap{
		"title": func(str string) string {
			return cases.Title(language.English).String(str)
		},

		"string": func(s fmt.Stringer) string {
			return s.String()
		},
	}

	tml, err := template.New("code").Funcs(funcMap).Parse(commandTemplate)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	if err := tml.Execute(&sb, cmd); err != nil {
		return "", err
	}

	return sb.String(), nil
}
