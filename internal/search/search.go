package search

import (
	"fmt"
	"forum/internal/forumEnv"
	"html/template"
	"net/http"
	"net/url"
)

type Searcher interface {
	GetPrimitive() SearchPrimitive

	ProcessString(string)
	ProcessRequest(*http.Request)
	ProcessRequestBasic(*http.Request)
	DoSearch(forumEnv.Env) error
	HasResult() bool
	GetResult() template.HTML
}

type SearchPrimitive struct {
	MyUserID   int
	template   *template.Template
	Name       string
	CurrentURL url.URL
	RawSearch  string

	Page      int
	PageLen   int
	PageCount int

	Order       string
	Ascending   bool
	OrderSqlKey string
}

// TEMPLATE FUNCTIONS

func (s SearchPrimitive) PagesBefore() (slc []int) {
	page := s.Page
	for i := page - 1; i >= 1; i-- {
		slc = append(slc, i)
	}
	return slc
}

func (s SearchPrimitive) PagesAfter() (slc []int) {
	page := s.Page
	total := s.PageCount
	for i := page + 1; i <= total; i++ {
		slc = append(slc, i)
	}
	return slc
}

func (s SearchPrimitive) ThreadColumnAttributes(filter string) (class template.HTMLAttr, err error) {
	var hrefUrl template.URL
	var orderPlus string
	if s.Order == filter {
		if s.Ascending {
			class = `class="active-order ascending"`
			orderPlus = "-desc"
		} else {
			class = `class="active-order descending"`
			orderPlus = "-asc"
		}
	}

	hrefUrl, err = s.Query("order", filter+orderPlus)
	if err != nil {
		return "", err
	}

	str := fmt.Sprintf(`%v href="%v#threads"`, class, hrefUrl)
	return template.HTMLAttr(str), nil
}

func (s SearchPrimitive) Query(kvp ...string) (template.URL, error) {
	if len(kvp)%2 == 1 {
		return "", fmt.Errorf(`need an even number of args`)
	}

	currentURL := s.CurrentURL

	query := currentURL.Query()
	for i := 0; i < len(kvp); i += 2 {
		if kvp[i+1] != "" {
			query.Set(kvp[i], kvp[i+1])
		}
	}

	currentURL.RawQuery = query.Encode()
	return template.URL(currentURL.RequestURI()), nil
}
