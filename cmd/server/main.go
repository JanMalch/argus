package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/janmalch/argus/internal/config"
	"github.com/janmalch/argus/internal/handler"
	"github.com/janmalch/argus/internal/tui"
)

func run(configFile string) error {
	if watcher, err := config.Watch(configFile); err != nil {
		return err
	} else {
		defer watcher.Close()
	}

	conf := config.GetConfig()
	directory := filepath.Join(filepath.Dir(configFile), conf.Directory)
	tuiApp := tui.NewApp(directory, conf.UI)
	config.SetListener(func(c config.Config) {
		tuiApp.SetUI(c.UI)
	})

	listeners := make([]net.Listener, 0)
	for i, server := range conf.Servers {
		listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", server.Port))
		if err != nil {
			return err
		}
		// usedPort := listener.Addr().(*net.TCPAddr).Port
		srv := handler.NewServer(tuiApp, func() config.Server {
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

	uiErr := tuiApp.Run()
	errs := make([]error, 0)
	if uiErr != nil {
		errs = append(errs, uiErr)
	}
	for _, l := range listeners {
		errs = append(errs, l.Close())
	}
	return errors.Join(errs...)
}

const VERSION = "0.4.0"

var cli struct {
	ConfigFile string           `arg:"" default:"argus.toml" type:"path" help:"Path to the configuration TOML file. Default is \"argus.toml\""`
	Version    kong.VersionFlag `short:"v" name:"version" help:"Print version information and quit"`
}

const helpExtraText = `A convenient proxy server for developers.

Example configuration file: https://github.com/JanMalch/argus/blob/main/argus.toml`

func main() {
	kong.Parse(&cli,
		kong.Name("argus"),
		kong.Description(helpExtraText),
		kong.Vars{
			"version": VERSION,
		},
	)

	if err := run(cli.ConfigFile); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
