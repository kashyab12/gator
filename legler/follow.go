package legler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kashyab12/gator/internal/database"
	"net/http"
	"time"
)

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FollowRequest struct {
	FeedID uuid.UUID `json:"feed_id"`
}

func (config *ApiConfig) PostFeedFollowLegler(w http.ResponseWriter, r *http.Request, user database.User) {
	followRequest := FollowRequest{}
	defer CloseIoReadCloser(r.Body)
	if decodeErr := json.NewDecoder(r.Body).Decode(&followRequest); decodeErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, decodeErr.Error())
	} else if feedFollow, feedFollowCreateErr := config.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    followRequest.FeedID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); feedFollowCreateErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, feedFollowCreateErr.Error())
	} else if responseErr := RespondWithJson(w, http.StatusCreated, FeedFollow{
		ID:        feedFollow.ID,
		FeedID:    feedFollow.FeedID,
		UserID:    feedFollow.UserID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
	}); responseErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, responseErr.Error())
	}
}

func (config *ApiConfig) DeleteFeedFollowLegler(w http.ResponseWriter, r *http.Request) {
	if feedFollowIDStr := r.PathValue("feedFollowID"); len(feedFollowIDStr) < 1 {
		_ = RespondWithError(w, http.StatusBadRequest, "invalid feedFollowID provided")
	} else if feedFollowID, convErr := uuid.Parse(feedFollowIDStr); convErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, convErr.Error())
	} else if _, deleteErr := config.DB.DeleteFeedFollow(r.Context(), feedFollowID); deleteErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, deleteErr.Error())
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (config *ApiConfig) GetAllFeedFollowsLegler(w http.ResponseWriter, r *http.Request, user database.User) {
	if allFeedFollows, fetchErr := config.DB.GetFeedFollows(r.Context(), user.ID); fetchErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, fetchErr.Error())
	} else if responseErr := RespondWithJson(w, http.StatusOK, dbFeedFollowToFeedFollowJson(allFeedFollows)); responseErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, responseErr.Error())
	}
}

func dbFeedFollowToFeedFollowJson(dbFeedFollows []database.FeedFollow) (feedFollowsJson []FeedFollow) {
	for _, dbFeedFollow := range dbFeedFollows {
		feedFollow := FeedFollow{
			ID:        dbFeedFollow.ID,
			FeedID:    dbFeedFollow.FeedID,
			UserID:    dbFeedFollow.UserID,
			CreatedAt: dbFeedFollow.CreatedAt,
			UpdatedAt: dbFeedFollow.UpdatedAt,
		}
		feedFollowsJson = append(feedFollowsJson, feedFollow)
	}
	return feedFollowsJson
}
