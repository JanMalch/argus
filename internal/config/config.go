package config

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

// TOML structs
type rawUI struct {
	Horizontal      bool     `toml:"horizontal"`
	GrowTimeline    int      `toml:"grow_timeline"`
	GrowExchange    int      `toml:"grow_exchange"`
	TimelineColumns []string `toml:"timeline_columns"`
}

type rawRequest struct {
	Headers    map[string]string `toml:"headers"`
	Parameters map[string]string `toml:"parameters"`
}

type rawResponse struct {
	Headers map[string]string `toml:"headers"`
	// any can be string for file, or int for status code
	Overwrites map[string]any `toml:"overwrites"`
}

type rawServer struct {
	Upstream string      `toml:"upstream"`
	Port     *int        `toml:"port"`
	Request  rawRequest  `toml:"request"`
	Response rawResponse `toml:"response"`
}

type rawConfig struct {
	UI        rawUI       `toml:"ui"`
	Directory string      `toml:"directory"`
	Servers   []rawServer `toml:"server"`
}

// sanitized structs

type UI struct {
	Horizontal      bool
	GrowTimeline    int
	GrowExchange    int
	TimelineColumns []string
}

type Overwrite struct {
	Method string
	Regex  *regexp.Regexp
	Exact  string
	File   string
	Status int
}

type Request struct {
	Headers    map[string]string
	Parameters map[string]string
}

type Response struct {
	Headers    map[string]string
	Overwrites []Overwrite
}

type Server struct {
	Upstream *url.URL
	Port     int
	Request  Request
	Response Response
}

type Config struct {
	UI        UI
	Directory string
	Servers   []Server
}

type ServerProvider func() Server

var (
	ErrEndpointEmptyPath = errors.New("an endpoint may not have an empty string for its path")
	ErrNoUpstream        = errors.New("upstream must be defined in configuration file")
	ErrNoServers         = errors.New("no servers configured")
	ErrNoFirstServerPort = errors.New("at least the first server must have a port greater than zero")
	ErrNoServerPort      = errors.New("server port must be greater than 0")
)

func parseRawServer(raw rawServer, fallbackPort int) (*Server, error) {
	if raw.Upstream == "" {
		return nil, ErrNoUpstream
	}
	upstream, err := url.Parse(raw.Upstream)
	if err != nil {
		return nil, err
	}
	var port int
	if raw.Port != nil {
		port = *raw.Port
	} else {
		if fallbackPort == 0 {
			return nil, ErrNoFirstServerPort
		}
		port = fallbackPort
	}
	if port == 0 {
		return nil, ErrNoServerPort
	}
	request := Request{
		Headers:    raw.Request.Headers,
		Parameters: raw.Request.Parameters,
	}

	overwrites := make([]Overwrite, 0)
	for k, v := range raw.Response.Overwrites {
		if k == "" {
			return nil, ErrEndpointEmptyPath
		}
		method := ""
		kspace := strings.Index(k, " ")
		if kspace != -1 {
			method = k[0:kspace]
		}
		path := k[kspace+1:]

		var regex *regexp.Regexp
		var exact string
		if strings.HasPrefix(path, "^") || strings.HasPrefix(path, "(?i)^") {
			regex = regexp.MustCompile(path)
		} else {
			exact = path
		}
		overwrite := Overwrite{
			Method: method,
			Regex:  regex,
			Exact:  exact,
		}
		switch value := v.(type) {
		case int64:
			if code, err := toStatusCode(value); err != nil {
				return nil, err
			} else {
				overwrite.Status = code
			}
		case string:
			overwrite.File = value
		default:
			return nil, fmt.Errorf("cannot parse config because of unknown type %T for endpoint definition %s", v, k)
		}
		overwrites = append(overwrites, overwrite)
	}
	response := Response{
		Headers:    raw.Response.Headers,
		Overwrites: overwrites,
	}

	return &Server{
		Upstream: upstream,
		Port:     port,
		Request:  request,
		Response: response,
	}, nil
}

func parse(tomlReader io.Reader) (*Config, error) {
	var raw rawConfig
	_, err := toml.NewDecoder(tomlReader).Decode(&raw)
	if err != nil {
		return nil, err
	}
	if len(raw.Servers) == 0 {
		return nil, ErrNoServers
	}

	directory := raw.Directory
	if directory == "" {
		directory = ".argus"
	}

	nextPort := 0
	servers := make([]Server, 0)
	for _, rs := range raw.Servers {
		s, err := parseRawServer(rs, nextPort)
		if err != nil {
			return nil, err
		}
		nextPort = s.Port + 1
		servers = append(servers, *s)
	}

	uiTimelineColumns := []string{"ID",
		"start",
		"method",
		"host",
		"request target",
		"end",
		"duration",
		"status_code",
		"status_Text",
	}
	if len(raw.UI.TimelineColumns) > 0 {
		uiTimelineColumns = raw.UI.TimelineColumns
	}

	return &Config{
		Directory: directory,
		Servers:   servers,
		UI: UI{
			Horizontal:      raw.UI.Horizontal,
			GrowTimeline:    max(1, raw.UI.GrowTimeline),
			GrowExchange:    max(1, raw.UI.GrowExchange),
			TimelineColumns: uiTimelineColumns,
		},
	}, nil
}

func toStatusCode(value int64) (int, error) {
	if value >= 600 {
		return 0, fmt.Errorf("invalid status code %d", value)
	}
	c := int(value)
	if http.StatusText(c) == "" {
		return 0, fmt.Errorf("invalid status code %d", value)
	}
	return c, nil
}
