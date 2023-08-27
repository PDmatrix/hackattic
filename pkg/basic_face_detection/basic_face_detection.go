package basic_face_detection

import (
	"encoding/json"
	"log"
	"os"

	pigo "github.com/esimov/pigo/core"
	"github.com/pdmatrix/hackattic/internal/utils"
)

type BasicFaceDetection struct{}

type Data struct {
	ImageUrl string `json:"image_url"`
}
type Output struct {
	FaceTiles [][]int `json:"face_tiles"`
}

func (d BasicFaceDetection) Solve(input string) (interface{}, error) {
	data := new(Data)
	output := new(Output)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	err = utils.DownloadFile(data.ImageUrl, "/tmp/image.png")
	if err != nil {
		return nil, err
	}

	cascadeFileUrl := "https://github.com/esimov/pigo/raw/e922e5442d3895b64cbb070a720c72a1a2d2c2da/cascade/facefinder"

	err = utils.DownloadFile(cascadeFileUrl, "/tmp/facefinder")
	if err != nil {
		return nil, err
	}

	cascadeFile, err := os.ReadFile("/tmp/facefinder")
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
