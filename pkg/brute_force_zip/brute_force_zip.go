package brute_force_zip

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type BruteForceZip struct{}

type Data struct {
	ZipUrl string `json:"zip_url"`
}

type Output struct {
	Secret string `json:"secret"`
}

// fcrackzip need to install
func (d BruteForceZip) Solve(input string) (interface{}, error) {
	data := new(Data)
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(data.ZipUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create("/tmp/package.zip")
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("fcrackzip", "-b", "-c", "a1", "-l", "4-6", "-u", "/tmp/package.zip")
	cmdOutput, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	re, _ := regexp.Compile("pw == (?P<pwd>.+)")
	matches := re.FindStringSubmatch(string(cmdOutput))
	pwdIndex := re.SubexpIndex("pwd")

	pwd := matches[pwdIndex]
	cmd = exec.Command("unzip", "-P", pwd, "/tmp/package.zip", "secret.txt", "-d", "/tmp")
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	dat, err := os.ReadFile("/tmp/secret.txt")
	if err != nil {
		return nil, err
	}

	output := new(Output)
	output.Secret = strings.TrimSuffix(string(dat), "\n")
	return output, nil
}
