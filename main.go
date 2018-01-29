package main

import (
	"os"
	"fmt"
	"net/url"
	"net/http"
	"time"
	"io/ioutil"
)

func main() {
	usage := `
highlights DESCRIPTION [SOURCE]

DESCRIPTION is the description which will be used for this highlight
SOURCE is an optional source under which this highlight will be grouped
`

	apiKey, ok := os.LookupEnv("RESCUETIME_API_KEY")
	if !ok {
		fmt.Println("RESCUETIME_API_KEY variable is missing")
		os.Exit(1)
	}

	if len(os.Args) < 2 || len(os.Args) > 3 {
		handleError(usage)
	}

	var src, desc string
	desc = os.Args[1]
	if len(os.Args) == 3 {
		src = os.Args[2]
	}

	rescueUrl := url.URL{
		Host:   "www.rescuetime.com",
		Scheme: "https",
		Path:   "anapi/highlights_post",
	}


	query := url.Values{}
	query.Add("key", apiKey)
	query.Add("description", desc)
	query.Add("highlight_date", time.Now().Format("2006-01-02"))
	if src != "" {
		query.Add("source", src)
	}

	req, err := http.NewRequest(http.MethodPost, rescueUrl.String(), nil)
	handleError(err)
	req.URL.RawQuery = query.Encode()

	client := http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	handleError(err)

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		handleError(err)
		handleError(string(body))
	}
}

func handleError(err interface{}) {
	if err == nil {
		return
	}

	switch err.(type) {
	case error:
		fmt.Println(err)
		os.Exit(1)
	case string:
		fmt.Println(err)
		os.Exit(1)
	}
}
