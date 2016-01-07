// instagrab fetches all instagram images belonging to an account
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	url := "http://instagram.com/"
	for _, profile := range os.Args[1:] {
		resp, err := http.Get(url + profile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
			os.Exit(1)
		}
		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: reading %s: %v\n", os.Args[0], url, err)
			os.Exit(1)
		}
		r := regexp.MustCompile("http.+?scontent.cdninstagram.com.+?jpg")
		//r, _ := regexp.Compile("http.+scontent.cdninstagram.com.+jpg")
		matches := r.FindAllString(string(b), -1)

		for _, match := range matches {
			match = strings.Replace(match, "\\/", "/", -1)
			splitted := strings.Split(match, "/")
			fname := splitted[len(splitted)-1]
			fmt.Println("fetching:", fname)

			fetchToFile(match, fname)
		}
	}
}

func fetchToFile(url, fname string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	fd, err := os.Create(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	defer fd.Close()
	defer resp.Body.Close()

	_, err = io.Copy(fd, resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	fmt.Println("Written:", fname)
}
