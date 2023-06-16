package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pdmatrix/hackattic/internal/client"
	"github.com/pdmatrix/hackattic/pkg/serving_dns"
	// "github.com/pdmatrix/hackattic/pkg/websocket_chit_chat"
	//"github.com/pdmatrix/hackattic/pkg/reading_qr"
	//"github.com/pdmatrix/hackattic/pkg/brute_force_zip"
	//"github.com/pdmatrix/hackattic/pkg/dockerized_solutions"
	//"github.com/pdmatrix/hackattic/pkg/backup_restore"
	//"github.com/pdmatrix/hackattic/pkg/jotting_jwts"
	//"github.com/pdmatrix/hackattic/pkg/password_hashing"
	//"github.com/pdmatrix/hackattic/pkg/tales_of_ssl"
	//"github.com/pdmatrix/hackattic/pkg/collision_course"
	//"github.com/pdmatrix/hackattic/pkg/help_me_unpack"
	//"github.com/pdmatrix/hackattic/pkg/mini_miner"
)

func main() {
	c := client.NewHackatticClient(os.Getenv("HACKATTIC_ACCESS_TOKEN"))
	input, err := c.GetString("serving_dns")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Input: %s\n", input)
	output, err := serving_dns.Run(input)
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Output: %s\n", string(data))
	res, err := c.PostSolution("serving_dns", data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Result: %s", res)
}
