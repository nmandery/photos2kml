package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

var client = &http.Client{
	// for timeout handling in go also refer to
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	Timeout: 10 * 60 * time.Second,
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func getJson(url string, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode >= 400 {
		return errors.New(fmt.Sprintf("upstream server returned HTTP status %d for %s: %s",
			r.StatusCode,
			url,
			http.StatusText(r.StatusCode)))
	}

	return json.NewDecoder(r.Body).Decode(target)
}

type NominatimResponse struct {
	DisplayName string `json:"display_name"`
}

func getNominatimName(p *Photo) (name string, err error) {
	Tell("Reverse-geocoding %s using http://nominatim.openstreetmap.org", p.Filename)
	data := new(NominatimResponse)
	err = getJson(fmt.Sprintf("http://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f&zoom=12", p.Lat, p.Lon), data)
	if err != nil {
		return
	}
	return data.DisplayName, nil
}
