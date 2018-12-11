# Contra

AJA's network device configuration tracker software built using GoLang.
Initially focused on routers and switches, but expected to track all
network device configurations.

As the MIT license states, use this code at your own risk. It has been
running smoothly for months on end without any issues, but something could
always go wrong.

### Features

* Automated configuration tracking for all supported network devices.
* Tracks changes in GIT, plus email notifications for changes.
* Email notifications when a device is unable to be backed up.

## Installation

```
git clone gitlab.aja.com/go/contra
make first
```

After the initial build, `make run` or `make` then execute `./bin/contra` will work.

### Configuration

Copy the `contra.example.conf` file to `contra.conf` and configure.

#### Encrypt Passwords

Note, all passwords are encrypted using an EncryptKey on the first run.
If an EncryptKey is not set, one will be randomly generated.
While the encryption key is in the conf file, this makes it so the passwords are not stored as plain
text.

We do this to prevent inadvertent password leak from someone standing behind you.
As well as to make it a little harder for someone who stumbles upon the file to determine the passwords.
Of course, with the source code, they can easily determine how to decrypt the passwords.

Set `EncryptPasswords = false` to disable this behavior.

#### Git Push

By default, changes will be committed to the defined workspace folder.

You can set up the workspace folder's git repo with a remote origin, and set GitPush to true.

By default, GitAuth is also set to true, and will use `.ssh/id_rsa` private keypair to attempt the push.
You can set GitAuth to false and it is effectively a `git push` from the workspace folder.

## Dependency Management

* `dep` is used for dependency tracking
* See https://golang.github.io/dep/docs/installation.html

## Supported Devices

### Current

* pfSense
* Cisco Small Business
* Vyatta
* HP/Procurve
* HP/Comware
    * For locked down devices use UnlockPass in the device configuration to unlock xtd-cli-mode

### Soon

* Cisco
* Juniper (JunOS)
* Arista

### Someday

* http
* custom scripting

## Device Configuration

If a password, or any field in the config file, contains a `#` or `;` character be sure to properly
quote the password with either a backtick ``` ` ``` or a set of three double-quotes ``` """ ``` for
example, if your password is `Some#pass;word` you will need one of the following formats:

```
Pass=`Some#pass;word`
Pass="""Some#pass;word"""
```

## Building
* Binary only: `make binaries`
* Linux packages (.deb and .rpm) `make packages`
  * requires `fpm` - http://fpm.readthedocs.io/en/latest/installing.html
* Release build (compressed binary and packages) `make release`
  * requires `fpm` - http://fpm.readthedocs.io/en/latest/installing.html
  * requires `upx` - https://upx.github.io/
* If you would like to build for another platform `GOOS=$platform GOARCH=$arch go build contra.go`
  * While Contra may work on platforms other than Linux it is untested.

## Contributing

If you have a device that is not yet supported, please review on of the existing devices and
consider making a pull request with support for your device. Alternatively, you can create an issue
and provide logs showing the steps necessary to pull the configs, and we can try to add support
for the device.

Be sure to `make test` to run fmt, vet, lint, and tests before each commit.

## License

Contra is licensed under the MIT License

## Acknowledgements

* This project was inspired by Splendid, Sweet, and Rancid.

