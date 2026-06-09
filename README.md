# jsonfmt

Like `gofmt`, but for JSON files.

Usage: `jsonfmt` or `jsonfmt -w` to autofix the issues.

## Install

**homebrew**:

```sh
brew install jsonfmt
```

**docker**:

```sh
docker run -v $PWD:/data --workdir /data caarlos0/jsonfmt -h
```

**apt**:

```sh
echo 'deb [trusted=yes] https://repo.caarlos0.dev/apt/ /' | sudo tee /etc/apt/sources.list.d/caarlos0.list
sudo apt update
sudo apt install jsonfmt
```

**yum**:

```sh
echo '[caarlos0]
name=caarlos0
baseurl=https://repo.caarlos0.dev/yum/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/caarlos0.repo
sudo yum install jsonfmt
```

**deb/rpm**:

Download the `.deb` or `.rpm` from the [releases page][releases] and
install with `dpkg -i` and `rpm -i` respectively.

**manually**:

Download the pre-compiled binaries from the [releases page][releases] or
clone the repo build from source.

[releases]: https://github.com/caarlos0/jsonfmt/releases


## Stargazers over time

[![Stargazers over time](https://starchart.cc/caarlos0/jsonfmt.svg)](https://starchart.cc/caarlos0/jsonfmt)

