package service

import (
	"io/ioutil"
	"strings"

	"os/exec"
)

func ImageCode(filename string) (string, error) {
	cmd := exec.Command("tesseract", filename, "stdout", "digits")
	r, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(body), " \n"), nil
}
