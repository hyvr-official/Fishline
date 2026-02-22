package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	binaryName := "fishline_test_bin"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	exePath, err := filepath.Abs(filepath.Join("..", binaryName))
	if err != nil {
		t.Fatalf("Failed to resolve absolute path: %v", err)
	}

	buildCmd := exec.Command("go", "build", "-o", exePath, "..")
	buildOut, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, string(buildOut))
	}
	defer os.Remove(exePath)

	configContent := `{
		"port": "18080",
		"logPath": ".",
		"debug": true,
		"commands": {
			"validproject": {
				"main": ["echo integration_test_passed"]
			}
		}
	}`

	configFilePath := filepath.Join("..", "test_config.json")
	if err := os.WriteFile(configFilePath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}
	defer os.Remove(configFilePath)

	defer os.RemoveAll(filepath.Join("..", "log"))

	runCmd := exec.Command(exePath, "-config", "test_config.json")
	runCmd.Dir = ".."

	if err := runCmd.Start(); err != nil {
		t.Fatalf("Failed to start application: %v", err)
	}

	defer func() {
		runCmd.Process.Kill()
		runCmd.Wait()
	}()

	time.Sleep(1 * time.Second)

	baseURL := "http://localhost:18080"
	client := &http.Client{Timeout: 5 * time.Second}

	sendReq := func(method, path string, body interface{}) (*http.Response, string) {
		var reqBody []byte
		if body != nil {
			reqBody, _ = json.Marshal(body)
		}

		req, err := http.NewRequest(method, baseURL+path, bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		return resp, string(respBody)
	}

	t.Run("Method Not Allowed", func(t *testing.T) {
		resp, _ := sendReq(http.MethodGet, "/validproject", nil)
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid JSON Payload", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, baseURL+"/validproject", bytes.NewBuffer([]byte("{invalidjson")))
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", resp.StatusCode)
		}
	})

	t.Run("Empty Project Name", func(t *testing.T) {
		resp, body := sendReq(http.MethodPost, "/", map[string]interface{}{
			"ref": "refs/heads/main",
			"repository": map[string]string{
				"ssh_url": "git@github.com:test/test.git",
			},
		})
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", resp.StatusCode)
		}
		if !strings.Contains(body, "Invalid project name") {
			t.Errorf("expected expected error message in body, got: %s", body)
		}
	})

	t.Run("Invalid Project Name", func(t *testing.T) {
		resp, body := sendReq(http.MethodPost, "/invalidproject", map[string]interface{}{
			"ref": "refs/heads/main",
			"repository": map[string]string{
				"ssh_url": "git@github.com:test/test.git",
			},
		})
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", resp.StatusCode)
		}
		if !strings.Contains(body, "Undefined project name") {
			t.Errorf("expected Undefined project name, got: %s", body)
		}
	})

	t.Run("Invalid Branch", func(t *testing.T) {
		resp, body := sendReq(http.MethodPost, "/validproject", map[string]interface{}{
			"ref": "refs/heads/dev",
			"repository": map[string]string{
				"ssh_url": "git@github.com:test/test.git",
			},
		})
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", resp.StatusCode)
		}
		if !strings.Contains(body, "Undefined branch") {
			t.Errorf("expected Undefined branch, got: %s", body)
		}
	})

	t.Run("Valid Pipeline Execution", func(t *testing.T) {
		resp, _ := sendReq(http.MethodPost, "/validproject", map[string]interface{}{
			"ref": "refs/heads/main",
			"repository": map[string]string{
				"ssh_url": "git@github.com:test/test.git",
			},
		})
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", resp.StatusCode)
		}
	})
}
