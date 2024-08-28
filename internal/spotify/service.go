package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
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
	client *http.Client
	token  token
}

func NewService() *service {
	return &service{
		client: &http.Client{},
		token:  token{},
	}
}

func (spotify *service) getAccessToken() (*token, error) {
	// check if token is still valid before getting a new one
	if spotify.token.AccessToken != "" && time.Now().Before(spotify.token.Expiration) {
		return &spotify.token, nil
	}

	spotifyClientID := os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	req, err := http.NewRequest(http.MethodPost, spotifyTokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(url.QueryEscape(spotifyClientID), url.QueryEscape(spotifyClientSecret))
	res, err := spotify.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	spotifyAuthResponse := SpotifyAuthResponse{}
	err = json.NewDecoder(res.Body).Decode(&spotifyAuthResponse)
	if err != nil {
		return nil, err
	}

	spotify.token = token{
		AccessToken: spotifyAuthResponse.AccessToken,
		Expiration:  time.Now().Add(time.Duration(spotifyAuthResponse.ExpiresIn) * time.Second),
	}
	return &spotify.token, nil
}

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
		return errors.New("Spotify HTTP Status: " + res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(item)
	if err != nil {
		return err
	}

	return nil
}

func (spotify *service) GetAlbums(albumIds []string) (AlbumResponse, error) {
	url := spotifyBaseURL + "/albums"
	albumResponse := AlbumResponse{}

	err := spotify.getItems(url, albumIds, &albumResponse)
	if err != nil {
		return albumResponse, err
	}

	return albumResponse, nil
}

func (spotify *service) GetTracks(trackIds []string) (TrackResponse, error) {
	url := spotifyBaseURL + "/tracks"
	trackResponse := TrackResponse{}

	err := spotify.getItems(url, trackIds, &trackResponse)
	if err != nil {
		return trackResponse, err
	}

	return trackResponse, nil
}

func (spotify *service) GetArtists(artistIds []string) (ArtistResponse, error) {
	url := spotifyBaseURL + "/artists"
	artistResponse := ArtistResponse{}

	err := spotify.getItems(url, artistIds, &artistResponse)
	if err != nil {
		return artistResponse, err
	}

	return artistResponse, nil
}

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

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := spotify.client.Do(req)
	if err != nil {
		return searchResponse, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return searchResponse, errors.New("Spotify HTTP Status: " + res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(&searchResponse)
	if err != nil {
		return searchResponse, err
	}

	return searchResponse, nil
}
