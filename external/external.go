package external

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrResponseNotOK error = errors.New("response not ok")
)

type (
	Data struct {
		UserId int    `json:"userId"`
		Id     int    `json:"id"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}

	External interface {
		FetchData(ctx context.Context, id string) ([]*Data, error)
	}

	v1 struct {
		baseURL string
		client  *http.Client
		timeout time.Duration
	}
)

func New(baseURL string, client *http.Client, timeout time.Duration) *v1 {
	return &v1{
		baseURL: baseURL,
		client:  client,
		timeout: timeout,
	}
}

func (v *v1) FetchData(ctx context.Context, id string) ([]*Data, error) {

	url := fmt.Sprintf("%s/?id=%s", v.baseURL, id)

	ctx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := v.client.Do(req)
	// handling unexpected error, like url is incorrect
	if err != nil {
		return nil, fmt.Errorf("%s", http.StatusText(http.StatusInternalServerError))
	}
	defer resp.Body.Close()

	// handling 400
	if resp.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(resp.StatusCode))
	}

	var d []*Data
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return nil, fmt.Errorf("%s", http.StatusText(http.StatusBadRequest))
	}

	return d, nil
}
