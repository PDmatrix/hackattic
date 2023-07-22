package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pdmatrix/hackattic/internal/client"
	"github.com/pdmatrix/hackattic/pkg/visual_basic_math"
	//"github.com/pdmatrix/hackattic/pkg/touch_tone_dialing"
	//"github.com/pdmatrix/hackattic/pkg/basic_face_detection"
	//"github.com/pdmatrix/hackattic/pkg/serving_dns"
	//"github.com/pdmatrix/hackattic/pkg/websocket_chit_chat"
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
	input, err := c.GetString("visual_basic_math")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Input: %s\n", input)
	output, err := visual_basic_math.Run(input)
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Output: %s\n", string(data))
	res, err := c.PostSolution("visual_basic_math", data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Result: %s", res)
}
