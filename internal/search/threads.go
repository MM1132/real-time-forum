package search

import (
	"bytes"
	"database/sql"
	fdb "forum/internal/forumDB"
	"forum/internal/forumEnv"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ThreadSearch is used to display a list of threads based on search parameters
type ThreadSearch struct {
	SearchPrimitive

	ThreadTitle sql.NullString
	TagName     sql.NullString
	Author      sql.NullString
	BoardID     sql.NullInt64
	BoardName   sql.NullString

	LatestAfter  sql.NullTime
	LatestBefore sql.NullTime
	OldestAfter  sql.NullTime
	OldestBefore sql.NullTime

	Results []fdb.Thread
}

// NewThreadSearch returns a ThreadSearch struct, which can be assigned filters and then executed with DoSearch
func NewThreadSearch(env forumEnv.Env, currentURL url.URL, myUserID int) *ThreadSearch {
	s := &ThreadSearch{}

	s.MyUserID = myUserID
	s.template = env.Templates["search"]
	s.CurrentURL = currentURL
	s.Page = 1
	s.PageLen = 16

	return s
}

func (s *ThreadSearch) GetPrimitive() SearchPrimitive {
	return s.SearchPrimitive
}

// Execute the search with currently set filters
func (s *ThreadSearch) DoSearch(env forumEnv.Env) error {
	params := make(map[string]interface{})

	params["page"] = s.Page
	params["pageLen"] = s.PageLen
	params["threadTitle"] = s.ThreadTitle
	params["tagName"] = s.TagName
	params["author"] = s.Author
	params["boardID"] = s.BoardID
	params["boardName"] = s.BoardName
	params["latestAfter"] = s.LatestAfter
	params["latestBefore"] = s.LatestBefore
	params["oldestAfter"] = s.OldestAfter
	params["oldestBefore"] = s.OldestBefore

	if s.Page <= 0 {
		params["page"] = 1
	}

	var err error
	s.Results, err = env.Threads.Search(s.OrderSqlKey, params)
	if err != nil {
		return err
	}

	count, err := env.Threads.SearchCount(params)
	if err != nil {
		return err
	}
	s.PageCount = count / s.PageLen
	if count%s.PageLen > 0 {
		s.PageCount++
	}

	return nil
}

func (s *ThreadSearch) HasResult() bool {
	return len(s.Results) > 0
}

// Gets the result as HTML, for direct use in templates
func (s *ThreadSearch) GetResult() template.HTML {
	buf := new(bytes.Buffer)

	tmpl := s.template
	if err := tmpl.ExecuteTemplate(buf, "threads", s); err != nil {
		log.Print(err)
	}

	return template.HTML(buf.String())
}

// Processes a search string by reading filters from it
func (s *ThreadSearch) ProcessString(str string) {
	s.RawSearch = str

	r1 := regexp.MustCompile(` (tag|board|author)\((.+)\)`)
	params := r1.FindAllStringSubmatch(str, -1)
	str = r1.ReplaceAllString(str, "")

	for _, param := range params {
		switch param[1] {
		case "board":
			s.BoardName.String, s.BoardName.Valid = param[2], true
		case "tag":
			s.TagName.String, s.TagName.Valid = param[2], true
		case "author":
			s.Author.String, s.Author.Valid = param[2], true
		}
	}

	str = strings.TrimSpace(str)
	if len(str) <= 0 {
		return
	}

	s.ThreadTitle.String, s.ThreadTitle.Valid = str, true
	return
}

// Processes a handler's request and gets filters from it
func (s *ThreadSearch) ProcessRequest(r *http.Request) {
	s.processRequest(r, false)
}

// Processes a handler's request and only gets page number and sorting order from it
func (s *ThreadSearch) ProcessRequestBasic(r *http.Request) {
	s.processRequest(r, true)
}

func (s *ThreadSearch) processRequest(r *http.Request, basic bool) {
	page, err := strconv.Atoi(r.FormValue("page"))
	if err == nil && page > 0 {
		s.Page = page
	}

	pageLen, err := strconv.Atoi(r.FormValue("pageLen"))
	if err == nil && pageLen > 0 && pageLen <= 100 {
		s.PageLen = pageLen
	}

	order := r.URL.Query().Get("order")
	if order == "" && r.Method == "POST" {
		err := r.ParseForm()
		if err == nil {
			a := r.Form.Get("order")
			b := r.Form.Get("ascDesc")
			order = a + "-" + b
		}
	}
	s.ProcessOrder(order)

	if basic {
		return
	}

	found := r.FormValue("title")
	if found != "" {
		s.ThreadTitle.String, s.ThreadTitle.Valid = found, found != ""
	}

	found = r.FormValue("tag")
	if found != "" {
		s.TagName.String, s.TagName.Valid = found, found != ""
	}

	found = r.FormValue("author")
	if found != "" {
		s.Author.String, s.Author.Valid = found, found != ""
	}

	boardID, err := strconv.Atoi(r.FormValue("boardID"))
	if err == nil && boardID >= 0 {
		s.BoardID.Int64, s.BoardID.Valid = int64(boardID), true
	}

	found = r.FormValue("boardName")
	if found != "" {
		s.BoardName.String, s.BoardName.Valid = found, found != ""
	}

	timeFmt := "2006-01-02"

	t, err := time.Parse(timeFmt, r.FormValue("after"))
	if err == nil {
		s.LatestAfter.Time, s.LatestAfter.Valid = t, true
	}

	t, err = time.Parse(timeFmt, r.FormValue("before"))
	if err == nil {
		s.LatestBefore.Time, s.LatestBefore.Valid = t, true
	}

	t, err = time.Parse(timeFmt, r.FormValue("updatedAfter"))
	if err == nil {
		s.OldestAfter.Time, s.OldestAfter.Valid = t, true
	}

	t, err = time.Parse(timeFmt, r.FormValue("updatedBefore"))
	if err == nil {
		s.OldestBefore.Time, s.OldestBefore.Valid = t, true
	}

	return
}

// Processes an order string (eg. "date-asc")
func (s *ThreadSearch) ProcessOrder(orderStr string) {
	slc := strings.SplitN(orderStr, "-", 2)

	s.Order = slc[0]
	switch s.Order {
	default:
		s.Order = "date"
		fallthrough
	case "date":
		s.OrderSqlKey = "pl.postID"
		s.Ascending = false
	case "name":
		s.OrderSqlKey = "t.Title"
		s.Ascending = true
	case "replies":
		s.OrderSqlKey = "CountPosts"
		s.Ascending = false
	}

	if len(slc) == 2 {
		switch slc[1] {
		case "asc":
			s.Ascending = true
		case "desc":
			s.Ascending = false
		}
	}

	if s.Ascending {
		s.OrderSqlKey += "-ASC"
	} else {
		s.OrderSqlKey += "-DESC"
	}
}
