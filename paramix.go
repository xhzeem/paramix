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
	var replaceMode bool
	flag.BoolVar(&replaceMode, "r", false, "Replace the value instead of appending it")

	var singleMode bool
	flag.BoolVar(&singleMode, "s", false, "Modify the parameters one by one")

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

		pp := make([]string, 0)
		for p, _ := range u.Query() {
			pp = append(pp, p)
		}
		sort.Strings(pp)

		key := fmt.Sprintf("%s%s?%s", u.Hostname(), u.EscapedPath(), strings.Join(pp, "&"))

		// Only output each host + path + params combination once
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = true

		// Modify the parameters one by one if the `-s` flag is set
		if singleMode {
			// Save the original URL
			originalURL := u.String()

			for i, _ := range pp {
				// Restore the original URL
				u, err = url.Parse(originalURL)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to parse url %s [%s]\n", originalURL, err)
					continue
				}

				qs := url.Values{}
				for j, p := range pp {
					if i == j {
						if replaceMode {
							qs.Set(p, flag.Arg(0))
						} else {
							qs.Set(p, u.Query().Get(p)+flag.Arg(0))
						}
					} else {
						qs.Set(p, u.Query().Get(p))
					}
				}
				u.RawQuery = qs.Encode()
				fmt.Printf("%s\n", u)
			}
		} else {
			qs := url.Values{}
			for param, vv := range u.Query() {
				if replaceMode {
					qs.Set(param, flag.Arg(0))
				} else {
					qs.Set(param, vv[0]+flag.Arg(0))
				}
			}

			u.RawQuery = qs.Encode()

			fmt.Printf("%s\n", u)
		}
	}
}
