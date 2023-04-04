# rpi-alpine-installer

alpine-linux installer(for just 1 minute!)

## Tool Install

```
go get github.com/SWITCHSCIENCE/rpi-alpine-installer
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

## First boot & SSH login

finalize script running message:

```
* Starting firstboot ... [ ok ]
* Starting local ...     [ ok ]

Welcome to Alpine Linux 3.13
```

```
ssh root@raspberrypi.local
```

## Backup to GitHub and Restore from GitHub

Backup:

1. microSD card backup into local `dist/` folder.
2. `dist/` folder initialize for git: `git init`.
3. `git add .` & `got commit -am "add files"`.
4. if you need git remote add: `git remote add origin git@github.com:...`.
5. push to gitHub main branch: `git push -u origin main`.

Restore:

1. microSD card format for FAT32 with label 'ALPINE'.
2. `git clone ...`
3. Change dir. into repos: `cd ...`
4. export file into microSD card volume: `git archive main | tar xv -C /Volumes/ALPINE`
