# Argus

_A convenient proxy for developers._

## Features

- HTTP-based reverse proxy (not a general purpose network proxy)
- Extremly easy to set up
- TUI to navigate the request timeline
- Skip upstream and respond directly based on file contents
- Record latest exchange per request target (soon™)
  
## Installation

Simply download the stand-alone binary for your platform from the [latest release](https://github.com/JanMalch/argus/releases/latest) and start it. No installation required.

You can also use the following script. Make sure you set the correct `platform`, which must be the name of a [release asset](https://github.com/JanMalch/argus/releases/latest).

```shell
#!/bin/sh
platform="argus-arm64-darwin"
url=$(curl -s https://api.github.com/repos/JanMalch/argus/releases/latest | grep "browser_download_url.*${platform}" | cut -d : -f 2,3 | tr -d '"' | xargs echo -n)
curl -sS -L -o argus "$url"
chmod +x ./argus
# for OS X (optional), see https://superuser.com/a/28400
# xattr -d com.apple.quarantine ./argus
./argus -v
```

Verify it works by running `argus -v` or `argus -h`.

## Usage

Argus behaviour is configured via a TOML configuration file.

```shell
argus # looks for argus.toml in the current working directory
argus path/to/my/config.toml
```

A minimal configuration must look like this:

```toml
[[server]]
upstream = "https://api.example.com"
port = 3000
```

Now simply change your app's upstream to `http://127.0.0.0:3000` and you are good to go!
Argus will pass all requests to `"https://api.example.com"` and display them in your terminal.

See [`argus.toml`](./argus.toml) for what a full configuration might look like.
Any changes besides adding or removing servers are loaded immediately without a server restart.

### Android

If you want to use `argus` for [Android development](https://developer.android.com/studio/run/emulator-networking), you have to use `http://10.0.2.2` in your Android app,
to access your machine's `localhost` from within the emulator. Remember to add the correct port!

## Development

Run with

```
go run cmd/server/main.go
```

and build with

```
GOOS=windows GOARCH=amd64 go build -o bin/argus-amd64.exe cmd/server/main.go
```
