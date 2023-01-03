package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
)

func main() {
	
	var newValue string
	flag.StringVar(&newValue, "v", "", "Set the custom value to modify the URLs upon")

	var addParam string
	flag.StringVar(&addParam, "p", "", "Add a custom parameter to the URL")

	var replaceMode bool
	flag.BoolVar(&replaceMode, "r", false, "Replace the value instead of appending it")

	var singleMode bool
	flag.BoolVar(&singleMode, "s", false, "Modify a single parameter at a time")

	var decodeMode bool
	flag.BoolVar(&decodeMode, "d", false, "URL decode the values of the paramters")

	flag.Parse()

	seen := make(map[string]bool)

	// read URLs on stdin, then modify the values in the query string with some user-provided value
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		u, err := url.Parse(sc.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse url %s [%s]\n", sc.Text(), err)
			continue
		}

		// Add the parameter to the URL if the `p` flag is specified
		if addParam != "" {
			if u.RawQuery == "" {
				// No parameters in the URL, so just add the new parameter
				u.RawQuery = addParam + "=" + newValue
			} else {
				// There are already parameters in the URL, so check if the specified parameter exists
				qs := u.Query()
				if _, exists := qs[addParam]; !exists {
					// The parameter doesn't exist, so add it
					qs.Set(addParam, newValue)
					if (decodeMode){
						u.RawQuery, _ = url.QueryUnescape(qs.Encode())
					} else {
						u.RawQuery = qs.Encode()
					}
				}
			}
		}

		param := make([]string, 0)
		for p, _ := range u.Query() {
			param = append(param, p)
		}

		sort.Strings(param)

		key := fmt.Sprintf("%s%s?%s", u.Hostname(), u.EscapedPath(), strings.Join(param, "&"))

		// Only output each host + path + params combination once
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = true

		// Modify the parameters one by one if the `-s` flag is set
		if singleMode {

			// Save the original URL
			originalURL := u.String()

			for i, _ := range param {

				// Restore the original URL
				u, _ = url.Parse(originalURL)

				qs := url.Values{}
				for j, p := range param {
					if i == j {
						if replaceMode {
							qs.Set(p, newValue)
						} else {
							if (u.Query().Get(p) != newValue) {
								qs.Set(p, u.Query().Get(p)+newValue)
							} else {
								qs.Set(p, newValue)
							}
						}
					} else {
						qs.Set(p, u.Query().Get(p))
					}
				}
				if (decodeMode){
					u.RawQuery, _ = url.QueryUnescape(qs.Encode())
				} else {
					u.RawQuery = qs.Encode()
				}

				// Use a buffered writer to write the modified URL to stdout
				w := bufio.NewWriter(os.Stdout)
				fmt.Fprintln(w, u)
				w.Flush()
			}
		} else {
			qs := url.Values{}
			for param, val := range u.Query() {
				if replaceMode {
					qs.Set(param, newValue)
				} else {
					qs.Set(param, val[0]+newValue)
				}
			}

			if (decodeMode){
				u.RawQuery, _ = url.QueryUnescape(qs.Encode())
			} else {
				u.RawQuery = qs.Encode()
			}

			// Use a buffered writer to write the modified URL to stdout
			w := bufio.NewWriter(os.Stdout)
			fmt.Fprintln(w, u)
			w.Flush()
		}
	}
}
