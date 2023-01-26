package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PageVersion struct {
	Message string `json:"message"`
	Number  int16  `json:"number"`
}

type PageStorage struct {
	Representation string `json:"representation"`
	Value          string `json:"value"`
}

type PageBody struct {
	Storage PageStorage `json:"storage"`
}

type PageSpace struct {
	Key string `json:"key"`
}

type ConfluencePage struct {
	Body     PageBody    `json:"body"`
	PageType string      `json:"type"`
	Space    PageSpace   `json:"space"`
	Status   string      `json:"status"`
	Title    string      `json:"title"`
	Version  PageVersion `json:"version"`
}

func getConfluencePage(pageId int32) (*ConfluencePage, error) {
	client := &http.Client{}
	var currentPage ConfluencePage

	url := fmt.Sprintf("https://%s.atlassian.net/wiki/rest/api/content/%v?expand=version,body.storage", CONFLUENCE_CLOUD_DOMAIN, pageId)
	fmt.Printf("Building GET confluence page request %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &currentPage, err
	}
	req.SetBasicAuth(CONFLUENCE_USER, CONFLUENCE_TOKEN)

	fmt.Println("Sending GET confluence page request")
	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		return &currentPage, fmt.Errorf("GET confluence page Failed %v", resp.StatusCode)
	}
	if err != nil {
		return &currentPage, err
	}
	defer resp.Body.Close()

	fmt.Println("Reading GET confluence page request")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &currentPage, err
	}

	fmt.Println("Casting GET confluence page request")
	json.Unmarshal(body, &currentPage)
	// fmt.Println(currentPage.Body.Storage.Value)

	return &currentPage, nil
}

func putConfluencePage(pageId int32, content *ConfluencePage) error {
	client := &http.Client{}

	// fmt.Println(content.Body.Storage.Value)
	fmt.Println("Casting to JSON PUT confluence page")
	var payload bytes.Buffer

	enc := json.NewEncoder(&payload)
	enc.SetEscapeHTML(false)
	enc.Encode(&content)
	fmt.Println(&payload)

	url := fmt.Sprintf("https://%s.atlassian.net/wiki/rest/api/content/%v", CONFLUENCE_CLOUD_DOMAIN, pageId)

	fmt.Println("Building PUT confluence page request")
	req, err := http.NewRequest("PUT", url, &payload)
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth(CONFLUENCE_USER, CONFLUENCE_TOKEN)
	if err != nil {
		return err
	}

	fmt.Println("Sending PUT confluence page request")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("PUT confluence page Failed (%v)\n", resp.StatusCode)
	}

	defer resp.Body.Close()

	return nil
}
