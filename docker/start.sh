#!/bin/sh

/usr/bin/journalctl -m -f -o json | journald2graylog
