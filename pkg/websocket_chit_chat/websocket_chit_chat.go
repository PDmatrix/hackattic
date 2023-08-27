package websocket_chit_chat

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketChitChat struct{}

type Data struct {
	Token string `json:"token"`
}

type Output struct {
	Secret string `json:"secret"`
}

func (d WebsocketChitChat) Solve(input string) (interface{}, error) {
	data := new(Data)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	u := url.URL{Scheme: "wss", Host: "hackattic.com", Path: fmt.Sprintf("/_/ws/%s", data.Token)}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	pastTime := time.Now()

	output := new(Output)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Printf("Received: %s", message)
		decodedMessage := string(message)
		if decodedMessage == "ping!" {
			// log.Printf("%v00", time.Since(pastTime).Milliseconds()/100)
			mes := fmt.Sprintf("%v", getTime(time.Since(pastTime).Milliseconds()))
			log.Printf("%s", mes)
			err = c.WriteMessage(websocket.TextMessage, []byte(mes))
			if err != nil {
				log.Println("Error during message writing:", err)
				break
			}
			pastTime = time.Now()
		}

		if strings.Contains(decodedMessage, "congratulations") {
			re, _ := regexp.Compile("congratulations! the solution to this challenge is \"(?P<answer>.+)\"")
			matches := re.FindStringSubmatch(decodedMessage)
			pwdIndex := re.SubexpIndex("answer")
			output.Secret = matches[pwdIndex]
			break
		}
	}

	return output, nil
}

func getTime(ms int64) int {
	//  700, 1500, 2000, 2500 or 3000
	switch {
	case ms < 900:
		return 700
	case ms >= 900 && ms < 1750:
		return 1500
	case ms >= 1750 && ms < 2250:
		return 2000
	case ms >= 2250 && ms < 2750:
		return 2500
	default:
		return 3000
	}
}
