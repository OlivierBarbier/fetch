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
    Status  boolean
    Url string
}

func unique (w http.ResponseWriter, r *http.Request) {
    // fmt.Println(r.URL.String())
    // fmt.Println(r.URL.Path)

    r.ParseForm()

    resto := r.Form.Get("r")

    combined := []uniqueResponse{}

    var u []string;

	for _, url := range urls {
	    u = append(u, url + r.URL.Path);
	}
	// fmt.Println(urls)
    
    results := asyncHttpGets(u, resto)

    for _, result := range results {
        if result == nil || result.response == nil {
            continue
        }

        if result.response.Status != "200 OK" {
            continue
        }

        // fmt.Printf("%s status: %s\n", result.url,
        //        result.response.Status)

        var response []uniqueResponse;

        parseHttpResponse(result, &response);

        for _, res := range response {
        	res.Url = result.url
            combined = append(combined, res);
        }
    }

    json, _ := json.Marshal(combined)

    fmt.Fprint(w, string(json))
}

func parseHttpResponse(result *HttpResponse, response *[]uniqueResponse) {
    raw_response, _ := ioutil.ReadAll(result.response.Body)

    json.Unmarshal(raw_response, response)
}
