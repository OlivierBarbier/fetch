package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  // "time"
  "io/ioutil"
)

var urls = []string{
  "http://www.mocky.io/v2/57b722090f0000b70e0b7d0e",
  "http://www.mocky.io/v2/57b722470f0000990e0b7d0f",
}

type HttpResponse struct {
  url      string
  response *http.Response
  err      error
}

type Response struct {
	Restaurant  string
}

func asyncHttpGets(urls []string) []*HttpResponse {
  ch := make(chan *HttpResponse)
  responses := []*HttpResponse{}
  client := http.Client{}
  for _, url := range urls {
      go func(url string) {
          // fmt.Printf("Fetching %s \n", url)
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

func parseHttpResponse(result *HttpResponse, response *[]Response) {
  raw_response, _ := ioutil.ReadAll(result.response.Body)
  
  json.Unmarshal(raw_response, response)
}

func main() {

  http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
	  combined := []Response{}

	  results := asyncHttpGets(urls)

	  for _, result := range results {
	      if result != nil && result.response != nil {
	          
	          // fmt.Printf("%s status: %s\n", result.url,
	          //        result.response.Status)

	          if result.response.Status != "200 OK" {
			  	continue
			  }
	          
			  var response []Response;

			  parseHttpResponse(result, &response);

		      for _, res := range response {
		    	combined = append(combined, res);
			  }

	      }
	  }
	  
	  j, _ := json.Marshal(combined)
	fmt.Fprint(w, string(j))
  })

  http.ListenAndServe(":8080", nil)
}
