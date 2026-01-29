package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type GitHubContent struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func GetPath(configPath string, exchange string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	if dir == "/root" || dir == "/home" {
		return "/config/" + configPath + exchange + ".json"
	}
	return os.Getenv("GOPATH") + "/src/github.com/diadata-org/lumina-library/config/" + configPath + exchange + ".json"
}

func GetConfig(directory string, exchange string, branch string) (data []byte, err error) {
	jsonFile, err := readFromRemote(directory, exchange, branch)
	if err != nil {
		log.Errorf("Failed to read %s from remote config for %s: %v", directory, exchange, err)
		jsonFile, err = getLocalConfig(directory, exchange)
		if err != nil {
			log.Errorf("Failed to read %s from local config for %s: %v", directory, exchange, err)
			return nil, err
		}
		log.Infof("Read %s from local config for %s", directory, exchange)
	} else {
		log.Infof("Read %s from remote config for %s", directory, exchange)
	}
	return jsonFile, nil
}

func readFromRemote(directory string, exchange string, branch string) ([]byte, error) {

	URL := "https://api.github.com/repos/diadata-org/decentral-feeder/contents/config/" + directory + "/" + exchange + ".json"
	if branch != "" {
		URL += "?ref=" + url.QueryEscape(branch)
	}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	// Optional authentication
	githubToken := Getenv("GITHUB_TOKEN", "")
	if githubToken != "" {
		log.Info("Set github token for API requests.")
		req.Header.Set("Authorization", "token "+githubToken)
	}

	time.Sleep(350 * time.Millisecond)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// --- Rate limit handling ---
	ratelimitRemaining, err := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
	if err != nil {
		log.Error("rateLimitRemaining for github API calls: ", err)
	}
	log.Info("remaining github API calls: ", ratelimitRemaining)

	// Only applies for anonymous calls or exhausted PAT limits
	if ratelimitRemaining == 0 {
		rateLimitReset, errParseInt := strconv.ParseInt(resp.Header.Get("X-RateLimit-Reset"), 10, 64)
		if errParseInt != nil {
			log.Error("rateLimitReset for github API calls: ", errParseInt)
		}

		timeWait := rateLimitReset - time.Now().Unix()
		if timeWait < 0 {
			timeWait = 0
		}

		log.Warnf("rate limit reached, waiting for refresh in %v", time.Duration(timeWait)*time.Second)
		time.Sleep(time.Duration(timeWait+30) * time.Second)

		resp.Body.Close()
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("GitHub API error: %s\n%s", resp.Status, string(body))
		return nil, err
	}

	var gh GitHubContent
	if err := json.NewDecoder(resp.Body).Decode(&gh); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(gh.Content)
}

func getLocalConfig(directory string, exchange string) (data []byte, err error) {
	configPath, err := getPath2Config(directory)
	if err != nil {
		return nil, err
	}
	path := configPath + exchange + ".json"
	jsonFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return jsonFile, nil
}

func getPath2Config(directory string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("error getting user: %v", err)
	}
	dir := usr.HomeDir
	configPath := "/config/" + directory + "/"
	if dir == "/root" || dir == "/home" {
		return configPath, nil
	}
	return os.Getenv("GOPATH") + "/src/github.com/diadata-org/decentral-feeder" + configPath, nil
}
