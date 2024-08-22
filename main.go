package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
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

	var onlyParamURLs bool
	flag.BoolVar(&onlyParamURLs, "p", false, "Only keep URLs with parameters")

	var splitMode bool
	flag.BoolVar(&splitMode, "s", false, "Split URL into all path levels")

	var blacklistMode bool
	flag.BoolVar(&blacklistMode, "b", false, "Enable blacklist to remove static URLs")

	var ext string
	flag.StringVar(&ext, "e", "", "Blacklist regex string (default is static extensions)")

	flag.Parse()

	if ext == "" {
		ext = `(?i)\.(png|apng|bmp|gif|ico|cur|jpg|jpeg|jfif|pjp|pjpeg|svg|tif|tiff|webp|xbm|3gp|aac|flac|mpg|mpeg|mp3|mp4|m4a|m4v|m4p|oga|ogg|ogv|mov|wav|webm|eot|woff|woff2|ttf|otf|css)(?:\?|#|$)`
	}

	seen := make(map[string]bool)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		u, err := url.Parse(sc.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse url %s [%s]\n", sc.Text(), err)
			continue
		}

		if blacklistMode && isBlacklisted(u.String(), ext) {
			continue
		}

		if splitMode {
			processPathLevels(u, seen, onlyParamURLs, blacklistMode, ext)
			// Process the original URL with parameters
			processURL(u, seen, newValue, addParam, rmParam, overrideMode, multiMode, decodeMode, onlyParamURLs, splitMode)
		} else {
			processURL(u, seen, newValue, addParam, rmParam, overrideMode, multiMode, decodeMode, onlyParamURLs, splitMode)
		}
	}
}

func processPathLevels(u *url.URL, seen map[string]bool, onlyParamURLs, blacklistMode bool, ext string) {
	paths := strings.Split(strings.Trim(u.Path, "/"), "/")
	for i := 0; i < len(paths); i++ {
		subPath := "/" + strings.Join(paths[:i+1], "/")
		subURL := *u // Create a copy of the URL
		subURL.Path = subPath
		subURL.RawQuery = "" // Remove query parameters for split URLs

		if blacklistMode && isBlacklisted(subURL.String(), ext) {
			continue
		}

		key := fmt.Sprintf("%s%s", subURL.Hostname(), subURL.EscapedPath())
		if _, exists := seen[key]; !exists {
			seen[key] = true
			outputURL(&subURL, onlyParamURLs, true)
		}
	}
}

func processURL(u *url.URL, seen map[string]bool, newValue, addParam, rmParam string, overrideMode, multiMode, decodeMode, onlyParamURLs, splitMode bool) {
	originalQuery := u.RawQuery

	// Remove parameters
	if u.RawQuery != "" && rmParam != "" {
		removeParameters(u, strings.Split(rmParam, ","), decodeMode)
	}

	// Add parameters
	if addParam != "" {
		addParameters(u, strings.Split(addParam, ","), newValue, decodeMode)
	}

	param := make([]string, 0)
	for p := range u.Query() {
		param = append(param, p)
	}
	sort.Strings(param)

	key := fmt.Sprintf("%s%s?%s", u.Hostname(), u.EscapedPath(), strings.Join(param, "&"))

	if _, exists := seen[key]; exists {
		return
	}
	seen[key] = true

	if multiMode {
		modifyAllParameters(u, newValue, addParam, overrideMode, decodeMode)
		outputURL(u, onlyParamURLs, splitMode)
	} else if len(param) > 0 {
		modifyParametersIndividually(u, param, newValue, overrideMode, decodeMode, onlyParamURLs, splitMode)
	} else {
		// If there are no parameters and we're not in split mode, output the URL
		if !splitMode {
			outputURL(u, onlyParamURLs, splitMode)
		}
	}

	// Restore the original query
	u.RawQuery = originalQuery
}

func removeParameters(u *url.URL, delParams []string, decodeMode bool) {
	qs := u.Query()
	for _, param := range delParams {
		qs.Del(param)
	}
	u.RawQuery = encodeQueryWithoutEmptyValues(qs)
	if decodeMode {
		u.RawQuery, _ = url.QueryUnescape(u.RawQuery)
	}
}

func addParameters(u *url.URL, newParams []string, newValue string, decodeMode bool) {
	for _, param := range newParams {
		qs := u.Query()
		if _, exists := qs[param]; !exists {
			qs.Set(param, newValue)
			u.RawQuery = encodeQueryWithoutEmptyValues(qs)
			if decodeMode {
				u.RawQuery, _ = url.QueryUnescape(u.RawQuery)
			}
		}
	}
}

func modifyAllParameters(u *url.URL, newValue, addParam string, overrideMode, decodeMode bool) {
	qs := u.Query()
	for p, val := range qs {
		if overrideMode {
			qs.Set(p, newValue)
		} else {
			if p != addParam && len(val) > 0 {
				qs.Set(p, val[0]+newValue)
			} else {
				qs.Set(p, newValue)
			}
		}
	}
	u.RawQuery = encodeQueryWithoutEmptyValues(qs)
	if decodeMode {
		u.RawQuery, _ = url.QueryUnescape(u.RawQuery)
	}
}

func modifyParametersIndividually(u *url.URL, param []string, newValue string, overrideMode, decodeMode, onlyParamURLs, splitMode bool) {
	originalURL := u.String()
	lastUrl := ""

	for i := range param {
		u, _ := url.Parse(originalURL)
		qs := u.Query()

		for j, p := range param {
			if i == j {
				if overrideMode {
					qs.Set(p, newValue)
				} else {
					if qs.Get(p) != newValue {
						if qs.Get(p) != "" {
							qs.Set(p, qs.Get(p)+newValue)
						} else {
							qs.Set(p, newValue)
						}
					}
				}
			}
		}

		u.RawQuery = encodeQueryWithoutEmptyValues(qs)
		if decodeMode {
			u.RawQuery, _ = url.QueryUnescape(u.RawQuery)
		}

		if lastUrl != u.String() {
			outputURL(u, onlyParamURLs, splitMode)
		}
		lastUrl = u.String()
	}
}

func encodeQueryWithoutEmptyValues(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			if v != "" {
				buf.WriteByte('=')
				buf.WriteString(url.QueryEscape(v))
			}
		}
	}
	return buf.String()
}

func outputURL(u *url.URL, onlyParamURLs, splitMode bool) {
	if !onlyParamURLs || u.RawQuery != "" {
		w := bufio.NewWriter(os.Stdout)
		fmt.Fprintln(w, u)
		w.Flush()
	}
}

func isBlacklisted(raw, ext string) bool {
	r, err := regexp.Compile(ext)
	if err != nil {
		return false
	}
	return r.MatchString(raw)
}
