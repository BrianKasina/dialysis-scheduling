package gateways

import (
    "database/sql"
    "github.com/BrianKasina/dialysis-scheduling/models"
)

type PostGateway struct {
    db *sql.DB
}

func NewPostGateway(db *sql.DB) *PostGateway {
    return &PostGateway{db: db}
}

func (pg *PostGateway) GetPosts(limit, offset int) ([]models.Post, error) {
    rows, err := pg.db.Query(`
        SELECT p.post_id, p.title, p.content, p.post_date, p.post_time, sa.name AS admin_name
        FROM posts p
        JOIN system_admin sa ON p.admin_id = sa.admin_id
        LIMIT ? OFFSET ?
    `, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []models.Post
    for rows.Next() {
        var post models.Post
        if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.PostDate, &post.PostTime, &post.AdminName); err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }
    return posts, nil
}

func (pg *PostGateway) SearchPosts(query string, limit, offset int) ([]models.Post, error) {
    searchQuery := "%" + query + "%"
    rows, err := pg.db.Query(`
        SELECT p.post_id, p.title, p.content, p.post_date, p.post_time, sa.name AS admin_name
        FROM posts p
        JOIN system_admin sa ON p.admin_id = sa.admin_id
        WHERE CONCAT(p.title, ' ', p.content, ' ', sa.name) LIKE ?
        LIMIT ? OFFSET ?
    `, searchQuery, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []models.Post
    for rows.Next() {
        var post models.Post
        if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.PostDate, &post.PostTime, &post.AdminName); err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }
    return posts, nil
}

func (pg *PostGateway) GetTotalPostCount(query string) (int, error) {
    var row *sql.Row
    if query != "" {
        searchQuery := "%" + query + "%"
        row = pg.db.QueryRow(`
            SELECT COUNT(*)
            FROM posts p
            JOIN system_admin sa ON p.admin_id = sa.admin_id
            WHERE CONCAT(p.title, ' ', p.content, ' ', sa.name) LIKE ?
        `, searchQuery)
    } else {
        row = pg.db.QueryRow("SELECT COUNT(*) FROM posts")
    }

    var count int
    err := row.Scan(&count)
    if err != nil {
        return 0, err
    }

    return count, nil
}

func (pg *PostGateway) CreatePost(post *models.Post) error {
    _, err := pg.db.Exec(
        `INSERT INTO posts (title, content, post_date, post_time, admin_id) 
        VALUES (?, ?, ?, ?, 
        (SELECT admin_id FROM system_admin WHERE name = ?)
        )`,
        post.Title, post.Content, post.PostDate, post.PostTime, post.AdminName)
    return err
}

func (pg *PostGateway) UpdatePost(post *models.Post) error {
    _, err := pg.db.Exec(
        `UPDATE posts SET title = ?, content = ?, post_date = ?, post_time = ?, 
         admin_id = (SELECT admin_id FROM system_admin WHERE name = ?) 
         WHERE post_id = ?`,
        post.Title, post.Content, post.PostDate, post.PostTime, post.AdminName, post.ID)
    return err
}

func (pg *PostGateway) DeletePost(postID string) error {
    _, err := pg.db.Exec("DELETE FROM posts WHERE post_id = ?", postID)
    return err
}