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

func (s *service) GetTodaysQuiz() Quiz {
	if s.todaysQuiz.CreatedAt.IsZero() || time.Since(s.todaysQuiz.CreatedAt) > 24*time.Hour {
		s.generateQuiz()
	}

	return s.todaysQuiz
}

func (s *service) generateQuiz() {
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

	randomTracks, err := s.spotifyService.Search(randomWildcard, "track")
	if err != nil {
		return
	}

	// eventually should make the quiz more complex and have it's own struct
	s.todaysQuiz = Quiz{
		Track:     randomTracks.Tracks.Items[r.IntN(len(randomTracks.Tracks.Items))],
		CreatedAt: time.Now(),
	}

	log.Printf("Generated Quiz with track: %s, from: %s", s.todaysQuiz.Track.Name, s.todaysQuiz.Track.Album.Artists[0].Name)

}
