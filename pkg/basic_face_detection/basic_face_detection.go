package basic_face_detection

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	pigo "github.com/esimov/pigo/core"
)

type Data struct {
	ImageUrl string `json:"image_url"`
}
type Output struct {
	FaceTiles [][]int `json:"face_tiles"`
}

// TODO: may need to run multiple times because result are not so relaible
func Run(input string) (*Output, error) {
	data := new(Data)
	output := new(Output)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	err = downloadFile(data.ImageUrl, "/tmp/image.png")
	if err != nil {
		return nil, err
	}

	cascadeFileUrl := "https://github.com/esimov/pigo/raw/e922e5442d3895b64cbb070a720c72a1a2d2c2da/cascade/facefinder"

	err = downloadFile(cascadeFileUrl, "/tmp/facefinder")
	if err != nil {
		return nil, err
	}

	cascadeFile, err := ioutil.ReadFile("/tmp/facefinder")
	if err != nil {
		log.Fatalf("Error reading the cascade file: %v", err)
	}

	src, err := pigo.GetImage("/tmp/image.png")
	if err != nil {
		return nil, err
	}

	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	cParams := pigo.CascadeParams{
		MinSize:     20,
		MaxSize:     150,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,

		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}

	pigo := pigo.NewPigo()
	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	classifier, err := pigo.Unpack(cascadeFile)
	if err != nil {
		return nil, err
	}

	angle := 0.0

	dets := classifier.RunCascade(cParams, angle)

	// Calculate the intersection over union (IoU) of two clusters.
	dets = classifier.ClusterDetections(dets, 0.2)

	arr := make([][]int, len(dets))
	for i, det := range dets {
		arr[i] = []int{(det.Row / 100), (det.Col / 100)}
	}
	output.FaceTiles = arr
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
