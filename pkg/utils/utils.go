package utils

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"regexp"
	"strings"
)

// GetWIFIName returns the name of the WIFI network the computer is connected to.
// This function is only implemented for macOS.
func GetWIFIName() string {
	const osxCmd = "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport"
	const osxArgs = "-I"

	cmd := exec.Command(osxCmd, osxArgs)
	stdout := bytes.NewBuffer(nil)
	cmd.Stdout = stdout

	err := cmd.Run()
	if err != nil {
		return ""
	}

	output := strings.TrimSpace(stdout.String())
	r := regexp.MustCompile(`SSID:\s*(.+)`)
	match := r.FindStringSubmatch(output)
	if len(match) < 2 {
		return ""
	}
	name := strings.SplitN(match[1], " ", 2)[1]

	return name
}

// Prettyfy prints the given object in a pretty format.
func Prettyfy(obj interface{}) string {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}
