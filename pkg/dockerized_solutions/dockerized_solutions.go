package dockerized_solutions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
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
	Secret string `json:"secret"`
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

func (d DockerizedSolutions) Solve(input string) (interface{}, error) {
	data := new(Data)
	output := new(Output)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	// To use this, login to ngrok and export NGROK_AUTHTOKEN variable
	listener, err := ngrok.Listen(context.TODO(),
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return nil, err
	}

	server := http.Server{
		Handler: http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
			log.Println(req.Method, " ", req.URL, " ")

			dur, _ := time.ParseDuration("5m")
			proxyClient := &http.Client{
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

			resp, err := proxyClient.Do(req)
			if err != nil {
				http.Error(wr, "Server Error", http.StatusInternalServerError)
				log.Fatal("ServeHTTP:", err)
			}
			defer resp.Body.Close()

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

	// Handle request from ngrok
	go server.Serve(listener)

	addr := strings.Replace(listener.Addr().String(), "https://", "", -1)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://hackattic.com/_/push/%s", data.TriggerToken), bytes.NewBuffer([]byte(fmt.Sprintf(`{"registry_host":"%s"}`, addr))))
	if err != nil {
		return nil, err
	}
	dur, _ := time.ParseDuration("5m")
	httpClient := &http.Client{
		Timeout: dur,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	bodyString := string(bodyBytes)
	re, _ := regexp.Compile(`\\"Tag\\":\\".+?\\"`)
	matches := re.FindStringSubmatch(bodyString)
	answer := ""
	for _, match := range matches {
		re2, _ := regexp.Compile(`\d{0,3}\.\d{0,3}\.\d{0,3}`)
		tag := re2.FindString(match)

		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		reader, err := cli.ImagePull(ctx, fmt.Sprintf("localhost:5000/hack:%s", tag), types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
		io.Copy(os.Stdout, reader)

		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: fmt.Sprintf("localhost:5000/hack:%s", tag),
			Tty:   false,
			Env:   []string{fmt.Sprintf("IGNITION_KEY=%s", data.IgnitionKey)},
		}, nil, nil, nil, "")
		if err != nil {
			return nil, err
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return nil, err
		}

		time.Sleep(time.Second * 2)

		reader, err = cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
			ShowStdout: true,
		})
		if err != nil {
			return nil, err
		}
		p := make([]byte, 8)
		reader.Read(p)
		content, _ := io.ReadAll(reader)
		l := strings.Trim(string(content), "\n")

		reader.Close()
		if l != "oops, wrong image!" {
			answer = l
		}

		if err := cli.ContainerStop(ctx, resp.ID, nil); err != nil {
			return nil, err
		}
	}

	output.Secret = answer

	return output, nil
}
