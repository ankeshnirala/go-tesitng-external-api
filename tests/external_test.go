package external_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/ankeshnirala/go-testing-extenal-api/external"
)

var (
	server *httptest.Server
	ext    external.External
)

func TestMain(m *testing.M) {
	fmt.Println("mocking server")

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFoundHandler().ServeHTTP(w, r)
	}))

	server.URL = "https://jsonplaceholder.typicode.com/posts"
	fmt.Println("mocking external", server.URL)
	ext = external.New(server.URL, http.DefaultClient, time.Second)

	fmt.Println("run tests")
	m.Run()
}

func fatal(t *testing.T, want, got interface{}) {
	t.Helper()
	t.Fatalf(`want: %v, got: %v`, want, got)
}

func TestExternal_FetchData(t *testing.T) {
	tt := []struct {
		name     string
		id       string
		wantData []*external.Data
		wantErr  error
	}{
		{
			name: "Testing jsonplaceholder success",
			id:   "1",
			wantData: []*external.Data{
				{
					Id:     1,
					UserId: 1,
					Title:  "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
					Body:   "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"},
			},
			wantErr: nil,
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotData, gotErr := ext.FetchData(context.Background(), tc.id)

			if !errors.Is(gotErr, tc.wantErr) {
				fatal(t, tc.wantErr, gotErr)
			}

			if !reflect.DeepEqual(gotData, tc.wantData) {
				fatal(t, tc.wantData, gotData)
			}
		})
	}
}

func TestFetchData_BadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFoundHandler().ServeHTTP(w, r)
	}))
	defer server.Close()

	server.URL = "https://jsonplaceholder.typicode.com/posts/1"
	ext := external.New(server.URL, http.DefaultClient, time.Second)

	tt := []struct {
		name     string
		id       string
		wantData []*external.Data
		wantErr  error
	}{
		{
			name:     "Testing jsonplaceholder failure for bad request",
			id:       "1",
			wantData: nil,
			wantErr:  fmt.Errorf("%s", http.StatusText(http.StatusBadRequest)),
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, gotErr := ext.FetchData(context.Background(), tc.id)
			if gotErr.Error() != tc.wantErr.Error() {
				fatal(t, tc.wantErr, gotErr)
			}
		})
	}
}

func TestFetchData_InternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFoundHandler().ServeHTTP(w, r)
	}))
	defer server.Close()

	server.URL = "https://jsonplacehder.typicode.com/postsss"
	ext := external.New(server.URL, http.DefaultClient, time.Second)

	tt := []struct {
		name     string
		id       string
		wantData []*external.Data
		wantErr  error
	}{
		{
			name:     "Testing jsonplaceholder failure for internal server error",
			id:       "1",
			wantData: nil,
			wantErr:  fmt.Errorf("%s", http.StatusText(http.StatusInternalServerError)),
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, gotErr := ext.FetchData(context.Background(), tc.id)
			if gotErr.Error() != tc.wantErr.Error() {
				fatal(t, tc.wantErr, gotErr)
			}
		})
	}
}
