package main

// User type
type User struct {
	Name     string
	Password string
}

// Message type
type Message struct {
	Text   string
	Time   string
	Sender string
}

// Post type
type Post struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Src         string `json:"src"`
	Description string `json:"description"`
	Likes       int    `json:"likes"`
}

// Comment type
type Comment struct {
	Text   string `json:"text"`
	PostID string `json:"postId"`
	Time   int64  `json:"time"`
	Sender string `json:"sender"`
}
