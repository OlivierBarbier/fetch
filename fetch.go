package main

import (
    // "fmt"
    "net/http"
)

var urls = []string{
    "http://app1.tillersystems.com/api",
    // "http://app2.tillersystems.com/api",    
}

type HttpResponse struct {
    url      string
    response *http.Response
    err      error
}

func main() {
    http.HandleFunc("/search", search)
    http.HandleFunc("/user-is-unique", unique)

    http.ListenAndServe(":8080", nil)
}

func asyncHttpGets(urls []string, qs string) []*HttpResponse {
    ch := make(chan *HttpResponse)
    responses := []*HttpResponse{}
    client := http.Client{}

    for _, url := range urls {
        go func(url string) {
            // fmt.Printf("Fetching %s \n", url + "?query=" + qs)
            resp, err := client.Get(url + "?query=" + qs)
            ch <- &HttpResponse{url, resp, err}
            if err != nil && resp != nil && resp.StatusCode == http.StatusOK {
                resp.Body.Close()
            }
        }(url)
    }

    for {
        select {
        case r := <-ch:
            // fmt.Printf("%s was fetched\n", r.url)
            if r.err != nil {
                // fmt.Println("with an error", r.err)
            }
            responses = append(responses, r)
            if len(responses) == len(urls) {
                return responses
            }
        // case <-time.After(50 * time.Millisecond):
        //     fmt.Printf(".")
        }
    }
    return responses
}

