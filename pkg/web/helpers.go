package web

import "net/http"

func CheckIfAvailable(url string) bool {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		return false
	}
	return true
}
