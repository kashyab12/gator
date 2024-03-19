package legler

import (
	"github.com/google/uuid"
	"github.com/kashyab12/gator/internal/database"
	"net/http"
	"strconv"
	"time"
)

type PostBody struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Url         string     `json:"url"`
	Description string     `json:"description"`
	PublishedAt *time.Time `json:"published_at"`
	FeedID      uuid.UUID  `json:"feed_id"`
}

func (config *ApiConfig) GetPostByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	var queryLimit int32 = 5
	if queryLimitStr := r.URL.Query().Get("queryLimit"); queryLimitStr != "" {
		if parsedQueryLimit, convErr := strconv.Atoi(queryLimitStr); convErr == nil {
			// todo: eeee, toast if crazy query limit provided.
			queryLimit = int32(parsedQueryLimit)
		}
	}
	if postsByUser, fetchPostsErr := config.DB.GetPostsByUserID(r.Context(), database.GetPostsByUserIDParams{
		UserID: user.ID,
		Limit:  queryLimit,
	}); fetchPostsErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, fetchPostsErr.Error())
	} else if responseErr := RespondWithJson(w, http.StatusOK, dbPostToPostJson(postsByUser)); responseErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, responseErr.Error())
	}
}

func dbPostToPostJson(posts []database.Post) (postBodies []PostBody) {
	for _, post := range posts {
		var publishedAtConv *time.Time = nil
		if post.PublishedAt.Valid {
			publishedAtConv = &post.PublishedAt.Time
		}
		postBody := PostBody{
			ID:          post.ID,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			Title:       post.Title,
			Url:         post.Url,
			Description: post.Description,
			PublishedAt: publishedAtConv,
			FeedID:      post.FeedID,
		}
		postBodies = append(postBodies, postBody)
	}
	return postBodies
}
