package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// https://gist.github.com/hSATAC/5343225
// [X](https://<domain>/<version>/<fragment>)
// eg. <domain>/<version>/<fragment>

// https://pxy.fi/5/rsanitize-address
// https://localhost:8080/5/rsanitize-address
// https://releases.llvm.org/5.0.0/tools/clang/docs/DiagnosticsReference.html#rsanitize-address

func main() {

	http.HandleFunc("/", redirect)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Received URL: >%s<\n", r.URL)

	url, parseErr := url.Parse(r.URL.String())
	if parseErr != nil {
		fmt.Printf("Unable to parse received URL: >%s<\n", r.URL)
		http.Error(w, "Unable to parse received URL", http.StatusInternalServerError)
	}

	fmt.Printf("Parsed URL: >%s<\n", url)

	if url.String() == "/" {
		http.ServeFile(w, r, "static/index.html")
	} else {
		newURL, assembleErr := assembleNewURL(url)
		if assembleErr == nil {
			fmt.Printf("Redirecting to: >%s<\n", newURL)
			http.Redirect(w, r, newURL, http.StatusFound)
		} else {
			fmt.Printf("Unable to assemble URL from: >%s< - %s\n", url, assembleErr)
			http.Error(w, "Unable to assemble URL", http.StatusBadRequest)
		}
	}
}

func assembleNewURL(url *url.URL) (string, error) {

	s := strings.SplitN(url.Path, "/", 3)

	// 0 is empty because we split on "/" and the URL begins with "/"
	// 1 == version
	// 2 == fragment

	if len(s) != 3 {
		err := fmt.Errorf("insufficient parts in provided url %+v", s)
		return "", err
	}

	_, err := strconv.Atoi(s[1])
	if err != nil {
		err := fmt.Errorf("first part of url is not a number: %+v", s)
		return "", err
	}

	if s[2] == "" {
		err := fmt.Errorf("second part of url is not a string: %+v", s)
		return "", err
	}

	return fmt.Sprintf("https://releases.llvm.org/%s.0.0/tools/clang/docs/DiagnosticsReference.html#%s", s[1], s[2]), nil
}
