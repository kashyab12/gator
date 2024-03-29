package legler

import (
	"context"
	"database/sql"
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

func (config *ApiConfig) FetchFeedMaster(intervalInSeconds time.Duration) {
	const queryFeedLimit = 10
	var waitGroup sync.WaitGroup
	for {
		if fetchedFeeds, fetchErr := config.DB.GetNextFeedsToFetch(context.Background(), queryFeedLimit); fetchErr != nil {
			return
		} else {
			for idx, fetchedFeed := range fetchedFeeds {
				waitGroup.Add(1)
				go config.fetchFeedInIntervalWorker(idx, &waitGroup, fetchedFeed)
			}
			waitGroup.Wait()
		}
		<-time.After(intervalInSeconds)
	}
}

func (config *ApiConfig) fetchFeedInIntervalWorker(id int, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	log.Printf("Worker %v - fetching %v\n", id, feed.Url)
	if feedData, fetchErr := fetcheed.FetchFeedData(feed.Url); fetchErr != nil {
		log.Printf("Worker %v - error fetching %v data", id, feed.Url)
	} else {
		log.Printf("Worker %v - feed Title: %v\n", id, feedData.Channel.Title)
		for _, postItem := range feedData.Channel.Items {
			var (
				pubDataTs = sql.NullTime{
					Time:  time.Time{},
					Valid: false,
				}
			)

			// Convert from string to time.Time
			if parsedPubData, parseErr := time.Parse(time.RFC1123Z, postItem.PubDate); parseErr != nil {
				log.Println(parseErr)
			} else {
				pubDataTs = sql.NullTime{
					Time:  parsedPubData,
					Valid: true,
				}
			}
			if createdPost, createErr := config.DB.CreatePost(context.Background(), database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       postItem.Title,
				Url:         postItem.Link,
				Description: postItem.Description,
				PublishedAt: pubDataTs,
				FeedID:      feed.ID,
			}); createErr != nil {
				log.Println(createErr)
			} else {
				log.Printf("Worker %v - created post %v\n", id, createdPost)
			}
			if _, updateErr := config.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
				LastFetchedAt: sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				},
				ID: feed.ID,
			}); updateErr != nil {
				log.Printf("Worker %v - error updating %v last_fetched_at and updated_at", id, feed.ID)
			}
		}
	}
}
