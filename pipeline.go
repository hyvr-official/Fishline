package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type repository struct {
	GitSSHURL string `json:"git_ssh_url"`
	SSHURL string `json:"ssh_url"`
}

type payload struct {
	Ref        string `json:"ref"`
	repository `json:"repository"`
}

func PipelineHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(map[string]string{
			"error":  "Only POST method is allowed",
			"status": "405",
		})

		return
	}

	var payload payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteLog("Pipeline called with a invalid JSON")

		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(map[string]string{
			"error":  "Invalid payload JSON",
			"status": "405",
		})

		return
	}

	var payloadRepositoryURL string

	if payload.repository.GitSSHURL != "" {
		payloadRepositoryURL = payload.repository.GitSSHURL
	}

	if payload.repository.SSHURL != "" {
		payloadRepositoryURL = payload.repository.SSHURL
	}

	WriteLog("")
	WriteLog(fmt.Sprintf("Pipeline call started by %s with %s branch payload", payloadRepositoryURL, payload.Ref))

	project := strings.Trim(r.URL.Path, "/")

	if project == "" {
		WriteLog("Pipeline called with a invalid project name")

		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(map[string]string{
			"error":  "Invalid project name",
			"status": "405",
		})

		return
	}

	branch := strings.TrimPrefix(payload.Ref, "refs/heads/")
	
	_, checkProject := ConfigValue.Commands[project]
	if !checkProject {
		WriteLog("Pipeline called with a undefined project name")

		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(map[string]string{
			"error":  "Undefined project name",
			"status": "405",
		})

		return
	}

	_, checkBranch := ConfigValue.Commands[project][branch]
	if !checkBranch {
		WriteLog("Pipeline called with a undefined branch")

		w.WriteHeader(http.StatusMethodNotAllowed)

		json.NewEncoder(w).Encode(map[string]string{
			"error":  "Undefined branch",
			"status": "405",
		})

		return
	}

	WriteLog("Starting pipeline commands execution")

	for _, command := range ConfigValue.Commands[project][branch] {
		WriteLog(fmt.Sprintf("Running the command %s", command))

		RunCommand(command)
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "200",
	})
}