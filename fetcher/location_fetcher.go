package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	locationsUrlFormat = "https://api.skypicker.com/locations?type=general&key=int_id&value=%d"
	allLocationsUrl    = "https://api.skypicker.com/locations/graphql?query=%7Bdump%20%28options%3A%20%7Blocation_types%3A%20%5B%22airport%22%5D%2C%20active_only%3A%20%22true%22%7D%29%20%7Bint_id%20id%7D%7D"
)

type locationsResponse struct {
	Locations []struct {
		ID    string `json:"id"`
		IntID int    `json:"int_id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"locations"`
	Meta struct {
		Locale struct {
			Code   string `json:"code"`
			Status string `json:"status"`
		} `json:"locale"`
	} `json:"meta"`
	LastRefresh      int `json:"last_refresh"`
	ResultsRetrieved int `json:"results_retrieved"`
}

func NewLocationsFetcher() Fetcher {
	return locationFetcher{}
}

type locationFetcher struct {
}

func (f locationFetcher) Fetch(id int) (string, error) {
	resp, err := http.Get(fmt.Sprintf(locationsUrlFormat, id))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return "", readErr
	}
	r := locationsResponse{}
	if err := json.Unmarshal(body, &r); err != nil {
		return "", err
	}
	if x := len(r.Locations); x == 0 {
		return "", errors.New("location not found")
	} else if x >= 2 {
		return "", errors.New("more locations with same id found")
	}
	return r.Locations[0].ID, nil

}

func (f locationFetcher) FetchAll() (map[int]string, error) {
	resp, err := http.Get(allLocationsUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	r := locationsResponse{}
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}
	if x := len(r.Locations); x == 0 {
		return nil, errors.New("location not found")
	}
	result := make(map[int]string)
	for _, loc := range r.Locations {
		result[loc.IntID] = loc.ID
	}
	return result, nil
}
