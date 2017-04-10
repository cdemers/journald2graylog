# journald2graylog

[![Build Status](https://travis-ci.org/cdemers/journald2graylog.svg?branch=master)](https://travis-ci.org/cdemers/journald2graylog)

Command line tool to help forward _systemd's_ _journald's_ logs to a [_Graylog_](https://www.graylog.org/) server taking advantage of the descriptive _GELF_ format.

The _journald2graylog_ command expects it's paramters to be provided as the environment variables, making it well suited for [_Docker_](https://www.docker.com/) or _systemd_ driven environments, and _PaaS_ platforms like [_Heroku_](https://www.heroku.com/) and the [_Twelve-Factor App_](http://12factor.net/config) approach.

## Usage

To use _journald2graylog_, you simply pipe the output of _journalctl_, while enabling it's _JSON_ output format, into the _jourald2graylog_ command.  It can be as simple this: `journalctl -o json | journald2graylog`, but usually you will require and want to provide more parameters.

Note that _journald2graylog_ **only supports UDP** for now, having TCP might be cool, but it's not in our short term plans.

There are four configuration parameters:

* The `J2G_HOSTNAME` is the _hostname_ or _IP_ of your _Graylog_ server, it has no default and **MUST** be specified.
* The `J2G_PORT` is the port of the **UDP GELF** input of the _Graylog_ server, it will default to `12201`, but this value will almost always differ depending on your _Graylog_ configuration, so you will most likely have to look it up in your own _Graylog_ server.
* The `J2G_PACKET_SIZE` is the maximum size of the TCP/IP packets you can use between the source (_journald2graylg_) and the destination (your _Graylog_ server). This will vary depending on your network capabilities, but the default value of _1420_ will be appropriate in the vast majority of situations.
* The `J2G_BLACKLIST` is a list containing regex identifying logs that must not be sent to _Graylog_, separated by a semicolon (`;`).

You can add debugging by specifying the `--verbose` (also `-v`) flag, it will display the configuration parameters sent to journald2graylog in stdout

Note that from version 0.2.0 onward, _journald2graylog_ will now exit if there is a network error, instead of looping forever. This makes a network problem more visible, and also gives Kubernetes (or a bash script, or systemd, etc) a chance to restart the application, which might end up resolving this kind of network problem.

### Example usage
This example uses all available configuration parameters, provided as environment variables:

``` bash
export J2G_HOSTNAME=graylog.example.com
export J2G_PORT=12201
export J2G_PACKET_SIZE=1420
export J2G_BLACKLIST="foo.*;bar.*"
sudo journalctl -o json -f | journald2graylog --verbose
```
Or you can simply do:

``` bash
journalctl -o json -f | J2G_HOSTNAME=graylog.example.com journald2graylog
```

And depending on your context, you might actually need to use something more among the line of:

``` bash
sudo journalctl -o json -f | J2G_HOSTNAME=graylog.example.com ./journald2graylog
```

Note that if a parameter specified via an environment variable will override the same parameter specified via command line. For example, in the following command, the _blacklist_ parameter will be set to _localhost_ and not _remotehost_ :

``` bash
J2G_BLACKLIST=localhost journald2graylog --blacklist remotehost
```

## Install

**From source**, you will have to already have a working _go_ development environment setup, with a proper _GOPATH_.

``` bash
go get github.com/cdemers/journald2graylog
```
The resulting binary should be compiled and placed in your GOPATH tree as `$GOPATH/bin/journald2graylog`.

**From binary**, you can download the latest precompiled binary (Linux AMD64) from the [release section](https://github.com/cdemers/journald2graylog/releases).

**Using _make_**, running `make` or `make all` will build a single binary for your current platform.

### Building for Docker
**Using _make_**, you can build a docker image by running `make docker`, it will build the _Linux_ binary and a docker image from _Dockerfile_, and attempt to push image to a docker registry. You must use your own registry by specifying `DOCKER_REGISTRY`, for example:

``` bash
make docker -e DOCKER_REGISTRY=private.registry.org
```
