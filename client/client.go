package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

const (
	loginHandle     = "/login"
	manageHandle    = "/wlmacflt.cmd?"
	sessionKeyRegex = "sessionKey=(\\d+)"
)

type Client struct {
	host   string
	client *http.Client
}

type ClientError struct {
	source string
	err    error
}

func (err *ClientError) Error() string {
	return fmt.Sprintf("[%s]: %s", err.source, err.err.Error())
}

func Login(ip string, username string, password string) (c *Client, e error) {
	host := "http://" + ip
	params := url.Values{"username": {username}, "password": {password}}

	jar, err := cookiejar.New(nil)
	client := http.Client{Jar: jar}

	resp, err := client.PostForm(host+"/login", params)
	if err != nil {
	}
	defer resp.Body.Close()

	return &Client{host, &client}, nil
}

func (c *Client) getSessionKey() (k string, e error) {
	resp, err := c.client.Get(c.host + manageHandle + "action=view")
	if err != nil {
		return "", &ClientError{"getSessionKey:http.Get", err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	regex := regexp.MustCompile(sessionKeyRegex)
	keyRes := regex.FindStringSubmatch(string(body))

	if len(keyRes) < 2 {
		return "", &ClientError{"getSessionKey", errors.New("no session key!")}
	}

	return keyRes[1], nil
}

func (c *Client) BanMac(mac string) (e error) {
	key, err := c.getSessionKey()
	if err != nil {
		return err
	}

	params := url.Values{
		"action":       {"add"},
		"wlFltMacAddr": {mac},
		"wlSyncNvram":  {"1"},
		"sessionKey":   {key},
	}.Encode()

	resp, err := c.client.Get(c.host + manageHandle + params)
	if err != nil {
		return &ClientError{"BanMac:http.Get", err}
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) UnbanMac(mac string) (e error) {
	key, err := c.getSessionKey()
	if err != nil {
		return err
	}

	params := url.Values{
		"action":     {"remove"},
		"rmLst":      {mac},
		"sessionKey": {key},
	}.Encode()

	resp, err := c.client.Get(c.host + manageHandle + params)
	if err != nil {
		return &ClientError{"UnbanMac:http.Get", err}
	}
	defer resp.Body.Close()

	return nil
}
