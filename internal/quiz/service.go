// Package quiz provides functions to generate and manage quizzes
// using data from Spotify's API.
package quiz

import (
	"backendProject/internal/spotify"
	"context"
	"log"
	"math/rand/v2"
	"time"
)

type service struct {
	repository     *Repository
	spotifyService spotify.Service
}

func NewService(repository *Repository, spotifyService spotify.Service) *service {
	return &service{
		spotifyService: spotifyService,
		repository:     repository,
	}
}

// GetTodaysQuiz generates a new quiz if one hasn't been created today
// or returns the already generated quiz. It searches for a random song
// in Spotify's API and maps the data to a Quiz object.
//
// Returns:
//   - A Quiz object containing the generated quiz data.
//   - An error if the quiz generation fails.
func (s *service) GetTodaysQuiz(ctx context.Context) (Quiz, error) {
	// early return if quiz was already generated today
	todaysQuiz, err := s.repository.GetQuiz(ctx, "quiz")
	if err != nil {
		log.Printf("Error getting today's quiz: %v", err)
		return Quiz{}, err
	}

	if !todaysQuiz.CreatedAt.IsZero() && time.Since(todaysQuiz.CreatedAt) < 24*time.Hour {
		log.Printf("Returning already generated quiz created at: %v", todaysQuiz.CreatedAt)
		return todaysQuiz, nil
	}

	randomTracks, err := s.spotifyService.RandomSearch("track")
	if err != nil {
		log.Printf("Error searching for a random song: %v", err)
		return Quiz{}, err
	}

	r := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	randomTrack := randomTracks.Tracks.Items[r.IntN(len(randomTracks.Tracks.Items))]

	artistIDs := make([]string, 5)
	for i, artist := range randomTrack.Album.Artists {
		if i == 5 {
			break
		}

		artistIDs[i] = artist.ID
	}

	// this prevents the quiz from being generated with a track that doesn't have a preview
	// right now it may loop indefinitely, making the spotify return 429, and this will basically
	// break the entire api. This is a temporary solution, and it should be fixed ASAP!!.
	var track spotify.Track
	for i := 0; i < 10; i++ {
		recommendedTracks, err := s.spotifyService.GetRecommendations(artistIDs, nil, []string{randomTrack.ID}, 80)
		if err != nil {
			log.Printf("Error getting recommendations from random song: %v", err)
			return Quiz{}, err
		}

		if len(recommendedTracks.Tracks) == 0 {
			continue
		}

		for j := 0; j < 10; j++ {
			track = recommendedTracks.Tracks[r.IntN(len(recommendedTracks.Tracks))]
			if track.PreviewURL != "" {
				break
			}
		}
		// TODO: get the property Retry-After from the response header and wait that amount of time
		time.Sleep(1 * time.Second)
	}

	recommmentedArtistIDs := make([]string, 5)
	for i, artist := range track.Album.Artists {
		if i == 5 {
			break
		}

		recommmentedArtistIDs[i] = artist.ID
	}

	artists, err := s.spotifyService.GetArtists(recommmentedArtistIDs)
	if err != nil {
		log.Printf("Error getting artists from random song: %v", err)
		return Quiz{}, err
	}

	todaysQuiz = buildQuiz(track, artists.Artists)
	err = s.repository.SetQuiz(ctx, "quiz", todaysQuiz)
	if err != nil {
		log.Printf("Error setting today's quiz: %v", err)
		return Quiz{}, err
	}

	artistNames := make([]string, len(todaysQuiz.Artists))
	for i, artist := range todaysQuiz.Artists {
		artistNames[i] = artist.Name
	}
	log.Printf("Generated Quiz with Track: %s, Album: %s, Artists: %v", todaysQuiz.Track.Name, todaysQuiz.Album.Name, artistNames)

	return todaysQuiz, nil
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
			ID:     artist.ID,
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
		ID:          album.ID,
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
		ID:           track.ID,
		Name:         track.Name,
		AudioPreview: track.PreviewURL,
	}
}
