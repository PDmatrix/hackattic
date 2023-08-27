package backup_restore

import (
	"bytes"
	"compress/gzip"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v4"
)

type BackupRestore struct{}

type Data struct {
	Dump string `json:"Dump"`
}

type Output struct {
	AliveSsns []string `json:"alive_ssns"`
}

func (d BackupRestore) Solve(input string) (interface{}, error) {
	data := new(Data)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	dump, err := b64.StdEncoding.DecodeString(data.Dump)
	if err != nil {
		return nil, err
	}

	unziped, err := gUnzipData(dump)
	if err != nil {
		return nil, err
	}

	f, err := os.Create("/tmp/dump.sql")
	if err != nil {
		return nil, err
	}

	defer f.Close()

	f.Write(unziped)

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, "postgres:10", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "postgres:10",
		Tty:   false,
		Env:   []string{"POSTGRES_PASSWORD=password"},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"5432/tcp": []nat.PortBinding{
				{
					HostIP:   "",
					HostPort: "5432",
				},
			},
		},
	}, nil, nil, "")
	if err != nil {
		return nil, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	defer cli.ContainerStop(ctx, resp.ID, nil)

	time.Sleep(time.Second * 2)

	var out bytes.Buffer
	cmd := exec.Command("/usr/bin/psql", "-p", "5432", "-h", "localhost", "-U", "postgres", "-w", "-f", "/tmp/dump.sql")
	cmd.Env = []string{"PGPASSWORD=password"}
	cmd.Stderr = &out
	cmd.Stdout = &out
	err = cmd.Run()
	fmt.Printf("%s\n", out.String())
	if err != nil {
		fmt.Printf("%+v", err)
		return nil, nil
	}

	fmt.Print("Dumped \n")

	urlExample := "postgres://postgres:password@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "select ssn from criminal_records where status='alive'")
	if err != nil {
		return nil, err
	}

	output := new(Output)
	for rows.Next() {
		var n string
		rows.Scan(&n)
		output.AliveSsns = append(output.AliveSsns, n)
	}

	return output, nil
}

func gUnzipData(data []byte) (resData []byte, err error) {
	b := bytes.NewBuffer(data)

	var r io.Reader
	r, err = gzip.NewReader(b)
	if err != nil {
		return
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return
	}

	resData = resB.Bytes()

	return
}
