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

type FeedBody struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	URL           string     `json:"url"`
	UserID        uuid.UUID  `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at,omitempty"`
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
	} else if responseErr := RespondWithJson(w, http.StatusCreated, &FeedBody{
		ID:            newFeed.ID,
		CreatedAt:     newFeed.CreatedAt,
		UpdatedAt:     newFeed.UpdatedAt,
		Name:          newFeed.Name,
		URL:           newFeed.Url,
		UserID:        newFeed.UserID,
		LastFetchedAt: nil,
	}); responseErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, responseErr.Error())
	}
}

func (config *ApiConfig) GetFeedsLegler(w http.ResponseWriter, r *http.Request) {
	if allFeeds, fetchFeedsErr := config.DB.GetFeeds(r.Context()); fetchFeedsErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, fetchFeedsErr.Error())
	} else if responseErr := RespondWithJson(w, http.StatusOK, dbFeedToFeedJson(allFeeds)); responseErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, responseErr.Error())
	}
}

func (config *ApiConfig) PostFeedFollowLegler(w http.ResponseWriter, r *http.Request, user database.User) {

}

func dbFeedToFeedJson(feeds []database.Feed) (feedBodies []FeedBody) {
	for _, feed := range feeds {
		var lastFetchedAtConv *time.Time = nil
		if feed.LastFetchedAt.Valid {
			lastFetchedAtConv = &feed.LastFetchedAt.Time
		}
		feedBody := FeedBody{
			ID:            feed.ID,
			CreatedAt:     feed.CreatedAt,
			UpdatedAt:     feed.UpdatedAt,
			Name:          feed.Name,
			URL:           feed.Url,
			UserID:        feed.UserID,
			LastFetchedAt: lastFetchedAtConv,
		}
		feedBodies = append(feedBodies, feedBody)
	}
	return feedBodies
}
