package spotify

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var spotifyService *service

func init() {
	godotenv.Load("../../.env")
}

func TestGetItems(t *testing.T) {
	spotifyService = NewService(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))

	type args struct {
		ids      []string
		itemType string
	}
	testCases := []struct {
		given    args
		expected string
	}{
		{args{[]string{"4LH4d3cOWNNsVw41Gqt2kv"}, "album"}, "The Dark Side of the Moon"},
		{args{[]string{"6mFkJmJqdDVQ1REhVfGgd1"}, "track"}, "Wish You Were Here"},
		{args{[]string{"0k17h0D3J5VfsdmQ1iZtE9"}, "artist"}, "Pink Floyd"},
	}
	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			switch tc.given.itemType {
			case "album":
				var albumResponse AlbumResponse
				err := spotifyService.getItems(spotifyBaseURL+"/albums", tc.given.ids, &albumResponse)
				if err != nil {
					t.Errorf("Error getting album: %v", err)
					return
				}
				if albumResponse.Albums[0].Name != tc.expected {
					t.Errorf("Expected album name to be %s, got %s", tc.expected, albumResponse.Albums[0].Name)
				}
			case "track":
				var trackResponse TrackResponse
				err := spotifyService.getItems(spotifyBaseURL+"/tracks", tc.given.ids, &trackResponse)
				if err != nil {
					t.Errorf("Error getting track: %v", err)
					return
				}
				if trackResponse.Tracks[0].Name != tc.expected {
					t.Errorf("Expected track name to be %s, got %s", tc.expected, trackResponse.Tracks[0].Name)
				}
			case "artist":
				var artistResponse ArtistResponse
				err := spotifyService.getItems(spotifyBaseURL+"/artists", tc.given.ids, &artistResponse)
				if err != nil {
					t.Errorf("Error getting artist: %v", err)
					return
				}
				if artistResponse.Artists[0].Name != tc.expected {
					t.Errorf("Expected artist name to be %s, got %s", tc.expected, artistResponse.Artists[0].Name)
				}
			}
		})
	}
}

func TestGetItemsError(t *testing.T) {
	spotifyService = NewService(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))

	type args struct {
		ids      []string
		itemType string
	}
	testCases := []struct {
		given    args
		expected string
	}{
		{args{[]string{"0000000000000000000"}, "album"}, "album"},
		{args{[]string{"%1%4$##@*&+(());'//"}, "track"}, "track"},
		{args{[]string{"*(!JM@#*()(!)))!!!#"}, "artist"}, "artist"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			switch tc.given.itemType {
			case "album":
				var albumResponse AlbumResponse
				err := spotifyService.getItems(spotifyBaseURL+"/albums", tc.given.ids, &albumResponse)
				if err == nil {
					t.Errorf("Expected error getting album, got nil")
				}
			case "track":
				var trackResponse TrackResponse
				err := spotifyService.getItems(spotifyBaseURL+"/tracks", tc.given.ids, &trackResponse)
				if err == nil {
					t.Errorf("Expected error getting track, got nil")
				}
			case "artist":
				var artistResponse ArtistResponse
				err := spotifyService.getItems(spotifyBaseURL+"/artists", tc.given.ids, &artistResponse)
				if err == nil {
					t.Errorf("Expected error getting artist, got nil")
				}
			}
		})
	}
}

