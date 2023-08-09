package goforum

import (
	"log"
	"os"
)

type Post struct {
	ID       int    // Unique ID of the post
	ThreadID int64  // Used to find other posts in the thread
	OrderID  int    // The order of this post in the thread
	Username string // User who made the post
	Title    string // Title of the thread the post is in
	Body     string // The actual message of the post
	Date     string // Date when the post was made
}

// Gets posts from thread ID
func GetPosts(threadId string) (posts []Post) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	rows, err := db.Query("SELECT * FROM posts WHERE threadId=?", threadId)

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var post Post

	for rows.Next() {

		err = rows.Scan(&post.ID, &post.ThreadID, &post.OrderID, &post.Username,
			&post.Title, &post.Body, &post.Date)

		if err != nil {
			log.Println(err)
		}
		posts = append(posts, post)
	}
	updateThreadViewCount(post.ThreadID)
	return posts
}

// Returns true if post was created with no errors
func (post *Post) createPost() bool {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	stmt, err := db.Prepare(`INSERT posts SET threadId=?, orderId=?, username=?, title=?
		, body=?, date=?`)
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = stmt.Exec(post.ThreadID, post.OrderID, post.Username, post.Title,
		post.Body, post.Date)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Updates forum and thread last post
func updateLastPost(username string, threadId int64, forumId int) {
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	stmt, err := db.Prepare("UPDATE threads SET lastpost=? WHERE id=?")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = stmt.Exec(username, threadId)
	if err != nil {
		log.Println(err)
	}

	stmt, err = db.Prepare("UPDATE forums SET recentuser=? WHERE id=?")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = stmt.Exec(username, forumId)
	if err != nil {
		log.Println(err)
	}
}

func updateThreadReplies(threadId int64) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	stmt, err := db.Prepare("UPDATE threads SET replies=replies+1 WHERE id=?")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = stmt.Exec(threadId)
	if err != nil {
		log.Println(err)
	}
}
