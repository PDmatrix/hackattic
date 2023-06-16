package collision_course

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type Data struct {
	Include string `json:"include"`
}

type Output struct {
	Files []string `json:"files"`
}

func Run(input string) (*Output, error) {
	data := new(Data)
	output := new(Output)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	f, err := os.Create("/tmp/input")

	if err != nil {
		return nil, err
	}

	defer f.Close()

	f.WriteString(data.Include)

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, "brimstone/fastcoll", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      "brimstone/fastcoll",
		Cmd:        []string{"--prefixfile", "input", "-o", "msg1.bin", "msg2.bin"},
		Tty:        false,
		WorkingDir: "/work",
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   "bind",
				Source: "/tmp",
				Target: "/work",
			},
		},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	b, err := os.ReadFile("/tmp/msg1.bin")
	if err != nil {
		return nil, err
	}

	var arr []string
	arr = append(arr, b64.StdEncoding.EncodeToString(b))

	b, err = os.ReadFile("/tmp/msg2.bin")
	if err != nil {
		return nil, err
	}
	arr = append(arr, b64.StdEncoding.EncodeToString(b))

	output.Files = arr

	fmt.Printf("%+v\n", output)

	return output, nil
}
