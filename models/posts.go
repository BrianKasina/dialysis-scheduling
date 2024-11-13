package models

type Post struct {
    ID       int    `json:"id" bson:"post_id"`
    Title    string `json:"title" bson:"title"`
    Content  string `json:"content" bson:"content"`
    AdminID  int    `json:"admin_id,omitempty" bson:"admin_id"`
    AdminName string `json:"admin_name,omitempty" bson:"admin_name"`
    PostDate string `json:"post_date" bson:"post_date"`
    PostTime string `json:"post_time" bson:"post_time"`
}