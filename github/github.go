package github

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ListOptions struct {
	// Page for pagination.
	Page int `url:"page,omitempty"`
	// PerPage is the number of items per page. Zero means omit per_page and
	// use the API default.
	PerPage int `url:"per_page,omitempty"`
}

// Response wraps an HTTP response and pagination positions from Link headers.
type Response struct {
	*http.Response

	NextPage  int
	PrevPage  int
	FirstPage int
	LastPage  int
	PerPage   int
}

func parsePagination(r *Response) {
	if r == nil || r.Response == nil {
		return
	}
	links := r.Header.Get("Link")
	if links == "" {
		return
	}
	for _, link := range strings.Split(links, ",") {
		link = strings.TrimSpace(link)
		parts := strings.Split(link, ";")
		if len(parts) < 2 {
			continue
		}
		uPart := strings.Trim(strings.TrimSpace(parts[0]), "<>")
		u, err := url.Parse(uPart)
		if err != nil {
			continue
		}
		q := u.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if pp := q.Get("per_page"); pp != "" {
			if n, err := strconv.Atoi(pp); err == nil && n > 0 {
				r.PerPage = n
			}
		}
		rel := ""
		for _, p := range parts[1:] {
			p = strings.TrimSpace(p)
			if strings.HasPrefix(p, "rel=") {
				rel = strings.Trim(strings.TrimPrefix(p, "rel="), `"`)
			}
		}
		switch rel {
		case "next":
			r.NextPage = page
		case "prev":
			r.PrevPage = page
		case "first":
			r.FirstPage = page
		case "last":
			r.LastPage = page
		}
	}
}
