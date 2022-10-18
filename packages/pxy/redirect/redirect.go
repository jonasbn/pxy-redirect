package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"net/url"
	"strings"
)

// https://gist.github.com/hSATAC/5343225
// [X](https://<domain>/<version>/<fragment>)
// eg. <domain>/<version>/<fragment>

// https://pxy.nu/5/rsanitize-address
// https://localhost:8080/5/rsanitize-address
// https://releases.llvm.org/5.0.0/tools/clang/docs/DiagnosticsReference.html#rsanitize-address

func main() {

	sysLog, err := syslog.Dial("tcp", "localhost:1234",
		syslog.LOG_WARNING|syslog.LOG_DAEMON, "pxy")

	log.SetOutput(sysLog)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", redirect)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Received URL: >%s<\n", r.URL)

	url, err := url.Parse(r.URL.String())
	if err != nil {
		log.Fatal(err)
	}

	newURL := assembleNewURL(url)

	fmt.Printf("Redirecting to: >%s<\n", newURL)

	http.Redirect(w, r, newURL, http.StatusFound)
}

func assembleNewURL(url *url.URL) string {

	s := strings.SplitN(url.Path, "/", 3)

	fmt.Printf("%+q\n", s)

	// 0 is empty because we split on "/" and the URL begins with "/"
	// 1 == version
	// 2 == fragment

	return fmt.Sprintf("https://releases.llvm.org/%s.0.0/tools/clang/docs/DiagnosticsReference.html#%s", s[1], s[2])
}
