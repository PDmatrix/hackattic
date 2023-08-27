package mini_miner

import (
	"crypto/sha256"
	"encoding/json"
)

type MiniMiner struct{}

type Data struct {
	Difficulty int `json:"difficulty"`
	Block      struct {
		Data  [][]interface{} `json:"data"`
		Nonce int             `json:"nonce"`
	} `json:"block"`
}

type Output struct {
	Nonce int `json:"nonce"`
}

func (d MiniMiner) Solve(input string) (interface{}, error) {
	data := new(Data)
	output := new(Output)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}
	byteArr := []uint8{1, 2, 4, 8, 16, 32, 64, 128}
	data.Block.Nonce = 0
	for {
		data.Block.Nonce = data.Block.Nonce + 1
		out, err := json.Marshal(data.Block)
		if err != nil {
			return nil, err
		}

		h := sha256.Sum256([]byte(out))
		if h[0] == 0 && h[1] < byteArr[16-data.Difficulty] {
			break
		}
	}

	output.Nonce = data.Block.Nonce

	return output, nil
}
