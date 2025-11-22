package url

import (
	"errors"
	"net/url"
	"strings"
)

func validateURL(raw string) error {
	if raw == "" {
		return errors.New("URL cannot be empty")
	}

	u, err := url.Parse(raw)
	if err != nil {
		return errors.New("invalid URL format")
	}

	u.Scheme = strings.ToLower(u.Scheme)
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("URL must start with http:// or https://")
	}

	if strings.HasPrefix(raw, "javascript:") ||
		strings.HasPrefix(raw, "data:") ||
		strings.HasPrefix(raw, "vbscript:") {
		return errors.New("unsafe URL scheme")
	}

	return nil
}

func normalizeURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	u.Host = strings.ToLower(u.Host)
	u.Fragment = ""

	if strings.HasSuffix(u.Path, "/") && len(u.Path) > 1 {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}

	return u.String(), nil
}
