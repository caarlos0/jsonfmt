# jsonfmt

Like `gofmt`, but for JSON files.

Usage: `jsonfmt` or `jsonfmt -w` to autofix the issues.

## Install

**homebrew**:

```sh
brew install caarlos0/tap/jsonfmt
```

**snapcraft**:

```sh
snap install jsonfmt
```

**docker**:

```sh
docker run -v $PWD:/data --workdir /data caarlos0/jsonfmt -h
```

**deb/rpm**:

Download the `.deb` or `.rpm` from the [releases page][releases] and
install with `dpkg -i` and `rpm -i` respectively.

**manually**:

Download the pre-compiled binaries from the [releases page][releases] or
clone the repo build from source.

[releases]: https://github.com/caarlos0/jsonfmt/releases
