# jsonfmt

Like `gofmt`, but for JSON files.

Usage: `jsonfmt` or `jsonfmt -w` to autofix the issues.

## Install

### macOS

```sh
brew install caarlos0/tap/jsonfmt
```

### ubuntu

```sh
snap install jsonfmt
```

### other linux

Download the `.deb` or `.rpm` from the [releases page][releases].

### Docker

```sh
docker run -v $PWD:/data --workdir /data caarlos0/jsonfmt -h
```

### other

Download the `tar.gz` file from the [releases page][releases] or build from
source.

[releases]: https://github.com/caarlos0/jsonfmt/releases
