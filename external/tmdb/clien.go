package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiKey, baseURL string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// TMDB filmo struktūra
type TMDBMovie struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title"`
	Overview      string  `json:"overview"`
	ReleaseDate   string  `json:"release_date"`
	PosterPath    string  `json:"poster_path"`
	BackdropPath  string  `json:"backdrop_path"`
	VoteAverage   float64 `json:"vote_average"`
	VoteCount     int     `json:"vote_count"`
	Popularity    float64 `json:"popularity"`
	Runtime       int     `json:"runtime"`
	Status        string  `json:"status"`
	ImdbID        string  `json:"imdb_id"`
	Genres        []Genre `json:"genres"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Paieška filmų pagal pavadinimą
func (c *Client) SearchMovies(ctx context.Context, query string, page int) (*SearchResponse, error) {
	endpoint := fmt.Sprintf("%s/search/movie", c.baseURL)

	params := url.Values{}
	params.Add("api_key", c.apiKey)
	params.Add("query", query)
	params.Add("page", strconv.Itoa(page))
	params.Add("language", "en-US")
	params.Add("include_adult", "false")

	url := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %s", resp.Status)
	}

	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Gauti filmo detales pagal TMDB ID
func (c *Client) GetMovieDetails(ctx context.Context, tmdbID int) (*TMDBMovie, error) {
	endpoint := fmt.Sprintf("%s/movie/%d", c.baseURL, tmdbID)

	params := url.Values{}
	params.Add("api_key", c.apiKey)
	params.Add("language", "en-US")
	params.Add("append_to_response", "credits") // gauti aktorius, režisierius

	url := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %s", resp.Status)
	}

	var movie TMDBMovie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, err
	}

	return &movie, nil
}

// Populiariausi filmai
func (c *Client) GetPopularMovies(ctx context.Context, page int) (*SearchResponse, error) {
	endpoint := fmt.Sprintf("%s/movie/popular", c.baseURL)

	params := url.Values{}
	params.Add("api_key", c.apiKey)
	params.Add("page", strconv.Itoa(page))
	params.Add("language", "en-US")

	url := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %s", resp.Status)
	}

	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type SearchResponse struct {
	Page         int         `json:"page"`
	Results      []TMDBMovie `json:"results"`
	TotalPages   int         `json:"total_pages"`
	TotalResults int         `json:"total_results"`
}
