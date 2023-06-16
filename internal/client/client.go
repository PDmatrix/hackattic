package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type HackatticClient struct {
	AccessToken string
	httpClient  *http.Client
}

func NewHackatticClient(accessToken string) *HackatticClient {
	return &HackatticClient{
		AccessToken: accessToken,
		httpClient:  &http.Client{},
	}
}

func (c *HackatticClient) GetString(challenge string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://hackattic.com/challenges/%s/problem?access_token=%s", challenge, c.AccessToken), nil)
	if err != nil {
		return "", err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *HackatticClient) PostSolution(challenge string, data []byte, additionalParams ...string) (string, error) {
	params := ""
	if len(additionalParams) > 0 {
		params = additionalParams[0]
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://hackattic.com/challenges/%s/solve?access_token=%s%s", challenge, c.AccessToken, params), bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
