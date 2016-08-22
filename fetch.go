package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

var urls map[string]string

type HttpResponse struct {
    url      string
    response *http.Response
    err      error
}

func main() {
    urls = map[string]string{
        "app1": "http://app1.tillersystems.com/api",
    }

    http.HandleFunc("/urls", func(w http.ResponseWriter, r *http.Request) {
        defer r.Body.Close() 

        if r.Method == "POST" {
            r.ParseForm()
            subDomain := r.FormValue("sub-domain")
            if "" != subDomain {
                urls[subDomain] = fmt.Sprintf("http://%s.tillersystems.com/api", subDomain)
            }
        }
        
        json, _ := json.Marshal(urls)
        fmt.Fprint(w, string(json))
    })
    
    http.HandleFunc("/search", search)
    http.HandleFunc("/user-is-unique", unique)

    http.ListenAndServe(":8080", nil)
}

func asyncHttpGets(urls []string) []*HttpResponse {
    ch := make(chan *HttpResponse)
    responses := []*HttpResponse{}
    client := http.Client{}

    for _, url := range urls {
        go func(url string) {
            // fmt.Printf("Fetching %s \n", url + "?query=" + qs)
            resp, err := client.Get(url)
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

