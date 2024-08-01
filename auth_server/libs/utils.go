package libs

import (
	"fmt"
	"net/url"
)

func GetAuthUrl(hostURL, port, token string) string {
	if port != "" {
		hostURL = hostURL + ":" + port
	}
	return fmt.Sprintf("%s/auth?token=%s&sign={sign}&public_key_hex={public_key_hex}", hostURL, url.QueryEscape(token))
}

func GetLoginHomeUrl(hostURL, port string) string {
	if port != "" {
		hostURL = hostURL + ":" + port
	}
	return fmt.Sprintf("%s/dashboard", hostURL)
}