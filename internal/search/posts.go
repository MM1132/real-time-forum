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

type PostSearch struct {
	SearchPrimitive

	Breadcrumbs bool

	PostID sql.NullInt64

	ThreadID     sql.NullInt64
	Title        sql.NullString
	TagName      sql.NullString
	BoardID      sql.NullInt64
	BoardName    sql.NullString
	LikedBy      sql.NullString
	LikedByID    sql.NullInt64
	DislikedBy   sql.NullString
	DislikedByID sql.NullInt64

	Content  sql.NullString
	Author   sql.NullString
	AuthorID sql.NullInt64

	After  sql.NullTime
	Before sql.NullTime

	Results []fdb.SearchPost
}

// NewPostSearch returns a PostSearch struct, which can be assigned filters and then executed with DoSearch
func NewPostSearch(env forumEnv.Env, currentURL url.URL, userID int) *PostSearch {
	s := &PostSearch{}

	s.MyUserID = userID
	s.template = env.Templates["search"]
	s.CurrentURL = currentURL
	s.Page = 1
	s.PageLen = 12

	s.Order = "date"
	s.OrderSqlKey = "p.postID-DESC"
	s.Ascending = false

	return s
}

func (s *PostSearch) GetPrimitive() SearchPrimitive {
	return s.SearchPrimitive
}

// Execute the search with currently set filters
func (s *PostSearch) DoSearch(env forumEnv.Env) error {
	params := make(map[string]interface{})

	params["myUserID"] = s.MyUserID
	params["page"] = s.Page
	params["pageLen"] = s.PageLen
	params["threadID"] = s.ThreadID
	params["threadTitle"] = s.Title
	params["tagName"] = s.TagName
	params["author"] = s.Author
	params["authorID"] = s.AuthorID
	params["boardID"] = s.BoardID
	params["boardName"] = s.BoardName
	params["likedBy"] = s.LikedBy
	params["likedByID"] = s.LikedByID
	params["dislikedBy"] = s.DislikedBy
	params["dislikedByID"] = s.DislikedByID
	params["content"] = s.Content
	params["after"] = s.After
	params["before"] = s.Before

	if s.Page <= 0 {
		params["page"] = 1
	}

	var err error
	s.Results, err = env.Posts.Search(env.Users, s.OrderSqlKey, params)
	if err != nil {
		return err
	}

	count, err := env.Posts.SearchCount(params)
	if err != nil {
		return err
	}
	s.PageCount = count / s.PageLen
	if count%s.PageLen > 0 {
		s.PageCount++
	}

	return nil
}

func (s *PostSearch) HasResult() bool {
	return len(s.Results) > 0
}

// Gets the result as HTML, for direct use in templates
func (s *PostSearch) GetResult() template.HTML {
	buf := new(bytes.Buffer)

	tmpl := s.template
	if err := tmpl.ExecuteTemplate(buf, "posts", s); err != nil {
		log.Print(err)
	}

	return template.HTML(buf.String())
}

// Processes a search string by reading filters from it
func (s *PostSearch) ProcessString(str string) {
	s.RawSearch = str

	r1 := regexp.MustCompile(` ?(tag|board|author|likedBy|dislikedBy|order)\((.+)\)`)
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
		case "likedBy":
			s.LikedBy.String, s.LikedBy.Valid = param[2], true
		case "dislikedBy":
			s.DislikedBy.String, s.DislikedBy.Valid = param[2], true
		case "order":
			s.ProcessOrder(param[2])
		}
	}

	str = strings.TrimSpace(str)
	if len(str) <= 0 {
		return
	}

	s.Content.String, s.Content.Valid = str, true
	return
}

// Processes a handler's request and gets filters from it
func (s *PostSearch) ProcessRequest(r *http.Request) {
	s.processRequest(r, false)
}

// Processes a handler's request and only gets page number and sorting order from it
func (s *PostSearch) ProcessRequestBasic(r *http.Request) {
	s.processRequest(r, true)
}

func (s *PostSearch) processRequest(r *http.Request, basic bool) {
	foundInt, err := strconv.Atoi(r.FormValue("page"))
	if err == nil && foundInt > 0 {
		s.Page = foundInt
	}

	foundInt, err = strconv.Atoi(r.FormValue("pageLen"))
	if err == nil && foundInt > 0 && foundInt <= 100 {
		s.PageLen = foundInt
	}

	foundInt, err = strconv.Atoi(r.FormValue("boardID"))
	if err == nil && foundInt >= 0 {
		s.PostID.Int64, s.PostID.Valid = int64(foundInt), true
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

	foundInt, err = strconv.Atoi(r.FormValue("threadID"))
	if err == nil && foundInt > 0 {
		s.ThreadID.Int64, s.ThreadID.Valid = int64(foundInt), true
	}

	found := r.FormValue("title")
	if found != "" {
		s.Title.String, s.Title.Valid = found, found != ""
	}

	found = r.FormValue("tag")
	if found != "" {
		s.TagName.String, s.TagName.Valid = found, found != ""
	}

	found = r.FormValue("author")
	if found != "" {
		s.Author.String, s.Author.Valid = found, found != ""
	}

	foundInt, err = strconv.Atoi(r.FormValue("authorID"))
	if err == nil && foundInt >= 0 {
		s.AuthorID.Int64, s.AuthorID.Valid = int64(foundInt), true
	}

	foundInt, err = strconv.Atoi(r.FormValue("boardID"))
	if err == nil && foundInt >= 0 {
		s.BoardID.Int64, s.BoardID.Valid = int64(foundInt), true
	}

	found = r.FormValue("boardName")
	if found != "" {
		s.BoardName.String, s.BoardName.Valid = found, found != ""
	}

	found = r.FormValue("likedBy")
	if found != "" {
		s.LikedBy.String, s.LikedBy.Valid = found, found != ""
	}

	foundInt, err = strconv.Atoi(r.FormValue("likedByID"))
	if err == nil && foundInt >= 0 {
		s.LikedByID.Int64, s.LikedByID.Valid = int64(foundInt), true
	}

	found = r.FormValue("dislikedBy")
	if found != "" {
		s.DislikedBy.String, s.DislikedBy.Valid = found, found != ""
	}

	foundInt, err = strconv.Atoi(r.FormValue("dislikedByID"))
	if err == nil && foundInt >= 0 {
		s.DislikedByID.Int64, s.DislikedByID.Valid = int64(foundInt), true
	}

	found = r.FormValue("content")
	if found != "" {
		s.Content.String, s.Title.Valid = found, found != ""
	}

	timeFmt := "2006-01-02 15:04:05"

	t, err := time.Parse(timeFmt, r.FormValue("after"))
	if err == nil {
		s.After.Time, s.After.Valid = t, true
	}

	t, err = time.Parse(timeFmt, r.FormValue("before"))
	if err == nil {
		s.Before.Time, s.Before.Valid = t, true
	}

	return
}

// Processes an order string (eg. "date-asc")
func (s *PostSearch) ProcessOrder(orderStr string) {
	slc := strings.SplitN(orderStr, "-", 2)

	s.Order = slc[0]
	switch s.Order {
	default:
		s.Order = "date"
		fallthrough
	case "date":
		s.OrderSqlKey = "p.postID"
		s.Ascending = false
	case "likes":
		s.OrderSqlKey = "likes"
		s.Ascending = false
	case "likeDate":
		s.OrderSqlKey = "l2.date"
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
