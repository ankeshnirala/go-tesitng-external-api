package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ankeshnirala/go-testing-extenal-api/external"
)

func main() {

	URL := "https://jsonplaceholder.typicode.com/posts/1"
	ext := external.New(URL, http.DefaultClient, time.Second)

	gotData, gotErr := ext.FetchData(context.Background(), "dfslg")

	fmt.Println(gotData, gotErr)

}
