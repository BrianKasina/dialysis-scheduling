package controllers

import (
    "encoding/json"
    "net/http"
    "math"
    "github.com/BrianKasina/dialysis-scheduling/gateways"
    "github.com/BrianKasina/dialysis-scheduling/utils"
    "github.com/BrianKasina/dialysis-scheduling/models"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/gorilla/mux"
)

type PostController struct {
    PostGateway *gateways.PostGateway
}

func NewPostController(db *mongo.Database) *PostController {
    return &PostController{
        PostGateway: gateways.NewPostGateway(db),
    }
}

// Handle GET requests for posts with pagination
func (pc *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
    limit, _ := r.Context().Value("limit").(int) // Retrieve limit from context
    page, _ := r.Context().Value("page").(int)   // Retrieve page from context
    identifier := r.URL.Query().Get("identifier")

    offset := (page - 1) * limit // Calculate the actual offset

    query := r.URL.Query().Get("query")
    var posts []models.Post
    var err error

    if identifier == "search" {
        posts, err = pc.PostGateway.SearchPosts(query, limit, offset)
    } else {
        posts, err = pc.PostGateway.GetPosts(limit, offset)
    }

    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch posts")
        return
    }

    totalEntries, err := pc.PostGateway.GetTotalPostCount(query)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch total posts count")
        return
    }

    totalPages := int(math.Ceil(float64(totalEntries) / float64(limit)))

    response := map[string]interface{}{
        "data":         posts,
        "total_pages":  totalPages,
        "page": page,
        "total_entries": totalEntries,
    }

    json.NewEncoder(w).Encode(response)
}

// Handle POST requests for posts
func (pc *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
    var post models.Post
    err := json.NewDecoder(r.Body).Decode(&post)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    err = pc.PostGateway.CreatePost(&post)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to create post")
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}

// Handle PUT requests for posts
func (pc *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
    var post models.Post
    err := json.NewDecoder(r.Body).Decode(&post)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    err = pc.PostGateway.UpdatePost(&post)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to update post")
        return
    }
    json.NewEncoder(w).Encode(post)
}

// Handle DELETE requests for posts
func (pc *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID := vars["id"]
    if postID == "" {
        utils.ErrorHandler(w, http.StatusBadRequest, nil, "Missing post ID")
        return
    }

    err := pc.PostGateway.DeletePost(postID)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to delete post")
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})
}