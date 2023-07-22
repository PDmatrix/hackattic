package visual_basic_math

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type Data struct {
	ImageUrl string `json:"image_url"`
}
type Output struct {
	Result string `json:"result"`
}

// TODO: may need to run multiple times because result are not so relaible
func Run(input string) (*Output, error) {
	data := new(Data)
	output := new(Output)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	err = downloadFile(data.ImageUrl, "/tmp/tesseract/image.png")
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
	result := 0
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, "x", "*", 1)
		line = strings.Replace(line, " ", "", -1)
		line = strings.Replace(line, "%", "7", -1)
		line = strings.Replace(line, "F", "7", -1)
		line = strings.Replace(line, "Z", "7", -1)
		line = strings.Replace(line, "®", "0", -1)
		line = strings.Replace(line, "q", "9", -1)
		line = strings.Replace(line, "¢", "", -1)
		lineWithoutPrefix := line[1:]
		num, err := strconv.Atoi(lineWithoutPrefix)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Original: %s. Parsed: %d\n", line, num)
		switch line[0] {
		case '+':
			result += num
		case '*':
			result *= num
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return output, nil
}

func recognize() error {
	ctx := context.Background()
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
		Cmd:        []string{"tesseract", "image_gray.png", "recognition"},
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
	grayImg := adjust.Gamma(img, 0)

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
