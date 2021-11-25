package pages

import (
	"fmt"
	fdb "forum/internal/forumDB"
	"html/template"
	"net/http"
	"strings"
)

// ThreadsPage represents a single page's worth of threads.
// To do this, it has all kinds of configuration variables, which are used to fill out the Threads field
type ThreadsPage struct {
	BoardID int

	Page         int
	PageSize     int
	PagesTotal   int
	ThreadsTotal int

	Order     string
	Ascending bool
	SQLKey    string

	Threads []fdb.Thread
}

// GetThreadsPage returns a ThreadsPage struct, which represents a single page's worth of threads.
func (env Board) GetThreadsPage(boardID int, r *http.Request) (ThreadsPage, error) {
	page, err := env.initThreadPage(boardID, r)
	if err != nil {
		return ThreadsPage{}, err
	}

	page.Threads, err = env.Threads.GetPageThreads(boardID, page.Page, page.PageSize, page.SQLKey)
	if err != nil {
		return ThreadsPage{}, err
	}

	return page, nil
}

// Initializes the ThreadsPage variables based on query values from the request.
func (env Board) initThreadPage(boardID int, r *http.Request) (ThreadsPage, error) {
	page := ThreadsPage{}
	page.BoardID = boardID

	// Parse current page and page size. Need both to build the UI correctly in the template
	page.Page, _ = GetQueryInt("p", r)
	page.PageSize, _ = GetQueryInt("plen", r)
	switch {
	case page.PageSize > 100:
		page.PageSize = 100
	case page.PageSize <= 0:
		page.PageSize = 16
	}

	page.ThreadsTotal = env.Threads.ThreadCount(boardID)
	page.PagesTotal = page.ThreadsTotal / page.PageSize
	if page.ThreadsTotal%page.PageSize != 0 {
		page.PagesTotal++
	}

	if page.Page > page.PagesTotal {
		page.Page = page.PagesTotal
	}

	if page.Page < 1 {
		page.Page = 1
	}

	// Order query for "date" can be "date", "date-asc", or "date-desc".
	// In the first case, it will default the second half.
	slc := strings.SplitN(r.URL.Query().Get("order"), "-", 2)

	page.Order = slc[0]
	switch page.Order {
	default:
		fallthrough
	case "date":
		page.SQLKey = "pl.postID"
		page.Ascending = false
	case "name":
		page.SQLKey = "t.Title"
		page.Ascending = true
	case "replies":
		page.SQLKey = "CountPosts"
		page.Ascending = false
	}

	if len(slc) == 2 {
		switch slc[1] {
		case "asc":
			page.Ascending = true
		case "desc":
			page.Ascending = false
		}
	}

	if page.Ascending {
		page.SQLKey += "-ASC"
	} else {
		page.SQLKey += "-DESC"
	}

	return page, nil
}

// Funcs for use in templates
func (data boardData) PagesBefore() (slc []int) {
	page := data.ThreadsPage.Page
	for i := page - 1; i >= 1; i-- {
		slc = append(slc, i)
	}
	return slc
}

func (data boardData) PagesAfter() (slc []int) {
	page := data.ThreadsPage.Page
	total := data.ThreadsPage.PagesTotal
	for i := page + 1; i <= total; i++ {
		slc = append(slc, i)
	}
	return slc
}

func (data boardData) ThreadColumnAttributes(filter string) (class template.HTMLAttr, err error) {
	var hrefUrl template.URL
	var orderPlus string
	if data.ThreadsPage.Order == filter {
		if data.ThreadsPage.Ascending {
			class = `class="active-order ascending"`
			orderPlus = "-desc"
		} else {
			class = `class="active-order descending"`
			orderPlus = "-asc"
		}
	}

	hrefUrl, err = data.Query("order", filter+orderPlus)
	if err != nil {
		return "", err
	}

	str := fmt.Sprintf(`%v href="%v#threads"`, class, hrefUrl)
	return template.HTMLAttr(str), nil
}
