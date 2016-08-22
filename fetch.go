package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
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

    router := mux.NewRouter().StrictSlash(true)

    router.HandleFunc("/urls", func(w http.ResponseWriter, r *http.Request) {
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
    
    router.HandleFunc("/urls/{id}", func(w http.ResponseWriter, r *http.Request) {
        defer r.Body.Close() 
        vars := mux.Vars(r) 
        id := vars["id"]

        if r.Method == "DELETE" {
            delete(urls, id)
        }

        if r.Method == "GET" {
            json, _ := json.Marshal(urls[id])
            fmt.Fprint(w, string(json))
        }        
    })

    router.HandleFunc("/search", search)
    router.HandleFunc("/user-is-unique", unique)

    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { 
        defer r.Body.Close()
        w.Header().Set("Access-Control-Allow-Origin", "*")
    })

    http.ListenAndServe(":8080", router)
}

func asyncHttpGets(urls []string) []*HttpResponse {
    ch := make(chan *HttpResponse)
    responses := []*HttpResponse{}
    client := http.Client{}

    for _, url := range urls {
        go func(url string) {
            fmt.Printf("Fetching %s \n", url)
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

