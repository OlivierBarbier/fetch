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

type searchResponse struct {
    Name  string
    Url string
    Id int
}

func search (w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    w.Header().Set("Access-Control-Allow-Origin", "*") 
    // fmt.Println(r.URL.String())
    // fmt.Println(r.URL.Path)

    r.ParseForm()

    // resto := r.Form.Get("r")

    combined := []searchResponse{}

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

        var response []searchResponse;

        parseHttpResponse(result, &response);

        for _, res := range response {
        	res.Url = result.url
            combined = append(combined, res);
        }
    }

    json, _ := json.Marshal(combined)

    fmt.Fprint(w, string(json))
}

func parseHttpResponse(result *HttpResponse, response *[]searchResponse) {
    raw_response, _ := ioutil.ReadAll(result.response.Body)

    json.Unmarshal(raw_response, response)
}
