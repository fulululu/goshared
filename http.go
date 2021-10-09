package goshared

import "net/url"

func ParseQuery(str string) (url.Values, error) {
	URLStr, e := url.QueryUnescape(str)
	if e != nil {
		return nil, e
	}
	u, e := url.Parse(URLStr)
	if e != nil {
		return nil, e
	}
	return u.Query(), nil
}
