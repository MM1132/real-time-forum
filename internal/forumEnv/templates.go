package forumEnv

import (
	"html/template"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
)

func CreateTemplates(rootPath string) map[string]*template.Template {
	// New file system for globbing
	fileSystem := os.DirFS(rootPath)

	// One path for every template we generate
	templPaths, _ := fs.Glob(fileSystem, "*.html")

	// We will be storing templates here
	templates := make(map[string]*template.Template)

	for _, templPath := range templPaths {
		// start with server/example.html
		base := path.Base(templPath)                     // example.html
		name := strings.TrimSuffix(base, path.Ext(base)) // example

		templates[name] = createTemplate(fileSystem, templPath, name)
	}

	return templates
}

func createTemplate(fileSystem fs.FS, templPath, name string) *template.Template {
	// Add all the HTML files in the specific folder into the slice as patterns
	specificFiles, _ := fs.Glob(fileSystem, path.Join("specific", name, "*.html"))
	patterns := append(
		[]string{templPath, "universal/*.html"},
		specificFiles...,
	)

	// Put all the files that match the patterns into the template
	newTemplate := template.Must(template.New(name).Funcs(getFuncMap()).ParseFS(fileSystem, patterns...))

	return newTemplate
}

// Set custom functions for all templates
func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"itoa": func(num int) string {
			return strconv.Itoa(num)
		},
		"divide": func(a, b int) int {
			return a / b
		},
	}
}
