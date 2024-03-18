package legler

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kashyab12/gator/internal/database"
	"github.com/kashyab12/gator/internal/fetcheed"
	"log"
	"net/http"
	"sync"
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
	defer CloseIoReadCloser(r.Body)
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
	} else if feedFollow, feedFollowCreateErr := config.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    newFeed.ID,
		UserID:    newFeed.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); feedFollowCreateErr != nil {
		_ = RespondWithError(w, http.StatusInternalServerError, feedFollowCreateErr.Error())
	} else if responseErr := RespondWithJson(w, http.StatusCreated, map[string]any{"feed": &FeedBody{
		ID:            newFeed.ID,
		CreatedAt:     newFeed.CreatedAt,
		UpdatedAt:     newFeed.UpdatedAt,
		Name:          newFeed.Name,
		URL:           newFeed.Url,
		UserID:        newFeed.UserID,
		LastFetchedAt: nil,
	}, "feed_follow": feedFollow}); responseErr != nil {
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

func (config *ApiConfig) FetchFeedMaster(intervalInSeconds time.Duration) error {
	const queryFeedLimit = 10
	var waitGroup sync.WaitGroup
	for {
		<-time.After(time.Second * intervalInSeconds)
		if fetchedFeeds, fetchErr := config.DB.GetNextFeedsToFetch(context.Background(), queryFeedLimit); fetchErr != nil {
			return fetchErr
		} else {
			for idx, fetchedFeed := range fetchedFeeds {
				waitGroup.Add(1)
				go fetchFeedInIntervalWorker(idx, &waitGroup, fetchedFeed)
			}
			waitGroup.Wait()
		}
	}
}

func fetchFeedInIntervalWorker(id int, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	log.Printf("Worker %v - fetching %v\n", id, feed.Url)
	if feedData, fetchErr := fetcheed.FetchFeedData(feed.Url); fetchErr != nil {
		log.Printf("Worker %v - error fetching %v data", id, feed.Url)
	} else {
		log.Printf("Worker %v - feed Title: %v\n", id, feedData.Channel.Title)
	}
}