func TestSearch(t *testing.T) {
	spotifyService = NewService(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))

	type args struct {
		query     string
		queryType string
	}
	testCases := []struct {
		given    args
		expected string
	}{
		{args{"The Dark Side of the Moon", "album"}, "The Dark Side of the Moon"},
		{args{"Wish You Were Here", "track"}, "Wish You Were Here"},
		{args{"Pink Floyd", "artist"}, "Pink Floyd"},
	}
	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			switch tc.given.queryType {
			case "album":
				albumResponse, err := spotifyService.Search(tc.given.query, "album")
				if err != nil {
					t.Errorf("Error searching for album: %v", err)
					return
				}
				if albumResponse.Albums.Items[0].Name != tc.expected {
					t.Errorf("Expected album name to be %s, got %s", tc.expected, albumResponse.Albums.Items[0].Name)
				}
			case "track":
				trackResponse, err := spotifyService.Search(tc.given.query, "track")
				if err != nil {
					t.Errorf("Error searching for track: %v", err)
					return
				}
				if trackResponse.Tracks.Items[0].Name != tc.expected {
					t.Errorf("Expected track name to be %s, got %s", tc.expected, trackResponse.Tracks.Items[0].Name)
				}
			case "artist":
				artistResponse, err := spotifyService.Search(tc.given.query, "artist")
				if err != nil {
					t.Errorf("Error searching for artist: %v", err)
					return
				}
				if artistResponse.Artists.Items[0].Name != tc.expected {
					t.Errorf("Expected artist name to be %s, got %s", tc.expected, artistResponse.Artists.Items[0].Name)
				}
			}
		})
	}
}

func TestGetAlbums(t *testing.T) {
	spotifyService = NewService(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))

	testCases := []struct {
		given    []string
		expected string
	}{
		{[]string{"4LH4d3cOWNNsVw41Gqt2kv"}, "The Dark Side of the Moon"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			albumResponse, err := spotifyService.GetAlbums(tc.given)
			if err != nil {
				t.Errorf("Error getting album: %v", err)
				return
			}
			if albumResponse.Albums[0].Name != tc.expected {
				t.Errorf("Expected album name to be %s, got %s", tc.expected, albumResponse.Albums[0].Name)
			}
		})
	}
}

func TestGetTracks(t *testing.T) {
	spotifyService = NewService(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))

	testCases := []struct {
		given    []string
		expected string
	}{
		{[]string{"6mFkJmJqdDVQ1REhVfGgd1"}, "Wish You Were Here"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			trackResponse, err := spotifyService.GetTracks(tc.given)
			if err != nil {
				t.Errorf("Error getting track: %v", err)
				return
			}
			if trackResponse.Tracks[0].Name != tc.expected {
				t.Errorf("Expected track name to be %s, got %s", tc.expected, trackResponse.Tracks[0].Name)
			}
		})
	}
}

func TestGetArtists(t *testing.T) {
	spotifyService = NewService(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))

	testCases := []struct {
		given    []string
		expected string
	}{
		{[]string{"0k17h0D3J5VfsdmQ1iZtE9"}, "Pink Floyd"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			artistResponse, err := spotifyService.GetArtists(tc.given)
			if err != nil {
				t.Errorf("Error getting artist: %v", err)
				return
			}
			if artistResponse.Artists[0].Name != tc.expected {
				t.Errorf("Expected artist name to be %s, got %s", tc.expected, artistResponse.Artists[0].Name)
			}
		})
	}
}

func TestGetItemsWithoutCredentials(t *testing.T) {
	spotifyService := NewService("", "")

	var albumResponse AlbumResponse
	err := spotifyService.getItems(spotifyBaseURL+"/albums", []string{"4LH4d3cOWNNsVw41Gqt2kv"}, &albumResponse)
	if err == nil {
		t.Errorf("Expected error getting album, got nil")
	}
}

func TestSearchWithoutCredentials(t *testing.T) {
	spotifyService := NewService("", "")

	_, err := spotifyService.Search("The Dark Side of the Moon", "album")
	if err == nil {
		t.Errorf("Expected error searching for album, got nil")
	}
}

func TestRandomSearch(t *testing.T) {
	spotifyService = NewService(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))

	_, err := spotifyService.RandomSearch("track")
	if err != nil {
		t.Errorf("Error searching for random track: %v", err)
	}
}

func TestGetRecommendations(t *testing.T) {
	spotifyService = NewService(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))

	_, err := spotifyService.GetRecommendations([]string{"0k17h0D3J5VfsdmQ1iZtE9"}, []string{"rock"}, []string{"6mFkJmJqdDVQ1REhVfGgd1"}, 80)
	if err != nil {
		t.Errorf("Error getting recommendations: %v", err)
	}
}
