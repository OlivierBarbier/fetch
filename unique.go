package main

/*
 * This file handle the business logic behind
 * /api/search
 */

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"
)

type uniqueResponse struct {
    Response  bool
    Url string
}

func unique (w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close() 
    w.Header().Set("Access-Control-Allow-Origin", "*")
    
    // fmt.Println(r.URL.String())
    // fmt.Println(r.URL.Path)

    r.ParseForm()

    // resto := r.Form.Get("u")

    combined := uniqueResponse{false, "*"}

    var u []string;

	for _, url := range urls {
	    u = append(u, url + r.URL.String());
	}
	// fmt.Println(urls)
    
    results := asyncHttpGets(u)

    for _, result := range results {
        if result == nil || result.response == nil {
            continue
        }

        if result.response.Status != "200 OK" {
            continue
        }

        // fmt.Printf("%s status: %s\n", result.url,
        //        result.response.Status)

        var response uniqueResponse;

        parseHttpUniqueResponse(result, &response);

        response.Url = result.url
// fmt.Println(response.Response)
        combined.Response = combined.Response || response.Response;
    }

    json, _ := json.Marshal(combined)

    fmt.Fprint(w, string(json))
}

func parseHttpUniqueResponse(result *HttpResponse, response *uniqueResponse) {
    raw_response, _ := ioutil.ReadAll(result.response.Body)

    json.Unmarshal(raw_response, response)
}
