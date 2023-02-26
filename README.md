# Amalgam CLI
Small CLI Download Tool written in Go to Download Episodes from amalgam-fansubs.moe or Music from detektiv-conan.ch.
Rewritten from the Original Repository: [Amalgam Detektiv Conan Downloader](https://gitlab.com/mauamy/amalgamdetektivconandownloader).

# Installation
## Download releases
Download the binary for your system from the [release page](https://github.com/kobayashi90/amalgam-cli/releases).

## Build it yourself
If you want to build it yourself, simply clone this repository and use the makefile
for building and installing it.
```bash
# local/test linux build
$ make build

# install it to your GOPATH (linux)
$ make install

# uninstall from your GOPATH (linux)
$ make uninstall

# build for windows
$ make windows

# build for mac
$ make mac
```

# Usage
#### General Usage
```bash
$ adcl <subcommand> <action> [flags] <parameter>
```
**Subcommands**: 
- episodes
- music

**Actions**:
- download
- list
 

#### Show Help
```bash
$ adcl -h
```

### Episodes
#### List Episodes
```bash
$ adcl episodes list
$ adcl episodes l
```
##### Flags
- **--dlink, -d**: show the default download link
- **--gdrive, -g**: show an indicator if a google drive download link is available
- **--format \<value\>, -f \<value\>**: set the output format. Available values: **csv, html, md** 

#### Download Episodes
```bash
$ adcl episodes download <episode_numbers>
$ adcl episodes d <episode_numbers>
```
*<episode_numbers>* is a separated list of episode numbers: **1 2 3 4**.
In addition you can provide ranges within this list: **1 2 3-8 10**. 
```bash
# example
$ adcl episodes d 710 840-845 870
```
##### Flags
- **--gdrive, -g**: download the episode from google drive


### Music
#### List Music
```bash
$ adcl music list
$ adcl music l
```
##### Flags
- **--format \<value\>, -f \<value\>**: set the output format. Available values: **csv, html, md** 

#### Download Music
```bash
$ adcl music download <music_ids>
$ adcl music d <music_ids>
```
*<music_ids>* is a separated list of music IDs matching the IDs in the list: **1 2 3 4**.
In addition you can provide ranges within this list: **1 2 3-8 10**. 
```bash
# example
$ adcl music d 1 3-6 22
```
##### Flags
- **--unzip, -u**: extract the zip archive after download
- **--keepArchive, -k**: do not remove the archive after extraction 
