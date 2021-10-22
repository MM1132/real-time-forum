package pages

import "fmt"

type noTemplateError string

func (err noTemplateError) Error() string {
	return fmt.Sprintf("Could not find the template \"%v\"", err)
}
