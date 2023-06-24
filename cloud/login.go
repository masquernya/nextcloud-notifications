package cloud

import (
	"context"
	"encoding/json"
	"github.com/masquernya/nextcloud-notifications/config"
	"github.com/masquernya/nextcloud-notifications/storage"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func (c *Cloud) IsLoginRequired() bool {
	if storage.Get().LoginPassword == "" {
		return true
	}
	return false
}

type LoginPollData struct {
	Token    string `json:"token"`
	Endpoint string `json:"endpoint"`
}

type LoginResponse struct {
	Poll     *LoginPollData `json:"poll"`
	LoginUrl string         `json:"login"`
}

func (c *Cloud) RequestLogin() *LoginResponse {
	for {
		req, err := http.NewRequest("POST", config.Get().CloudUrl+"/index.php/login/v2", nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("user-agent", "todo-notifications/1.0")
		res, err := c.client.Do(req)
		if err != nil {
			panic(err)
		}
		if res.StatusCode != 200 {
			log.Fatal("Unexpected status code: ", res.StatusCode)
		}
		bits, _ := io.ReadAll(res.Body)
		var loginResponse LoginResponse
		err = json.Unmarshal(bits, &loginResponse)
		if err != nil {
			panic(err)
		}
		return &loginResponse
	}
}

type PollLoginResponse struct {
	Server   string `json:"server"`
	Username string `json:"loginName"`
	Password string `json:"appPassword"`
}

func (c *Cloud) PollLogin(poll *LoginPollData) (*PollLoginResponse, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute*10)) // actual max len is 20m
	for {
		time.Sleep(time.Second * 2)
		if ctx.Err() != nil {
			cancel()
			return nil, ctx.Err()
		}
		// curl -X POST https://cloud.example.com/login/v2/poll -d "token=mQUYQdffOSAMJYtm8pVpkOsVqXt5hglnuSpO5EMbgJMNEPFGaiDe8OUjvrJ2WcYcBSLgqynu9jaPFvZHMl83ybMvp6aDIDARjTFIBpRWod6p32fL9LIpIStvc6k8Wrs1"
		body := strings.NewReader("token=" + poll.Token)
		req, err := http.NewRequest("POST", poll.Endpoint, body)
		if err != nil {
			panic(err)
		}
		req.Header.Set("user-agent", "todo-notifications/1.0")
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
		res, err := c.client.Do(req)
		if err != nil {
			panic(err)
		}
		if res.StatusCode == 404 {
			//log.Info("polling for login...")
			continue
		}
		bits, _ := io.ReadAll(res.Body)
		if res.StatusCode != 200 {
			os.WriteFile("poll.html", bits, 0644)
			log.Fatal("Unexpected status code:", res.StatusCode, "body", string(bits))
		}
		var pollResponse PollLoginResponse
		err = json.Unmarshal(bits, &pollResponse)
		if err != nil {
			panic(err)
		}
		cancel()
		return &pollResponse, nil
	}
}
