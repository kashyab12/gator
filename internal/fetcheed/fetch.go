package fetcheed

import (
	"encoding/xml"
	"net/http"
)

func FetchFeedData(url string) (any, error) {
	feedData := struct {
		any
	}{}
	if resp, reqErr := http.Get(url); reqErr != nil {
		return feedData, reqErr
	} else if decodingErr := xml.NewDecoder(resp.Body).Decode(&feedData); decodingErr != nil {
		return feedData, decodingErr
	}
	return feedData, nil
}
