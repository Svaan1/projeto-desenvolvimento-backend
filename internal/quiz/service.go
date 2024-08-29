// Package quiz provides functions to generate and manage quizzes
// using data from Spotify's API.
package quiz

import (
	"backendProject/internal/spotify"
	"log"
	"math/rand/v2"
	"time"
)

type service struct {
	todaysQuiz     Quiz
	spotifyService spotify.Service
}

func NewService(spotifyService spotify.Service) *service {
	return &service{
		spotifyService: spotifyService,
	}
}

// GetTodaysQuiz generates a new quiz if one hasn't been created today
// or returns the already generated quiz. It searches for a random song
// in Spotify's API and maps the data to a Quiz object.
//
// Returns:
//   - A Quiz object containing the generated quiz data.
//   - An error if the quiz generation fails.
func (s *service) GetTodaysQuiz() (Quiz, error) {
	// early return if quiz was already generated today
	if !s.todaysQuiz.CreatedAt.IsZero() && time.Since(s.todaysQuiz.CreatedAt) < 24*time.Hour {
		return s.todaysQuiz, nil
	}

	// since spotify doesn't have a random search,
	// we can use wildcards to search for an *almost*
	// random result

	// if not random enough, we can add more wildcards
	// or make the search more complex, since right now
	// it only gets results from the first page

	// also, it's probably a good idea to narrow results
	// based on region, otherwise we might get some
	// impossible to guess songs
	r := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
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

	randomTracks, err := s.spotifyService.Search(randomWildcard, "track")
	if err != nil {
		log.Printf("Error searching for a random song: %v", err)
		return Quiz{}, err
	}

	randomTrack := randomTracks.Tracks.Items[r.IntN(len(randomTracks.Tracks.Items))]
	artists, err := s.spotifyService.GetArtists([]string{randomTrack.Album.Artists[0].Id})
	if err != nil {
		log.Printf("Error getting artists from random song: %v", err)
		return Quiz{}, err
	}

	s.todaysQuiz = buildQuiz(randomTrack, artists.Artists)
	log.Printf("Generated Quiz with track: %s, from: %v", s.todaysQuiz.Track.Name, s.todaysQuiz.Artists)

	return s.todaysQuiz, nil
}

// buildQuiz creates a Quiz object from a given Spotify track and a list
// of Spotify artists. It maps the provided data to the appropriate quiz models
// and returns the constructed Quiz.
//
// Parameters:
//   - track: A Spotify track to be used in the quiz.
//   - artists: A slice of Spotify artists to be included in the quiz.
//
// Returns:
//   - A Quiz object containing the mapped track, album, and artist data.
func buildQuiz(track spotify.Track, artists []spotify.Artist) Quiz {
	return Quiz{
		Artists:   mapArtists(artists),
		Album:     mapAlbum(track.Album),
		Track:     mapTrack(track),
		CreatedAt: time.Now(),
	}
}

// mapArtists converts a slice of Spotify artists to a slice of quizArtist objects.
//
// Parameters:
//   - artists: A slice of Spotify artists to be mapped.
//
// Returns:
//   - A slice of quizArtist objects containing the mapped artist data.
func mapArtists(artists []spotify.Artist) []quizArtist {
	mappedArtists := make([]quizArtist, len(artists))
	for i, artist := range artists {
		mappedArtists[i] = quizArtist{
			Id:     artist.Id,
			Name:   artist.Name,
			Genres: artist.Genres,
		}
	}
	return mappedArtists
}

// mapAlbum converts a Spotify album to a quizAlbum object.
//
// Parameters:
//   - album: A Spotify album to be mapped.
//
// Returns:
//   - A quizAlbum object containing the mapped album data.
func mapAlbum(album spotify.Album) quizAlbum {
	return quizAlbum{
		Id:          album.Id,
		Name:        album.Name,
		Image:       album.Images[0].URL,
		ReleaseDate: album.ReleaseDate,
	}
}

// mapTrack converts a Spotify track to a quizSong object.
//
// Parameters:
//   - track: A Spotify track to be mapped.
//
// Returns:
//   - A quizSong object containing the mapped track data.
func mapTrack(track spotify.Track) quizSong {
	return quizSong{
		Id:           track.Id,
		Name:         track.Name,
		AudioPreview: track.PreviewURL,
	}
}
