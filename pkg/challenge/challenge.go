package challenge

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pdmatrix/hackattic/internal/client"
	"github.com/pdmatrix/hackattic/pkg/backup_restore"
	"github.com/pdmatrix/hackattic/pkg/basic_face_detection"
	"github.com/pdmatrix/hackattic/pkg/brute_force_zip"
	"github.com/pdmatrix/hackattic/pkg/collision_course"
	"github.com/pdmatrix/hackattic/pkg/dockerized_solutions"
	"github.com/pdmatrix/hackattic/pkg/help_me_unpack"
	"github.com/pdmatrix/hackattic/pkg/jotting_jwts"
	"github.com/pdmatrix/hackattic/pkg/mini_miner"
	"github.com/pdmatrix/hackattic/pkg/password_hashing"
	"github.com/pdmatrix/hackattic/pkg/reading_qr"
	"github.com/pdmatrix/hackattic/pkg/serving_dns"
	"github.com/pdmatrix/hackattic/pkg/tales_of_ssl"
	"github.com/pdmatrix/hackattic/pkg/the_redis_one"
	"github.com/pdmatrix/hackattic/pkg/touch_tone_dialing"
	"github.com/pdmatrix/hackattic/pkg/visual_basic_math"
	"github.com/pdmatrix/hackattic/pkg/websocket_chit_chat"
)

type Challenge interface {
	Solve(input string) (interface{}, error)
}

var challenges map[string]Challenge = map[string]Challenge{
	"backup_restore":       new(backup_restore.BackupRestore),
	"basic_face_detection": new(basic_face_detection.BasicFaceDetection),
	"brute_force_zip":      new(brute_force_zip.BruteForceZip),
	"collision_course":     new(collision_course.CollisionCourse),
	"dockerized_solutions": new(dockerized_solutions.DockerizedSolutions),
	"help_me_unpack":       new(help_me_unpack.HelpMeUnpack),
	"jotting_jwts":         new(jotting_jwts.JottingJwts),
	"mini_miner":           new(mini_miner.MiniMiner),
	"password_hashing":     new(password_hashing.PasswordHashing),
	"reading_qr":           new(reading_qr.ReadingQr),
	"serving_dns":          new(serving_dns.ServingDns),
	"tales_of_ssl":         new(tales_of_ssl.TalesOfSsl),
	"the_redis_one":        new(the_redis_one.TheRedisOne),
	"touch_tone_dialing":   new(touch_tone_dialing.TouchToneDialing),
	"visual_basic_math":    new(visual_basic_math.VisualBasicMath),
	"websocket_chit_chat":  new(websocket_chit_chat.WebsocketChitChat),
}

func GetSolution(challenge string, playground bool) (string, error) {
	c := client.NewHackatticClient(os.Getenv("HACKATTIC_ACCESS_TOKEN"))
	input, err := c.GetString(challenge)
	if err != nil {
		return "", err
	}
	fmt.Printf("Input: %s\n", input)
	output, err := challenges[challenge].Solve(input)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	fmt.Printf("Output: %s\n", string(data))

	additionalParams := ""
	if playground {
		additionalParams = "&playground=1"
	}

	res, err := c.PostSolution(challenge, data, additionalParams)
	if err != nil {
		return "", err
	}

	fmt.Printf("Result: %s\n", res)

	return res, nil
}
