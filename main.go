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
	flag.StringVar(&newValue, "v", "", "Value to modify the parameters upon")

	var addParam string
	flag.StringVar(&addParam, "a", "", "Add custom parameters to the URLs, comma separated")

	var rmParam string
	flag.StringVar(&rmParam, "r", "", "Remove parameters from the URLs, comma separated")

	var overrideMode bool
	flag.BoolVar(&overrideMode, "o", false, "Replace the value instead of appending")

	var multiMode bool
	flag.BoolVar(&multiMode, "m", false, "Modify all parameters at once")

	var decodeMode bool
	flag.BoolVar(&decodeMode, "d", false, "URLdecode the values of the parameters")

	var keepAllUrls bool
	flag.BoolVar(&keepAllUrls, "k", false, "Keep the URLs with no parameters")


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

		// Remove parameters from the URL if the `r` flag is specified
		delParams := strings.Split(rmParam, ",")
		if u.RawQuery != "" && rmParam != "" {
			qs := u.Query()
			for _, param := range delParams {
				if _, exists := qs[param]; exists {
					// The parameter exists, so delete it
					delete(qs, param)
				}
			}
			// Rebuild the query string without the deleted parameter(s)
			u.RawQuery = ""
			for p, v := range qs {
				u.RawQuery += p
				if v[0] != "" {
					u.RawQuery += "=" + v[0]
				}
				u.RawQuery += "&"
			}
			// Trim the trailing ampersand
			u.RawQuery = u.RawQuery[:len(u.RawQuery)-1]
			if decodeMode {
				u.RawQuery, _ = url.QueryUnescape(u.RawQuery)
			}
		}

		// Add parameters to the URL if the `a` flag is specified
		newParams := strings.Split(addParam, ",")
		if addParam != "" {
			for _, param := range newParams {
				if u.RawQuery == "" {
					// No parameters in the URL, so just add the new parameter
					u.RawQuery = param + "=" + newValue
				} else {
					// There are already parameters in the URL, so check if the specified parameter exists
					qs := u.Query()
					if _, exists := qs[param]; !exists {
						// The parameter doesn't exist, so add it
						qs.Set(param, newValue)

						u.RawQuery = qs.Encode()
						if decodeMode {
							u.RawQuery, _ = url.QueryUnescape(u.RawQuery)
						}
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

		// Only output each host + path + newParams combination once
		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = true

		if multiMode {
			qs := url.Values{}
			for p, val := range u.Query() {
				if overrideMode {
					qs.Set(p, newValue)
				} else {
					if (p != addParam) {
						qs.Set(p, val[0]+newValue)
					} else {
						qs.Set(p, newValue)
					}
				}
			}

			u.RawQuery = qs.Encode()
			if decodeMode {
				u.RawQuery, _ = url.QueryUnescape(u.RawQuery)
			}

			if !(u.RawQuery == "") || keepAllUrls {
				w := bufio.NewWriter(os.Stdout)
				fmt.Fprintln(w, u)
				w.Flush()
			}
			
		} else {
			// Save the original URL
			originalURL := u.String()
			lastUrl := ""

			for i, _ := range param {

				// Restore the original URL
				u, _ = url.Parse(originalURL)

				qs := url.Values{}

				for j, p := range param {
					if i == j {
						if overrideMode {
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

				u.RawQuery = qs.Encode()
				if decodeMode {
					u.RawQuery, _ = url.QueryUnescape(u.RawQuery)
				}

				// Make sure no duplicates
				if lastUrl != u.String() {
					w := bufio.NewWriter(os.Stdout)
					fmt.Fprintln(w, u)
					w.Flush()
				}
				lastUrl = u.String()
			}
			if keepAllUrls && (u.RawQuery == "") {
				w := bufio.NewWriter(os.Stdout)
				fmt.Fprintln(w, u)
				w.Flush()
			}
		}
	}
}
