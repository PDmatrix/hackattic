package serving_dns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/localtunnel/go-localtunnel"
	"github.com/miekg/dns"
)

type Data struct {
	Records []struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Data string `json:"data"`
	} `json:"records"`
}
type Output struct {
	DnsIp   string `json:"dns_ip"`
	DnsPort string `json:"dns_port"`
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
func Run(input string) (*Output, error) {
	data := new(Data)
	output := new(Output)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	listener, err := localtunnel.Listen(localtunnel.Options{})
	if err != nil {
		panic(err)
	}
	server := &dns.Server{Addr: ":8000", Net: "tcp", Listener: listener}
	go server.ListenAndServe()
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		fmt.Printf("GotRequest")
	})

	addr := strings.Replace(listener.Addr().String(), "https://", "", -1)
	fmt.Println(addr)

	return output, nil
}
