package sitemap

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	link "github.com/elozev/gophercises/link/pkg"
)

const errorMessageFormat = "\nexpected: %+v;\ngot: %+v\n"

var baseUrl, _ = url.Parse("https://localhost:9000")

func TestOnlyInternal(t *testing.T) {
	parsedLinks := []link.Link{
		{Href: "/test"},
		{Href: "https://localhost:9000/test"},
		{Href: "https://localhost:8102/test"},
		{Href: "https://www.google.com/test"},
		{Href: "https://www.yahoo.co.uk/test"},
	}

	expected := []link.Link{
		{Href: "/test"},
		{Href: "https://localhost:9000/test"},
	}

	got := onlyInternal(parsedLinks, baseUrl.String())

	if !reflect.DeepEqual(expected, got) {
		t.Fatalf(errorMessageFormat, expected, got)
	}

	parsedLinks = []link.Link{
		{Href: "/test"},
		{Href: "http://localhost:9000/test"},
		{Href: "http://localhost:9000/2021/09/22/terms-of-reference/"},
		{Href: "https://localhost:8102/test"},
		{Href: "https://www.google.com/test"},
		{Href: "https://www.yahoo.co.uk/test"},
	}

	expected = []link.Link{
		{Href: "/test"},
		{Href: "http://localhost:9000/test"},
		{Href: "http://localhost:9000/2021/09/22/terms-of-reference/"},
	}

	got = onlyInternal(parsedLinks, baseUrl.String())

	if !reflect.DeepEqual(expected, got) {
		t.Fatalf(errorMessageFormat, expected, got)
	}
}

func TestCleanUpUrl(t *testing.T) {
	check := func(exp any, g any) {
		if exp != g {
			t.Fatalf(errorMessageFormat, exp, g)
		}
	}

	base := "https://www.test.com/path_123?q=1"
	expected := "https://www.test.com"
	got := cleanUrl(base)
	check(expected, got)

	base = "http://localhost:9113/hello"
	expected = "http://localhost:9113"
	got = cleanUrl(base)
	check(expected, got)
}

func TestIsLinkInternal(t *testing.T) {
	errorMessage := "expected %+v with host %s to be internal"

	host := "localhost"
	internal := link.Link{
		Href: "/test",
	}

	if !isLinkInternal(internal, host) {
		t.Fatalf(errorMessage, internal, host)
	}

	internal.Href = "http://localhost:9111/test"

	if !isLinkInternal(internal, host) {
		t.Fatalf(errorMessage, internal, host)
	}

	internal.Href = "http://localhost/v1/api/boba"

	if !isLinkInternal(internal, host) {
		t.Fatalf(errorMessage, internal, host)
	}

	// test with https

	internal.Href = "https://localhost/v1/api/https"

	if !isLinkInternal(internal, host) {
		t.Fatalf(errorMessage, internal, host)
	}

	// external links
	external := link.Link{
		Href: "https://google.com",
	}

	if isLinkInternal(external, host) {
		t.Fatalf(errorMessage, external, host)
	}
}

func TestCleanLinkPrefix(t *testing.T) {
	host := "localhost:9090"
	expected := "/test123"
	url := fmt.Sprintf("http://%s%s", host, expected)

	got := cleanLinkPrefix(url, host)
	if expected != got {
		t.Fatalf(errorMessageFormat, expected, got)
	}

	host = "localhost"
	expected = "/test123"
	url = fmt.Sprintf("https://%s%s", host, expected)

	got = cleanLinkPrefix(url, host)
	if expected != got {
		t.Fatalf(errorMessageFormat, expected, got)
	}

	expected = "/test123"
	url = expected

	got = cleanLinkPrefix(url, host)
	if expected != got {
		t.Fatalf(errorMessageFormat, expected, got)
	}
}
