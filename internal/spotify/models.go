package spotify

type SpotifyAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type Album struct {
	Id          string             `json:"id"`
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
	Id         string `json:"id"`
	Album      Album  `json:"album"`
	Name       string `json:"name"`
	PreviewURL string `json:"preview_url"`
}
type TrackResponse struct {
	Tracks []Track `json:"tracks"`
}

type Artist struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Genres    []string `json:"genres"`
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
}
type SimplifiedArtist struct {
	Id   string `json:"id"`
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
