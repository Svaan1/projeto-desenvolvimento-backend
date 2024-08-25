package spotify

type SpotifyAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type Album struct {
	Name    string `json:"name"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
}
type AlbumResponse struct {
	Albums []Album `json:"albums"`
}

type Track struct {
	Name  string `json:"name"`
	Album Album  `json:"album"`
}
type TrackResponse struct {
	Tracks []Track `json:"tracks"`
}

type Artist struct {
	Name string `json:"name"`
}
type ArtistResponse struct {
	Artists []Artist `json:"artists"`
}
