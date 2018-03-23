# Contra

AJA's network device configuration tracker software built using GoLang.
Initially focused on routers and switches, but expected to track all
network device configurations.

### Features

- Automated configuration watching for all supported network devices.
- Tracks changes in GIT.
- Review status and download config files via web portal.

## Installation

```
git clone gitlab.aja.com/go/contra
dep ensure
make run
```

Alternatively, simply run `make` then execute `./bin/contra`

## Dependency Management

Considering `dep` for tracking dependencies.
See https://golang.github.io/dep/docs/installation.html

## Supported Devices

### Current

- None

### Soon

- PFSense
- Cisco CSB

### Someday

- Cisco
- Comware
- Vyatta
- Juniper (JunOS)

## Device Configuration

If a password, or any field in the config file, contains a `#` or `;` character be sure to properly
quote the password with either a backtick ``` ` ``` or a set of three double-quotes ``` """ ``` for
example, if your password is `Some#pass;word` you will need one of the following formats:

    Pass=`Some#pass;word`
    Pass="""Some#pass;word"""

## License

Contra is licensed under the MIT License

## Acknowledgements

- This project was inspired by Splendid, Sweet, and Rancid.

