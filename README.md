# Argus

_A convenient proxy for developers._

## Usage

Install the CLI from the [latest release](https://github.com/JanMalch/argus/releases) and start it.

```shell
argus # looks for argus.toml in the current working directory
argus path/to/my/config.toml
```

A minimal configuration must look like this:

```toml
[[server]]
upstream = "https://jsonplaceholder.typicode.com"
```

This will open the server on a random port and log all requests.
Any changes to the configuration file are loaded immediately without a server restart.

Checkout [`argus.toml`](./argus.toml) for what a full configuration might look like.

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
