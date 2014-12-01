// Package ohmy pulls data from the ohmyrockness.com API
package ohmy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Shows is an array of shows, the response we get from the ohmyrockness
// API
type Shows []*Show

// Show is a list of bands, a venue and some other info. Not everything
// returned in the API is included here.
type Show struct {
	Bands  []*Band `json:"cached_bands"`
	Venue  *Venue
	Starts *time.Time `json:"starts_at"`
}

// Venue is where the show takes place.
type Venue struct {
	Address   string `json:"full_address"`
	Latitude  string
	Longitude string
	Name      string
	Slug      string
}

// Band is one of the bands playing at the show
type Band struct {
	Name string
	Slug string
}

const (
	// base URL of the API and the full URL of the index page. We use
	// the index page to get the CSRF token.
	base = "http://www.ohmyrockness.com/"

	// api is the endpoint of the shows API, relative to base.
	api = "api/shows.json"

	// regioned value for NY
	regioned = 1

	// maxPer is the max number of records per page to request
	maxPer = 50

	// authToken is a value hidden in the javascript code. TODO: get this
	// dynamically.
	authToken = "3b35f8a73dabd5f14b1cac167a14c1f6"
)

var (
	// NoCSRF is a generic error for any problem getting the CSRF token
	NoCSRF = errors.New("problem getting CSRF token")
)

// getIndexData gets the CSRF token and cookies from the index page
func getIndexData() (csrf string, cookies []*http.Cookie, err error) {

	// Get the content via HTTP
	response, err := http.Get(base)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()

	// Parse as a goquery document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	// Find the token tag
	s := doc.Find(`meta[name="csrf-token"]`).First()
	if s == nil {
		err = NoCSRF
		log.Println(err)
		return
	}

	// Extract the value from the tag
	csrf, exists := s.Attr("content")
	if !exists {
		err = NoCSRF
		log.Println(err)
		return
	}

	cookies = response.Cookies()

	// Success! We have the CSRF and cookies.
	return
}

// apiURL returns the full HTTP URL to given the page, etc.
func apiURL(page, per int) string {
	// Combine base and api path to get the URL (minus the query string)
	u := fmt.Sprint(base, api)

	// Create the query string object
	v := url.Values{}
	v.Set("index", "true")
	v.Set("page", fmt.Sprint(page))
	v.Set("per", fmt.Sprint(per))
	v.Set("regioned", fmt.Sprint(regioned))

	// Return the full URL
	fullURL := fmt.Sprint(u, "?", v.Encode())
	return fullURL
}

func callAPI(csrf string, cookies []*http.Cookie, page, per int) (shows Shows, err error) {
	// Create a client so we can modify headers of the request
	client := &http.Client{}

	// Create a request to modify its headers
	request, err := http.NewRequest("GET", apiURL(page, per), nil)

	// Add authorization to request headers
	request.Header.Add("X-CSRF-Token", csrf)
	request.Header.Add("Authorization", fmt.Sprintf(
		`Token token="%s"`, authToken,
	))

	// Add cookies from the index page to request
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}

	// Make the API call
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()

	// Read the body of the response
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	// Encode body into an array of shows
	err = json.Unmarshal(b, &shows)
	if err != nil {
		log.Println(err)
		return
	}

	// Success! We have the API response of shows
	return
}

// GetShows will get at most n number of shows (if less than n available,
// will give the most available)
func GetShows(n int) (allShows Shows, err error) {
	// Get the data we need to call the API
	csrf, cookies, err := getIndexData()
	if err != nil {
		return
	}

	// List of all shows across different API calls
	allShows = Shows{}

	// Start with page one
	page := 1

	// Save our original n for limit later
	origN := n

	// Loop until no more results to get
	for n > 0 {

		// Determine number of records to request
		var per int
		if n > maxPer {
			per = maxPer
		} else {
			per = n
		}

		// Get shows for this page
		shows, callErr := callAPI(csrf, cookies, page, per)
		if callErr != nil {
			err = callErr
			return
		}

		// Append to all of our shows
		allShows = append(allShows, shows...)

		// Increment/decrement for next loop
		page++
		n = n - per
	}

	// It seems that the API doesn't actually listen to the per value, so
	// limit the returned shows here.
	allShows = allShows[0:origN]

	return
}
