# cider

`cider` is a simple library/cli tool for reporting information about CIDR blocks.

## Overview

The tool expects a CIDR block as input and behaves as follows:

 - With no additional arguments, display information about the network block
 - With `-s SIZE`, display all CIDR blocks of `SIZE` than can fit within the network
 - With `-a`, list all CIDR blocks which can fit within the network

## Installation

Install `cider` to your $GOPATH using `make install` or build it using `make build`.

## Tests

Test with `make test`

## Example usage

Display information about a CIDR block:

```
$ cider 192.168.0.0/24
Network:    192.168.0.0
First:      192.168.0.1
Last:       192.168.0.254
Broadcast:  192.168.0.255
Netmask:    255.255.255.0
Usable:     254
Next:       192.168.1.0
```

Show the /26 networks that can fit within CIDR block 192.168.0.0/24:

```
$ cider -s 26 192.168.0.0/24
192.168.0.0/26
192.168.0.64/26
192.168.0.128/26
192.168.0.192/26
```

Show all networks that can fit within CIDR block 192.168.0.0/27:
```
$ cider -a 192.168.0.0/27
/28:
192.168.0.0/28
192.168.0.16/28

/29:
192.168.0.0/28
192.168.0.8/28
192.168.0.16/28
192.168.0.24/28
```
