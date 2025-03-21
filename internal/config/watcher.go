package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// TODO: this doesn't need to be a global singleton
var (
	config   *Config = nil
	mu       sync.Mutex
	listener func(c Config) = nil
)

func GetConfig() Config {
	if config == nil {
		log.Panic("Tried accessing the config before it was loaded.")
	}
	return *config
}

func setConfig(path string) error {
	mu.Lock()
	defer mu.Unlock()
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	pconfig, err := parse(file)
	if err != nil {
		return err
	} else {
		if config != nil {
			if len(config.Servers) != len(pconfig.Servers) {
				log.Panic("Add or removing servers at runtime is not supported.")
			}
		}
		config = pconfig
		return nil
	}
}

func SetListener(f func(c Config)) {
	listener = f
}

func Watch(path string) (*fsnotify.Watcher, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		argusHome := os.Getenv("ARGUS_HOME")
		if argusHome == "" {
			return nil, err
		}
		if _, err = os.Stat(filepath.Join(argusHome, path)); errors.Is(err, os.ErrNotExist) {
			return nil, err
		} else {
			path = filepath.Join(argusHome, path)
		}
	}
	if err := setConfig(path); err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return watcher, err
	}

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
				if event.Has(fsnotify.Write) && event.Name == path {
					timer.Reset(time.Millisecond * 100)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				panic(err)

			case <-timer.C:
				setConfig(path)
				if listener != nil {
					listener(GetConfig())
				}
			}

		}
	}()

	err = watcher.Add(filepath.Dir(path))
	return watcher, err
}
