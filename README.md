# journald2graylog
Command line tool to forward systemd's journald's logs to a Graylog server taking advantage of the descriptive GELF format.

The journald2graylog command expects it's paramters to be provided as the environment variables.

## Example usage
```
export J2G_HOSTNAME=graylog.example.com
export J2G_PORT=12201
export J2G_PACKET_SIZE=1420
journalctl -o json -f | journald2graylog 
```
Or simply:
```
journalctl -o json -f | J2G_HOSTNAME=graylog.example.com journald2graylog 
```
And depending on your context, you might actually need to use something more among the line of:
```
sudo journalctl -o json -f | J2G_HOSTNAME=graylog.example.com ./journald2graylog
```

## Install

From source, you will have to already have a working _go_ development environment setup, with a proper _GOPATH_.
```
go get github.com/robertkowalski/graylog-golang
go get github.com/cdemers/journald2graylog
```

From binary, you can download the latest precompiled binary (Linux AMD64) from the [release section](https://github.com/cdemers/journald2graylog/releases).

