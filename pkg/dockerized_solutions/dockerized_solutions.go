package dockerized_solutions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/localtunnel/go-localtunnel"
)

type DockerizedSolutions struct{}

type Data struct {
	Credentials struct {
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"credentials"`
	IgnitionKey  string `json:"ignition_key"`
	TriggerToken string `json:"trigger_token"`
}

type Output struct {
	AppUrl string `json:"app_url"`
}

var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

// 413 with localtunnel :(
func (d DockerizedSolutions) Solve(input string) (interface{}, error) {
	data := new(Data)
	output := new(Output)
	var wg sync.WaitGroup
	wg.Add(1)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	listener, err := localtunnel.Listen(localtunnel.Options{})
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Handler: http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
			log.Println(req.Method, " ", req.URL, " ")

			dur, _ := time.ParseDuration("5m")
			client := &http.Client{
				Timeout: dur,
			}

			req.Header.Set("Host", req.Host)
			req.Header.Set("X-Forwarded-For", req.Header.Get("X-Forwarded-For"))
			//http: Request.RequestURI can't be set in client requests.
			//http://golang.org/src/pkg/net/http/client.go
			req.RequestURI = ""
			req.URL.Host = "localhost:5000"
			req.URL.Scheme = "http"
			req.Header.Set("X-Forwarded-Proto", "https")
			req.Header.Set("Docker-Distribution-Api-Version", "registry/2.0")

			delHopHeaders(req.Header)

			resp, err := client.Do(req)
			if err != nil {
				http.Error(wr, "Server Error", http.StatusInternalServerError)
				log.Fatal("ServeHTTP:", err)
			}
			defer resp.Body.Close()

			log.Println(req.RemoteAddr, " ", resp.Status)

			delHopHeaders(resp.Header)

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			fmt.Printf("%v", bodyString)

			copyHeader(wr.Header(), resp.Header)
			wr.WriteHeader(resp.StatusCode)
			io.Copy(wr, resp.Body)
		}),
	}

	// Handle request from localtunnel
	go server.Serve(listener)

	addr := strings.Replace(listener.Addr().String(), "https://", "", -1)

	fmt.Println(addr)
	req, err := http.NewRequest("POST", fmt.Sprintf("https://hackattic.com/_/push/%s", data.TriggerToken), bytes.NewBuffer([]byte(fmt.Sprintf(`{"registry_host":"%s"}`, addr))))
	if err != nil {
		return nil, err
	}
	dur, _ := time.ParseDuration("5m")
	client := &http.Client{
		Timeout: dur,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// bodyBytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// bodyString := string(bodyBytes)
	// fmt.Printf("%v", bodyString)
	output.AppUrl = listener.Addr().String()
	return output, nil
}
