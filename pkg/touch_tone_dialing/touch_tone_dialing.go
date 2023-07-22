package touch_tone_dialing

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	dtmf "github.com/Hallicopter/go-dtmf/dtmf"
)

type Data struct {
	WavUrl string `json:"wav_url"`
}
type Output struct {
	Sequence string `json:"sequence"`
}

// TODO: may need to run multiple times because result are not so relaible
func Run(input string) (*Output, error) {
	data := new(Data)
	output := new(Output)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	err = downloadFile(data.WavUrl, "/tmp/file.wav")
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

func downloadFile(url string, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
