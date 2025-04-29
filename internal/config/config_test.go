package config

import (
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseValidConfig(t *testing.T) {
	toml := strings.NewReader(`
directory = ".argus"

[[server]]
upstream = "https://jsonplaceholder.typicode.com"
port = 3000

[server.request.headers]
"X-Test" = "Hi"
"Cache-Control" = "None"

[server.request.parameters]
"__argus_id" = "{{id}}"
"__cache_buster" = "{{rng.uuid}}"

[server.response.headers]
"X-Argus-ID" = "{{id}}"

[server.response.overwrites]
"^/todos/\\d+" = 403
"GET /comments/6" = "custom_get_6.json"
"POST /comments/6" = "custom_post_6.json"
"/images/banner.png" = "dummy_landscape.png"

[[server]]
upstream = "https://postman-echo.com"

`)

	placeholder, err := url.Parse("https://jsonplaceholder.typicode.com")
	require.NoError(t, err)
	echo, err := url.Parse("https://postman-echo.com")
	require.NoError(t, err)

	conf, err := parse(toml)
	if assert.NoError(t, err) {
		expected := Config{
			Directory: ".argus",
			Servers: []Server{
				{
					Upstream: placeholder,
					Port:     3000,
					Request: Request{
						Headers: map[string]string{
							"X-Test":        "Hi",
							"Cache-Control": "None",
						},
						Parameters: map[string]string{
							"__argus_id":     "{{id}}",
							"__cache_buster": "{{rng.uuid}}",
						},
					},
					Response: Response{
						Headers: map[string]string{
							"X-Argus-ID": "{{id}}",
						},
						// FIXME: non-deterministic test because map order is non-deterministic ...
						Overwrites: []Overwrite{
							{
								Method: "",
								Regex:  regexp.MustCompile(`^/todos/\d+`),
								Exact:  "",
								File:   "",
								Status: 403,
							},
							{
								Method: "GET",
								Regex:  nil,
								Exact:  "/comments/6",
								File:   "custom_get_6.json",
								Status: 0,
							},
							{
								Method: "POST",
								Regex:  nil,
								Exact:  "/comments/6",
								File:   "custom_post_6.json",
								Status: 0,
							},
							{
								Method: "",
								Regex:  nil,
								Exact:  "/images/banner.png",
								File:   "dummy_landscape.png",
								Status: 0,
							},
						},
					},
				},
				{
					Upstream: echo,
					Port:     3001,
					Response: Response{
						Overwrites: make([]Overwrite, 0),
					},
				},
			},
		}
		assert.EqualValues(t, expected, conf)
	}
}
