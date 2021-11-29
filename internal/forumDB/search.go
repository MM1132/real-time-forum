package forumDB

import (
	"database/sql"
	"fmt"
)

// THREAD SEARCH //

// Make a search based on parameters set in the ThreadSearch struct
func (m ThreadModel) Search(orderKey string, params map[string]interface{}) ([]Thread, error) {
	key := "SearchThreads-" + orderKey
	stmt, ok := m.statements[key]
	if !ok {
		panic(fmt.Errorf("could not find the %v statement in threads.sql", key))
	}

	rows, err := stmt.Query(
		params["page"],
		params["pageLen"],
		params["threadTitle"],
		params["tagName"],
		params["author"],
		params["boardID"],
		params["boardName"],
		params["latestAfter"],
		params["latestBefore"],
		params["oldestAfter"],
		params["oldestBefore"],
	)
	if err != nil {
		return nil, err
	}

	var threads []Thread
	for rows.Next() {
		thread := Thread{}
		thread.Extras = &ThreadExtras{}
		err = rows.Scan(
			&thread.ThreadID,
			&thread.Title,
			&thread.BoardID,

			&thread.Extras.CountPosts,
			&thread.Extras.CountUsers,

			&thread.Extras.LatestID,
			&thread.Extras.LatestAuthorID,
			&thread.Extras.LatestAuthor,
			&thread.Extras.LatestDate,

			&thread.Extras.OldestID,
			&thread.Extras.OldestAuthorID,
			&thread.Extras.OldestAuthor,
			&thread.Extras.OldestDate,
		)
		if err != nil {
			return nil, err
		}

		thread.Tags = m.Tags.GetByThread(thread.ThreadID)

		threads = append(threads, thread)
	}

	return threads, nil
}

// Count how many results a thread search with certain parameters would return
func (m ThreadModel) SearchCount(params map[string]interface{}) (int, error) {
	key := "SearchThreadsCount"
	stmt, ok := m.statements[key]
	if !ok {
		panic(fmt.Errorf("could not find the %v statement in threads.sql", key))
	}

	row := stmt.QueryRow(
		params["page"],
		params["pageLen"],
		params["title"],
		params["tagName"],
		params["author"],
		params["boardID"],
		params["boardName"],
		params["latestAfter"],
		params["latestBefore"],
		params["oldestAfter"],
		params["oldestBefore"],
	)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// POST SEARCH //

type SearchPost struct {
	Post
	ThreadName string
	BoardID    int
	BoardName  string
}

// Make a search based on parameters set in the PostSearch struct
func (m PostModel) Search(orderKey string, params map[string]interface{}) ([]SearchPost, error) {
	key := "SearchPosts-" + orderKey
	stmt, ok := m.statements[key]
	if !ok {
		panic(fmt.Errorf("could not find the %v statement in posts.sql", key))
	}

	rows, err := stmt.Query(
		params["myUserID"],
		params["page"],
		params["pageLen"],
		params["threadID"],
		params["threadTitle"],
		params["tagName"],
		params["author"],
		params["authorID"],
		params["boardID"],
		params["boardName"],
		params["likedBy"],
		params["likedByID"],
		params["dislikedBy"],
		params["dislikedByID"],
		params["content"],
		params["After"],
		params["Before"],
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	var posts []SearchPost
	for rows.Next() {
		post := SearchPost{}
		err = rows.Scan(
			&post.PostID,
			&post.ThreadID,
			&post.UserID,
			&post.Content,
			&post.Date,
			&post.User.UserID,
			&post.User.Name,
			&post.User.Email,
			&post.User.Password,
			&post.User.Image,
			&post.User.Description,
			&post.User.Creation,
			&post.Likes,
			&post.MyLike,

			&post.ThreadName,
			&post.BoardID,
			&post.BoardName,
		)
		if err != nil {
			return nil, err
		}

		post.User.Extras = &UserExtras{}

		posts = append(posts, post)
	}

	return posts, nil
}

// Count how many results a post search with certain parameters would return
func (m PostModel) SearchCount(params map[string]interface{}) (int, error) {
	key := "SearchPostsCount"
	stmt, ok := m.statements[key]
	if !ok {
		panic(fmt.Errorf("could not find the %v statement in posts.sql", key))
	}

	row := stmt.QueryRow(
		params["myUserID"],
		params["page"],
		params["pageLen"],
		params["threadID"],
		params["threadTitle"],
		params["tagName"],
		params["author"],
		params["authorID"],
		params["boardID"],
		params["boardName"],
		params["likedBy"],
		params["likedByID"],
		params["dislikedBy"],
		params["dislikedByID"],
		params["content"],
		params["After"],
		params["Before"],
	)

	var count int
	err := row.Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return count, nil
}
