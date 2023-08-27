package a_global_presence

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type AGlobalPresence struct{}

type Data struct {
	PresenceToken string `json:"presence_token"`
}

type Output struct{}

// TODO: may need to run multiple times because result are not so relaible
func (d AGlobalPresence) Solve(input string) (interface{}, error) {
	data := new(Data)
	output := new(Output)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)

	path := filepath.Join(exPath, "..", "..", "pkg", "a_global_presence", "proxy.txt")
	proxyList, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer proxyList.Close()

	var arr []string

	scanner := bufio.NewScanner(proxyList)
	for scanner.Scan() {
		arr = append(arr, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	ch := make(chan string)
	for _, p := range arr {
		go func(p string) {
			err := sendProxyRequest(p, data.PresenceToken)
			if err != nil {
				fmt.Printf("Got error: %v\n", err)
			}
			if err == nil {
				ch <- p
			}
		}(p)
	}

	time.Sleep(time.Second * 10)

	return output, nil
}

func sendProxyRequest(proxy string, presenceToken string) error {
	fmt.Printf("Using proxy: %s\n", proxy)
	proxyURL, err := url.Parse("http://" + proxy)

	if err != nil {
		log.Println(err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 5,
	}

	res, err := client.Get(fmt.Sprintf("https://hackattic.com/_/presence/%s", presenceToken))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Got response: %s\n", string(b))

	return nil
}
