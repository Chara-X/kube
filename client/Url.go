package client

import (
	"net/url"
)

type Url struct {
	Group, Version, Kind, Namespace, Name string
	Query                                 url.Values
}

func (u *Url) String() string {
	var url = "/api/"
	if u.Group != "" {
		url += u.Group + "/"
	}
	url += u.Version + "/"
	if u.Namespace != "" {
		url += "namespaces/" + u.Namespace + "/"
	}
	url += u.Kind
	if u.Name != "" {
		url += "/" + u.Name
	}
	if u.Query != nil {
		url += "?" + u.Query.Encode()
	}
	return url
}
