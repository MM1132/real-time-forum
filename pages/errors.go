package pages

import "fmt"

type noTemplateError struct {
	name string
}

func (err noTemplateError) Error() string {
	return fmt.Sprintf("Could not find the template \"%s\"", err.name)
}
