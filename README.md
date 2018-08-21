# Contra

AJA's network device configuration tracker software built using GoLang.
Initially focused on routers and switches, but expected to track all
network device configurations.

### Features

* Automated configuration watching for all supported network devices.
* Tracks changes in GIT.
* Review status and download config files via web portal.

## Installation

```
git clone gitlab.aja.com/go/contra
make first
```

After the initial build, `make run` or `make` then execute `./bin/contra` will work.

## Dependency Management

* `dep` is used for dependency tracking
* See https://golang.github.io/dep/docs/installation.html

## Supported Devices

### Current

* pfSense
* Cisco Small Business
* Vyataa
* HP/Procurve
* HP/Comware
    * For locked down devices use UnlockPass in the device configuration to unlock xtd-cli-mode

### Soon

### Someday

* Cisco
* Juniper (JunOS)

## Device Configuration

If a password, or any field in the config file, contains a `#` or `;` character be sure to properly
quote the password with either a backtick ``` ` ``` or a set of three double-quotes ``` """ ``` for
example, if your password is `Some#pass;word` you will need one of the following formats:

```
Pass=`Some#pass;word`
Pass="""Some#pass;word"""
```

## Building

Requires `fpm` found at: http://fpm.readthedocs.io/en/latest/installing.html

## License

Contra is licensed under the MIT License

## Acknowledgements

* This project was inspired by Splendid, Sweet, and Rancid.

