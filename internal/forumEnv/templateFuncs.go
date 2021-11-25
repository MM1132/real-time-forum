package forumEnv

import (
	"fmt"
	"html/template"
)

// Query returns the current URI, except with new key/value pairs added. Setting an already existing key will replace its value.
// This function is meant to be used in templates for hrefs.
func (data GenericData) Query(kvp ...string) (template.URL, error) {
	if len(kvp)%2 == 1 {
		return "", fmt.Errorf(`need an even number of args`)
	}

	currentURL := data.CurrentURL

	query := currentURL.Query()
	for i := 0; i < len(kvp); i += 2 {
		query.Set(kvp[i], kvp[i+1])
	}

	currentURL.RawQuery = query.Encode()
	return template.URL(currentURL.RequestURI()), nil
}
