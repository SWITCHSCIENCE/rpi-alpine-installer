# rpi-alpine-installer

## Tool Install

```
go get github.com/144lab/rpi-alpine-installer
```

## Format micro-SD card for FAT32

for macOS

```
diskutil eraseDisk FAT32 ALPINE MBRFormat /dev/diskN
```

## Install alpine linux into micro-SD card

for aarch64(RaspberryPi 3 or 4)

```shell
curl https://github.com/<GitHub-UserID>.keys > keys
rpi-alpine-installer -version=v3.13.4 -arch=aarch64 \
	-ssid=<SSID> -passphrase=<Passphrase> \
	-authorized_keys=keys \
	-dist=/Volumes/ALPINE
```

for armhf(RaspberryPi 3 or 4)

```shell
curl https://github.com/<GitHub-UserID>.keys > keys
rpi-alpine-installer -version=v3.13.4 -arch=armhf \
	-ssid=<SSID> -passphrase=<Passphrase> \
	-authorized_keys=keys \
	-dist=/Volumes/ALPINE
```
