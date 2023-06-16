package reading_qr

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"

	"github.com/liyue201/goqr"
)

type Data struct {
	ImageUrl string `json:"image_url"`
}

type Output struct {
	Code string `json:"code"`
}

// fcrackzip need to install
func Run(input string) (*Output, error) {
	data := new(Data)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(data.ImageUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create("/tmp/image.png")
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, err
	}

	file, err := os.Open("/tmp/image.png")
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// prepare BinaryBitmap
	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		fmt.Printf("Recognize failed: %v\n", err)
		return nil, err
	}

	output := new(Output)
	output.Code = string(qrCodes[0].Payload)
	return output, nil
}
