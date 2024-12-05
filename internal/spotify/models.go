package spotify

import (
	"strconv"
	"time"
)

type Service interface {
	GetAlbums(albumIds []string) (AlbumResponse, error)
	GetTracks(trackIds []string) (TrackResponse, error)
	GetArtists(artistIds []string) (ArtistResponse, error)
	Search(query, queryType string) (SearchResponse, error)
	RandomSearch(queryType string) (SearchResponse, error)
	GetRecommendations(seedArtists, seedGenres, seedTracks []string, popularity int) (RecommendationsResponse, error)
}

type Token struct {
	AccessToken string
	Expiration  time.Time
}

type SpotifyAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type Album struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Artists     []SimplifiedArtist `json:"artists"`
	ReleaseDate string             `json:"release_date"`
	Images      []struct {
		URL string `json:"url"`
	} `json:"images"`
}
type AlbumResponse struct {
	Albums []Album `json:"albums"`
}

type Track struct {
	ID         string `json:"id"`
	Album      Album  `json:"album"`
	Name       string `json:"name"`
	PreviewURL string `json:"preview_url"`
}
type TrackResponse struct {
	Tracks []Track `json:"tracks"`
}

type Artist struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Genres    []string `json:"genres"`
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
}
type SimplifiedArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type ArtistResponse struct {
	Artists []Artist `json:"artists"`
}

type SearchResponse struct {
	Albums struct {
		Items []Album `json:"items"`
	} `json:"albums"`
	Tracks struct {
		Items []Track `json:"items"`
	} `json:"tracks"`
	Artists struct {
		Items []Artist `json:"items"`
	} `json:"artists"`
}

type RecommendationsResponse struct {
	Tracks []Track `json:"tracks"`
}

func (r RecommendationsResponse) String() string {
	var result string
	for _, track := range r.Tracks {
		result += track.String() + " by "
		for i, artist := range track.Album.Artists {
			if i > 0 {
				result += ", "
			}
			result += artist.Name
		}
		result += "\n"
	}
	return result
}

func (a Album) String() string {
	var result string
	result += "Album: " + a.Name + " by "
	for i, artist := range a.Artists {
		if i > 0 {
			result += ", "
		}
		result += artist.Name
	}
	result += " released on " + a.ReleaseDate + "\n"
	return result
}

func (t Track) String() string {
	var result string
	result += "Track: " + t.Name + " from album " + t.Album.Name + " by "
	for i, artist := range t.Album.Artists {
		if i > 0 {
			result += ", "
		}
		result += artist.Name
	}
	result += "\n"
	return result
}

func (a Artist) String() string {
	var result string
	result += "Artist: " + a.Name + " with genres "
	for i, genre := range a.Genres {
		if i > 0 {
			result += ", "
		}
		result += genre
	}
	result += " and " + strconv.Itoa(a.Followers.Total) + " followers\n"
	return result
}
