package touch_tone_dialing

import (
	"encoding/json"

	dtmf "github.com/Hallicopter/go-dtmf/dtmf"
	"github.com/pdmatrix/hackattic/internal/utils"
)

type TouchToneDialing struct{}

type Data struct {
	WavUrl string `json:"wav_url"`
}
type Output struct {
	Sequence string `json:"sequence"`
}

// TODO: may need to run multiple times because result are not so relaible
func (d TouchToneDialing) Solve(input string) (interface{}, error) {
	data := new(Data)
	output := new(Output)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	err = utils.DownloadFile(data.WavUrl, "/tmp/file.wav")
	if err != nil {
		return nil, err
	}

	touch, err := dtmf.DecodeDTMFFromFile("/tmp/file.wav", 4000, 5)
	if err != nil {
		return nil, err
	}
	output.Sequence = touch
	return output, nil
}
