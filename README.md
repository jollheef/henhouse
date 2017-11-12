[![Build Status](https://travis-ci.org/jollheef/henhouse.svg?branch=master)](https://travis-ci.org/jollheef/henhouse)
[![GoDoc](https://godoc.org/github.com/jollheef/henhouse?status.svg)](http://godoc.org/github.com/jollheef/henhouse)
[![Coverage Status](https://coveralls.io/repos/jollheef/henhouse/badge.svg?branch=master&service=github)](https://coveralls.io/github/jollheef/henhouse?branch=master)
[![Go Report Card](http://goreportcard.com/badge/jollheef/henhouse)](http://goreportcard.com/report/jollheef/henhouse)

# Henhouse

Scoreboard for jeopardy-style CTFs.

Fundamental principle: if henhouse is not helping you make jeopardy-style CTF easily, then there is a bug in henhouse.

![Imgur image](https://i.imgur.com/uMCFPw7.png)

## Install

### Packagecloud
[![packagecloud](https://packagecloud.io/assets/packagecloud-badge-fbea7fd09f5aab38e8d59fec16f2268c.png)](https://packagecloud.io/jollheef/henhouse)

Built for Ubuntu 16.04.

    $ curl -s https://packagecloud.io/install/repositories/jollheef/henhouse/script.deb.sh | sudo bash
    $ sudo apt install henhouse

### Build deb package from source

    $ apt install golang build-essential binutils upx-ucl
    $ export GOPATH=$(realpath ./) && go get github.com/jollheef/henhouse/...
    $ cd ${GOPATH}/src/github.com/jollheef/henhouse
    $ ./package.sh

## Development

### Depends

#### Gentoo

    $ sudo emerge dev-lang/go dev-db/postgresql

#### Ubuntu

    $ sudo apt install golang postgresql

### Build

First you need set GOPATH environment variable.

    $ export GOPATH=$(realpath ./)

After you need download and build henhouse with depends.

    $ go get github.com/jollheef/henhouse

### Run

    $ sudo psql -U postgres
    postgres=# CREATE DATABASE henhouse;
    postgres=# CREATE USER henhouse WITH password 'STRENGTH_PASSWORD';
    postgres=# GRANT ALL privileges ON DATABASE henhouse TO henhouse;

After that you need to fix 'connection' parameter in configuration file.
(And other parameters, of course)

Now, run it!

    $ ${GOPATH}/bin/henhouse ${GOPATH}/src/github.com/jollheef/henhouse/config/henhouse.toml --reinit
