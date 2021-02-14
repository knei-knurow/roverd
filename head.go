package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// Blink allows to control 2 diodes on the ESP. Wohoo.
func Blink(host string, port string, red string, green string) error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	flag.Parse()

	if host == "" {
		return errors.New("head: host is empty")
	}

	if port == "" {
		return errors.New("head: host is empty")
	}

	espURL := url.URL{
		Scheme: "http",
		Host:   host + ":" + port,
		Path:   "/diode",
	}

	// Construct request.
	q := url.Values{}
	q.Add("green", green)
	q.Add("red", red)

	req, err := http.NewRequest("GET", espURL.String(), nil)
	if err != nil {
		return errors.Wrap(err, "head: create request")
	}
	req.URL.RawQuery = q.Encode()

	// Make request.
	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "head: make request")
	}

	// Get response.
	bodyContent, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "head: read response body")
	}

	log.Printf("head: got response, status: %s, body: %s\n", res.Status, string(bodyContent))
	return nil
}

// MoveServo moves servo to targetRotation.
func MoveServo(targetRotation int) error {
	return errors.New("not implemented")
}
