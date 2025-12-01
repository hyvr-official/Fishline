package main

import (
	"fmt"
	"os"
	"io"
	"time"
	"log"
	"net/http"
	"os/exec"
	"encoding/json"
	"strings"
)

func GetPublicIP() string {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return string(body)
}

func WriteLog(line string) {
	f, err := os.OpenFile(ConfigValue.LogPath+string(os.PathSeparator)+"log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer f.Close()

	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatalf("%v", err)
	}

	timestamp := time.Now().In(ist).Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s IST] %s\n", timestamp, line)

	if ConfigValue.Debug {
		fmt.Print(logLine)
	}

	if _, err := f.WriteString(logLine); err != nil {
		log.Fatalf("%v", err)
	}
}

func RunCommand(command string) {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()

	WriteLog(string(output))

	if err != nil {
		WriteLog(err.Error())
	}
}

func JsonFieldExists(field string, cfg Config) bool {
    raw := make(map[string]json.RawMessage)
    file, _ := json.Marshal(cfg)

    json.Unmarshal(file, &raw)
    _, ok := raw[field]

    return ok
}