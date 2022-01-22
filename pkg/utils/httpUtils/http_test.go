package httpUtils

import (
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

const (
	contentTypeJson = "application/json"
)

func TestSendRequest(t *testing.T) {
	request := Request{
		Url:         "https://catfact.ninja/fact",
		Method:      http.MethodGet,
		ContentType: contentTypeJson,
	}
	type response struct {
		Fact   string `json:"fact"`
		Length int    `json:"length"`
	}
	var responseBody response
	code, err := SendRequest(request, &responseBody, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, code)
	require.NotEmpty(t, responseBody.Fact)
	require.NotEqual(t, 0, responseBody.Length)
}

func TestSendRequest_QueryParameters(t *testing.T) {
	request := Request{
		Url:    "https://api.agify.io",
		Method: http.MethodGet,
		QueryParameters: map[string]string{
			"name": "toto",
		},
		ContentType: contentTypeJson,
	}
	type response struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Count int    `json:"count"`
	}
	var responseBody response
	code, err := SendRequest(request, &responseBody, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, code)
	require.NotEmpty(t, responseBody.Name)
	require.NotEqual(t, 0, responseBody.Age)
	require.NotEqual(t, 0, responseBody.Count)
}

func TestSendRequest_HeaderParameters(t *testing.T) {
	request := Request{
		Url:    "https://api.themoviedb.org/3/movie/155",
		Method: http.MethodGet,
		HeadersParameters: map[string]string{
			"Authorization": "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIyNDAxNjYxOWFiMjNiMDYzNjMzYzgwZTY4MzFlN2NjYyIsInN1YiI6IjYwOGI3ZjZlOGM0MGY3MDA1N2U3ZDg4MCIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.TqHh6OC7IZ0s7err6njtR054Pi87kG6UaaER5WL04k0",
		},
		QueryParameters: map[string]string{
			"language": "fr-FR",
		},
		ContentType: contentTypeJson,
	}
	type response struct {
		BackdropPath string `json:"backdrop_path"`
		Budget       int    `json:"budget"`
		Genres       []struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"genres"`
		Homepage         string  `json:"homepage"`
		Id               int     `json:"id"`
		ImdbId           string  `json:"imdb_id"`
		OriginalLanguage string  `json:"original_language"`
		OriginalTitle    string  `json:"original_title"`
		Overview         string  `json:"overview"`
		Popularity       float64 `json:"popularity"`
		PosterPath       string  `json:"poster_path"`
		ReleaseDate      string  `json:"release_date"`
		Revenue          int     `json:"revenue"`
		Runtime          int     `json:"runtime"`
	}
	var responseBody response
	code, err := SendRequest(request, &responseBody, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, "The Dark Knight", responseBody.OriginalTitle)
}

func TestSendRequest_Body(t *testing.T) {
	type user struct {
		Name string `json:"name"`
		Job  string `json:"job"`
	}
	fakeUser := user{
		Name: fake.UserName(),
		Job:  fake.Words(),
	}
	request := Request{
		Url:         "https://reqres.in/api/users",
		Method:      http.MethodPost,
		ContentType: contentTypeJson,
		Body:        fakeUser,
	}
	type response struct {
		Name      string    `json:"name"`
		Job       string    `json:"job"`
		Id        string    `json:"id"`
		CreatedAt time.Time `json:"createdAt"`
	}
	var responseBody response
	code, err := SendRequest(request, &responseBody, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, code)
	require.Equal(t, fakeUser.Name, responseBody.Name)
	require.Equal(t, fakeUser.Job, responseBody.Job)
	require.NotEqual(t, "0", responseBody.Id)
	require.True(t, responseBody.CreatedAt.Local().Before(time.Now().Local().Add(1*time.Second)),
		"Creation time %s should be before now time %s", responseBody.CreatedAt.Local(),
		time.Now().Local().Add(1*time.Second))
}

func TestSendRequest_HandleBadHttpStatusCodeWithExpectedResponse(t *testing.T) {
	request := Request{
		Url:         "http://blabla.customers.isi.nc",
		Method:      http.MethodPost,
		ContentType: contentTypeJson,
	}
	type response struct {
		Something string `json:"something"`
	}
	var responseBody response
	code, err := SendRequest(request, &responseBody, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, code)
}

// TestSendRequest_Timeout will send a request to website that will wait X sec before sending response.
// In this way, we can check if timeout is working.
func TestSendRequest_Timeout(t *testing.T) {
	request := Request{
		Url:         "http://httpstat.us/200?sleep=5000",
		Method:      http.MethodGet,
		ContentType: contentTypeJson,
	}

	// Try with deadline exceed (timeout < 5s)
	customClient := http.Client{Timeout: 2 * time.Second}

	_, err := SendRequest(request, nil, &customClient)
	require.Error(t, err)

	// Try with deadline exceed (timeout > 5s)
	customClient = http.Client{Timeout: 15 * time.Second}

	_, err = SendRequest(request, nil, &customClient)
	require.NoError(t, err)
}
