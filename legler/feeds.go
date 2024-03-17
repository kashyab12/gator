package legler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kashyab12/gator/internal/database"
	"net/http"
	"time"
)

type PostFeedsBody struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (config *ApiConfig) PostFeedsLegler(w http.ResponseWriter, r *http.Request, user database.User) {
	var (
		decoder  = json.NewDecoder(r.Body)
		feedBody = PostFeedsBody{}
	)
	if decodeErr := decoder.Decode(&feedBody); decodeErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, decodeErr.Error())
	} else if newFeed, feedInsertErr := config.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedBody.Name,
		Url:       feedBody.URL,
		UserID:    user.ID,
	}); feedInsertErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, feedInsertErr.Error())
	} else if responseErr := RespondWithJson(w, http.StatusCreated, struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		URL       string    `json:"url"`
		UserID    uuid.UUID `json:"user_id"`
	}{newFeed.ID, newFeed.CreatedAt, newFeed.UpdatedAt, newFeed.Name, newFeed.Url, newFeed.UserID}); responseErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, responseErr.Error())
	}
}

func (config *ApiConfig) GetFeedsLegler(w http.ResponseWriter, r *http.Request) {
	if allFeeds, fetchFeedsErr := config.DB.GetFeeds(r.Context()); fetchFeedsErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, fetchFeedsErr.Error())
	} else if responseErr := RespondWithJson(w, http.StatusOK, allFeeds); responseErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, responseErr.Error())
	}
}
