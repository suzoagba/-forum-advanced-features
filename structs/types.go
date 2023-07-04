package structs

type Post struct {
	ID            int
	Username      string
	Title         string
	Description   string
	CreationDate  string
	Tags          []string
	ImageFileName string
	Likes         int
	Dislikes      int
	Edited        bool
	TimeEdited    string
}

type Comment struct {
	ID           string
	Content      string
	PostID       int
	UserID       int
	Username     string
	CreationDate string
	Likes        int
	Dislikes     int
	Edited       bool
	TimeEdited   string
}

type User struct {
	ID       string
	Username string // Display the name of the user who is logged in
	LoggedIn bool
}

type ErrorMessage struct {
	Error   bool
	Message string
	Field1  string
	Field2  string
	Field3  []string
	Image   string
}

type Tag struct {
	ID   int
	Name string
}

type CommentListing struct {
	Post     Post
	Comments []Comment
}

type Activity struct {
	Posts           []Post
	Comments        []CommentListing
	PostLikes       []Post
	PostDislikes    []Post
	CommentLikes    []CommentListing
	CommentDislikes []CommentListing
}

type Notification struct {
	ID          int
	User        string
	Who         string
	ActionDone  string
	IsPost      bool
	IsComment   bool
	PostID      int
	CommentID   int
	IsRead      bool
	CreatedDate string
}

type ForPage struct {
	Error               ErrorMessage
	User                User
	Posts               []Post
	Tags                []Tag
	Comments            []Comment
	OAuth               OAuth
	Activity            Activity
	Notifications       []Notification
	UnreadNotifications int
}

type OAuth struct {
	GoogleID string
	GitHubID string
}

type OAuthUser struct {
	email string
	other string
}
