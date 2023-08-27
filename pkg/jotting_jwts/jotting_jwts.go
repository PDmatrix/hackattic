package jotting_jwts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/localtunnel/go-localtunnel"
)

type JottingJwts struct{}

type Data struct {
	JwtSecret string `json:"jwt_secret"`
}

type Output struct {
	AppUrl string `json:"app_url"`
}

func (d JottingJwts) Solve(input string) (interface{}, error) {
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

	type MyCustomClaims struct {
		Append string `json:"append"`
		jwt.StandardClaims
	}

	var solution = ""

	server := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			defer r.Body.Close()

			b, _ := io.ReadAll(r.Body)
			token, _ := jwt.ParseWithClaims(string(b), &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(data.JwtSecret), nil
			})

			if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
				if claims.Append != "" {
					solution = fmt.Sprintf("%s%s", solution, claims.Append)
				} else {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(fmt.Sprintf(`{"solution": "%s"}`, solution)))
				}
				fmt.Printf("Got append value: %s\n", claims.Append)
			} else {
				fmt.Println("Wrong token supplied")
			}

		}),
	}

	// Handle request from localtunnel
	go server.Serve(listener)

	output.AppUrl = listener.Addr().String()

	return output, nil
}
