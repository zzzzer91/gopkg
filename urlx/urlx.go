package urlx

import "net/url"

// UnescapeURL unescapes url. If there's an error, return the original URL.
func UnescapeURL(rawURL string) string {
	if u, err := url.QueryUnescape(rawURL); err == nil {
		return u
	}
	return rawURL
}
