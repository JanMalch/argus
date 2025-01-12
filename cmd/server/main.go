package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/janmalch/argus/internal/config"
	"github.com/janmalch/argus/internal/handler"
	"github.com/janmalch/argus/internal/ui"
)

func run(args []string) error {
	configFile := consumeArgs(args)

	sessionDir, err := os.MkdirTemp("", "argus-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(sessionDir)

	if watcher, err := config.Watch(configFile); err != nil {
		return err
	} else {
		defer watcher.Close()
	}

	conf := config.GetConfig()

	tui := ui.NewTerminalUI(conf.Directory, sessionDir, filepath.Join(conf.Directory, "log.txt"))

	listeners := make([]net.Listener, 0)
	for i, server := range conf.Servers {
		listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", server.Port))
		if err != nil {
			return err
		}
		// usedPort := listener.Addr().(*net.TCPAddr).Port
		srv := handler.NewServer(tui, func() config.Server {
			// NOTE: important to always get a fresh config
			return config.GetConfig().Servers[i]
		})
		listeners = append(listeners, listener)
		go func() {
			if err := http.Serve(listener, srv); err != nil && !errors.Is(err, http.ErrServerClosed) && !errors.Is(err, net.ErrClosed) {
				panic(fmt.Errorf("error listening and serving: %w", err))
			}
		}()
	}

	uiErr := tui.Run()
	errs := make([]error, 0)
	if uiErr != nil {
		errs = append(errs, uiErr)
	}
	for _, l := range listeners {
		errs = append(errs, l.Close())
	}
	return errors.Join(errs...)
}

const VERSION = "0.1.0"

const helpText = `
Usage: argus [<config>]

A convenient proxy server for developers.

Arguments:
  [<config>]    Path to the TOML configuration file. Default: argus.toml
  
Flags:
  -h, --help       Show this help.
  -v, --version    Print Argus version.
  
Example configuration file:

directory = ".argus"

[[server]]
upstream = "https://jsonplaceholder.typicode.com"
port = 3000

[server.response.overwrites]
"^/todos/\\d+" = 403
"GET /comments/6" = "custom_get_6_response.json"
"POST /comments/6" = "custom_post_6_response.json"
"^/images/.+" = "img/16_9.png"

[[server]]
upstream = "https://postman-echo.com"
# you can specify the port explicitly, otherwise +1 to previous server

[server.request.headers]
"X-Test" = "Hi"
"Cache-Control" = "None"

[server.request.parameters]
# combine request headers and query parameters for cache busting
"__cache_buster" = "{{rng.uuid}}"

[server.response.headers]
# add response headers to ease debugging with your application
"X-Argus-ID" = "{{id}}"`

func consumeArgs(args []string) string {
	if slices.Contains(args, "-v") || slices.Contains(args, "--version") {
		fmt.Println("Argus version " + VERSION)
		os.Exit(0)
		return ""
	}

	if slices.Contains(args, "-h") || slices.Contains(args, "--help") {
		fmt.Println(helpText)
		os.Exit(0)
		return ""
	}

	config := "argus.toml"
	for _, a := range args[1:] {
		if a != "" && !strings.HasPrefix(a, "-") {
			config = a
		}
	}
	return config
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
