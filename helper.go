package main

import (
	"fmt"
	"io"
	"time"
	"log"
	"net/http"
	"os/exec"
	"encoding/json"
	"strings"
	"runtime"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger = &lumberjack.Logger{
    Filename:   ConfigValue.LogPath + "/log",
    MaxSize:    1,
    MaxBackups: 3,
    MaxAge:     28,
    Compress:   true, 
}

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
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
    logLine := fmt.Sprintf("[%s] %s\n", timestamp, line)

    if ConfigValue.Debug {
        fmt.Print(logLine)
    }

    if _, err := Logger.Write([]byte(logLine)); err != nil {
        log.Printf("Failed to write to log: %v", err)
    }
}

func RunCommands(commands []string) {
	var shell, flag string

	if runtime.GOOS == "windows" {
		shell = "powershell"
		flag = "-Command"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}

	command := strings.Join(commands, " && ")
	cmd := exec.Command(shell, flag, command)
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