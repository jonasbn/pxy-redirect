package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// https://gist.github.com/hSATAC/5343225
// [X](https://<domain>/<version>/<fragment>)
// eg. <domain>/<version>/<fragment>

// https://pxy.fi/5/rsanitize-address
// https://localhost:8080/5/rsanitize-address
// https://releases.llvm.org/5.0.0/tools/clang/docs/DiagnosticsReference.html#rsanitize-address

func main() {

	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)

	http.HandleFunc("/", redirect)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {

	log.Infof("Received URL: >%s<\n", r.URL)

	url, parseErr := url.Parse(r.URL.String())
	if parseErr != nil {
		log.Errorf("Unable to parse received URL: >%s<\n", r.URL)
		http.Error(w, "Unable to parse received URL", http.StatusInternalServerError)

		return
	}

	log.Debugf("Parsed URL: >%s<\n", url)

	if url.String() == "/robots.txt" {
		http.ServeFile(w, r, "static/robots.txt")
		return
	}

	if url.String() == "/favicon.ico" {
		http.ServeFile(w, r, "static/favicon.ico")
		return
	}

	if url.String() == "/" {
		http.ServeFile(w, r, "static/index.html")
		return
	}

	newURL, assembleErr := assembleNewURL(url)
	if assembleErr == nil {
		log.Infof("Redirecting to: >%s<\n", newURL)
		http.Redirect(w, r, newURL, http.StatusFound)
	} else {
		log.Errorf("Unable to assemble URL from: >%s< - %s\n", url, assembleErr)
		http.Error(w, "Unable to assemble URL", http.StatusBadRequest)
	}
}

func assembleNewURL(url *url.URL) (string, error) {

	s := strings.SplitN(url.Path, "/", 3)

	log.Debugf("Parsed following parts: >%#v<\n", s)

	// 0 is empty because we split on "/" and the URL begins with "/"
	// 1 == version
	// 2 == fragment

	if len(s) != 3 {
		err := fmt.Errorf("insufficient parts in provided url %q", s)
		return "", err
	}

	_, err := strconv.Atoi(s[1])
	if err != nil {
		err := fmt.Errorf("first part of url is not a number: %q", s)
		return "", err
	}

	if s[2] == "" {
		err := fmt.Errorf("second part of url is not a string: %q", s)
		return "", err
	}

	return fmt.Sprintf("https://releases.llvm.org/%s.0.0/tools/clang/docs/DiagnosticsReference.html#%s", s[1], s[2]), nil
}
