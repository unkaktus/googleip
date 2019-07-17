// googleip.go - determine publicly routable IP using Google services.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of googleip, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package googleip

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
)

// GetIP determines publicly routable IP using Google services.
func GetIP(t http.RoundTripper) (*net.IP, error) {
	return getIPSorryPage(t)
}

// Fetch a YouTube video page and extract IP from it.
func getIPYouTubeVideo(t http.RoundTripper, videoID string) (*net.IP, error) {
	if videoID == "" {
		// Use a video that will higly unlikely be ever removed ("Me at the zoo").
		videoID = "jNQXAC9IVRw"
	}
	u := "https://www.youtube.com/watch?v=" + videoID

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := t.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	re, err := regexp.Compile(`Find(\%26ip\%3D[\d\.]+)|(\d{1,3}\.){3}\d{1,3}`)
	if err != nil {
		return nil, err
	}
	ips := re.FindAll(body, 10)
	var ipa net.IP
	for _, ip := range ips {
		ipa = net.ParseIP(string(ip))
		if ipa != nil {
			break
		}
	}
	if ipa == nil {
		return nil, errors.New("no valid IP found")
	}
	return &ipa, nil
}

// Go to Google Sorry Page and grab IP from there.
func getIPSorryPage(t http.RoundTripper) (*net.IP, error) {
	u := "https://www.google.com/sorry/index"

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := t.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusTooManyRequests, http.StatusServiceUnavailable:
	default:
		return nil, errors.New("unexpected status code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	re, err := regexp.Compile(`Find(IP\saddress:\s[\d\.]+)|(\d{1,3}\.){3}\d{1,3}`)
	if err != nil {
		return nil, err
	}
	ips := re.FindAll(body, 10)
	var ipa net.IP
	for _, ip := range ips {
		ipa = net.ParseIP(string(ip))
		if ipa != nil {
			break
		}
	}
	if ipa == nil {
		return nil, errors.New("no valid IP found")
	}
	return &ipa, nil
}
