package models

type Post struct {
    ID       int    `json:"id" db:"post_id"`
    Title    string `json:"title" db:"title"`
    Content  string `json:"content" db:"content"`
    AdminID  int    `json:"admin_id" db:"admin_id"`
    PostDate string `json:"post_date" db:"post_date"`
    PostTime string `json:"post_time" db:"post_time"`
}