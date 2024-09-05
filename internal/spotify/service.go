package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	spotifyBaseURL  = "https://api.spotify.com/v1"
	spotifyTokenURL = "https://accounts.spotify.com/api/token"
)

type token struct {
	AccessToken string
	Expiration  time.Time
}
type service struct {
	client              *http.Client
	token               token
	spotifyClientID     string
	spotifyClientSecret string
}

func NewService(spotifyClientID, spotifyClientSecret string) *service {
	return &service{
		client:              &http.Client{},
		token:               token{},
		spotifyClientID:     spotifyClientID,
		spotifyClientSecret: spotifyClientSecret,
	}
}

// getAccessToken retrieves a new access token from Spotify's API.
// If the current token is still valid, it returns the current token.
//
// Returns:
//   - A token object containing the access token and its expiration time.
//   - An error if the request fails.
func (s *service) getAccessToken() (*token, error) {
	// check if token is still valid before getting a new one
	if s.token.AccessToken != "" && time.Now().Before(s.token.Expiration) {
		return &s.token, nil
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	req, err := http.NewRequest(http.MethodPost, spotifyTokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(url.QueryEscape(s.spotifyClientID), url.QueryEscape(s.spotifyClientSecret))
	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	spotifyAuthResponse := SpotifyAuthResponse{}
	err = json.NewDecoder(res.Body).Decode(&spotifyAuthResponse)
	if err != nil {
		return nil, err
	}

	s.token = token{
		AccessToken: spotifyAuthResponse.AccessToken,
		Expiration:  time.Now().Add(time.Duration(spotifyAuthResponse.ExpiresIn) * time.Second),
	}
	log.Println("New Spotify access token generated")
	return &s.token, nil
}

// getItems retrieves items from Spotify's API based on the given URL and IDs.
// the items can be of type Album, Track or Artist.
//
// Parameters:
//   - url: The URL to send the request to.
//   - ids: A slice of IDs to retrieve from the API.
//   - item: A pointer to the item struct to decode the response into. (Album, Track or Artist)
//
// Returns:
//   - An error if the request or data parsing fails.
func (spotify *service) getItems(url string, ids []string, item interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	token, err := spotify.getAccessToken()
	if err != nil {
		return err
	}

	params := req.URL.Query()
	params.Set("ids", strings.Join(ids, ","))
	req.URL.RawQuery = params.Encode()

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := spotify.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return errors.New("Spotify HTTP Status: " + res.Status)
		}

		return errors.New("Spotify HTTP Status: " + res.Status + "\n" + string(body))
	}

	err = json.NewDecoder(res.Body).Decode(item)
	if err != nil {
		return err
	}

	return nil
}

// GetAlbums retrieves albums from Spotify's API based on the given album IDs.
//
// Parameters:
//   - albumIds: A slice of album IDs to retrieve from the API.
//
// Returns:
//   - An AlbumResponse object containing the retrieved albums.
//   - An error if the request or data parsing fails.
func (spotify *service) GetAlbums(albumIds []string) (AlbumResponse, error) {
	url := spotifyBaseURL + "/albums"
	albumResponse := AlbumResponse{}

	err := spotify.getItems(url, albumIds, &albumResponse)
	if err != nil {
		return albumResponse, err
	}

	return albumResponse, nil
}

// GetTracks retrieves tracks from Spotify's API based on the given track IDs.
//
// Parameters:
//   - trackIds: A slice of track IDs to retrieve from the API.
//
// Returns:
//   - A TrackResponse object containing the retrieved tracks.
//   - An error if the request or data parsing fails.
func (spotify *service) GetTracks(trackIds []string) (TrackResponse, error) {
	url := spotifyBaseURL + "/tracks"
	trackResponse := TrackResponse{}

	err := spotify.getItems(url, trackIds, &trackResponse)
	if err != nil {
		return trackResponse, err
	}

	return trackResponse, nil
}

// GetArtists retrieves artists from Spotify's API based on the given artist IDs.
//
// Parameters:
//   - artistIds: A slice of artist IDs to retrieve from the API.
//
// Returns:
//   - An ArtistResponse object containing the retrieved artists.
//   - An error if the request or data parsing fails.
func (spotify *service) GetArtists(artistIds []string) (ArtistResponse, error) {
	url := spotifyBaseURL + "/artists"
	artistResponse := ArtistResponse{}

	err := spotify.getItems(url, artistIds, &artistResponse)
	if err != nil {
		return artistResponse, err
	}

	return artistResponse, nil
}

// Search retrieves search results from Spotify's API based on the given query and query type.
//
// Parameters:
//
//   - query: The search query to retrieve from the API.
//
//   - queryType: The type of search query to perform. (track, album, artist)
//
// Returns:
//   - A SearchResponse object containing the search results.
//   - An error if the request or data parsing fails.
func (spotify *service) Search(query, queryType string) (SearchResponse, error) {
	url := spotifyBaseURL + "/search"
	var searchResponse SearchResponse

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return searchResponse, err
	}

	token, err := spotify.getAccessToken()
	if err != nil {
		return searchResponse, err
	}

	params := req.URL.Query()
	params.Set("q", query)
	params.Set("type", queryType)
	req.URL.RawQuery = params.Encode()

	log.Println("Searching for:", query)

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := spotify.client.Do(req)
	if err != nil {
		return searchResponse, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return searchResponse, errors.New("Spotify HTTP Status: " + res.Status)
		}

		return searchResponse, errors.New("Spotify HTTP Status: " + res.Status + "\n" + string(body))
	}

	err = json.NewDecoder(res.Body).Decode(&searchResponse)
	if err != nil {
		return searchResponse, err
	}

	return searchResponse, nil
}

// RandomSearch retrieves search results from Spotify's API based on a random query
//
// Parameters:
//   - queryType: The type of search query to perform. (track, album, artist)
//
// Returns:
//   - A SearchResponse object containing the search results.
//   - An error if the request or data parsing fails.
func (s *service) RandomSearch(queryType string) (SearchResponse, error) {
	r := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

	// since spotify doesn't have a random search,
	// we can use wildcards to search for an *almost*
	// random result

	// if not random enough, we can add more wildcards
	// or make the search more complex, since right now
	// it only gets results from the first page

	// also, it's probably a good idea to narrow results
	// based on region, otherwise we might get some
	// impossible to guess songs

	letters := "abcdefghijklmnopqrstuvwxyz1234567890"
	var wildcards []string

	// %aa, %aa%, aa%, %ab, %ab%, ab% ...
	for i := 0; i < len(letters); i++ {
		for j := 0; j < len(letters); j++ {
			combination := string(letters[i]) + string(letters[j])
			wildcards = append(wildcards, "%"+combination, "%"+combination+"%", combination+"%")
		}
	}
	randomWildcard := wildcards[r.IntN(len(wildcards))]

	return s.Search(randomWildcard, queryType)
}
