package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Provider struct {
	watcher  *fsnotify.Watcher
	listener func(c Config)
	config   Config
	Path     string
}

func resolvePath(path string) (string, error) {
	if path == "" {
		return "argus.toml", nil
	}
	if filepath.IsAbs(path) {
		return path, nil
	}
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		argusHome := os.Getenv("ARGUS_HOME")
		if argusHome == "" {
			return "", err
		}
		atHome := filepath.Join(argusHome, path)
		if _, err = os.Stat(atHome); errors.Is(err, os.ErrNotExist) {
			return "", err
		} else {
			return atHome, nil
		}
	}
	return "", err
}

func NewProvider(path string) (*Provider, error) {
	p, err := resolvePath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return nil, fmt.Errorf("failed to determine absolute path from '%s': %w", p, err)
	}
	config, err := parseFile(abs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file from '%s': %w", abs, err)
	}
	return &Provider{config: config, Path: abs}, nil
}

func (p *Provider) Get() Config {
	return p.config
}

func (p *Provider) set(new Config) {
	old := p.Get()
	if len(old.Servers) != len(new.Servers) {
		log.Panic("Add or removing servers at runtime is not supported.")
	}
	p.config = new
}

func (p *Provider) SetListener(f func(c Config)) {
	p.listener = f
}

func (p *Provider) Close() error {
	if p.watcher != nil {
		return p.watcher.Close()
	}
	return nil
}

func (p *Provider) Watch() error {
	if p.watcher != nil {
		return nil
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	p.watcher = watcher

	go func() {
		// Duplicate WRITE events on Windows
		// https://github.com/fsnotify/fsnotify/issues/122#issuecomment-1065925569
		timer := time.NewTimer(time.Millisecond)
		<-timer.C // timer should be expired at first

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) && event.Name == p.Path {
					timer.Reset(time.Millisecond * 100)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				panic(err)

			case <-timer.C:
				c, err := parseFile(p.Path)
				if err != nil {
					log.Panicf("failed to read config file after changes were detected.\n%s", err)
				}
				p.set(c)
				if p.listener != nil {
					p.listener(p.Get())
				}
			}

		}
	}()

	err = watcher.Add(filepath.Dir(p.Path))
	return err
}
