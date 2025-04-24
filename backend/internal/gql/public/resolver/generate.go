//go:build ignore

package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
)

var jsonTagRegexp = regexp.MustCompile(`json:"([^"]+)"`)

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

func snakeCaseMutateHook(b *modelgen.ModelBuild) *modelgen.ModelBuild {
	for _, model := range b.Models {
		for _, field := range model.Fields {
			matches := jsonTagRegexp.FindStringSubmatch(field.Tag)
			if len(matches) == 2 {
				snake := toSnakeCase(matches[1])
				field.Tag = jsonTagRegexp.ReplaceAllString(field.Tag, fmt.Sprintf(`json:"%s"`, snake))
			}
		}
	}
	return b
}

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}

	p := modelgen.Plugin{
		MutateHook: snakeCaseMutateHook,
	}

	err = api.Generate(cfg, api.ReplacePlugin(&p))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}
