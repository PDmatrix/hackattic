package visual_basic_math

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/pdmatrix/hackattic/internal/utils"
)

type VisualBasicMath struct{}

type Data struct {
	ImageUrl string `json:"image_url"`
}
type Output struct {
	Result string `json:"result"`
}

// TODO: may need to run multiple times because result are not so relaible, see run_untill_passed()
func (d VisualBasicMath) Solve(input string) (interface{}, error) {
	data := new(Data)
	output := new(Output)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}
	os.Mkdir("/tmp/tesseract", os.ModePerm)

	err = utils.DownloadFile(data.ImageUrl, "/tmp/tesseract/image.png")
	if err != nil {
		return nil, err
	}
	err = convertToGray()
	if err != nil {
		return nil, err
	}

	err = recognize()
	if err != nil {
		return nil, err
	}

	f, err := os.Open("/tmp/tesseract/recognition.txt")
	if err != nil {
		return nil, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	var result int64 = 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if len(line) == 8 && strings.HasPrefix(line, "1") {
			after, _ := strings.CutPrefix(line, "1")
			line = after
			line = "+" + line
		}
		if len(line) == 8 && strings.HasPrefix(line, "2") {
			after, _ := strings.CutPrefix(line, "2")
			line = after
			line = "÷" + line
		}
		lineWithoutPrefix := strings.Replace(strings.Replace(strings.Replace(strings.Replace(line, "÷", "", 1), "x", "", 1), "-", "", 1), "+", "", 1)
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "+") || strings.HasPrefix(line, "x") || strings.HasPrefix(line, "÷") {
		} else {
			return nil, fmt.Errorf("bad string %s", line)
		}
		num_parsed, err := strconv.Atoi(lineWithoutPrefix)
		if err != nil {
			return nil, err
		}
		num := int64(num_parsed)

		if len(line) == 0 {
			return nil, fmt.Errorf("len(line) is equal 0")
		}
		fmt.Printf("Original: %s. Parsed: %d\n", line, num)
		switch line[0] {
		case '+':
			result += num
		case 'x':
			result *= num
		case '-':
			result -= num
		default:
			result = int64(math.Floor(float64(result) / float64(num)))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	output.Result = fmt.Sprintf("%d", result)

	return output, nil
}

func recognize() error {
	ctx := context.Background()
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	exPath := filepath.Dir(ex)
	path := filepath.Join(exPath, "..", "..", "pkg", "visual_basic_math", "eng_nums5.traineddata")
	trainedDataFile, err := os.Open(path)
	if err != nil {
		return err
	}
	tempData, err := os.Create("/tmp/tesseract/eng_nums5.traineddata")
	if err != nil {
		return err
	}

	_, err = io.Copy(tempData, trainedDataFile)
	if err != nil {
		return err
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	reader, err := cli.ImagePull(ctx, "clearlinux/tesseract-ocr", types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      "clearlinux/tesseract-ocr",
		Cmd:        []string{"tesseract", "-l", "eng_nums5", "--psm", "6", "-c", "tessedit_char_whitelist=-+0123456789x÷", "image_gray.png", "recognition"},
		Tty:        false,
		WorkingDir: "/app",
		Env:        []string{"TESSDATA_PREFIX=/app"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   "bind",
				Source: "/tmp/tesseract",
				Target: "/app",
			},
		},
	}, nil, nil, "")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})

	return nil
}

func convertToGray() error {
	file, err := os.Open("/tmp/tesseract/image.png")
	if err != nil {
		return err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	// grayImg := adjust.Gamma(img, 0)
	//Converting image to grayscale
	grayImg := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			R, G, B, _ := img.At(x, y).RGBA()

			if R == 65535 && G == 65535 && B == 65535 {
				grayImg.Set(x, y, color.White)
			} else {
				grayImg.Set(x, y, color.Black)
			}
		}
	}

	f, err := os.Create("/tmp/tesseract/image_gray.png")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := png.Encode(f, grayImg); err != nil {
		return err
	}

	return nil
}
