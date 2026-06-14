---
title: "Installation"
description: "Install regex101 from a release, with go install, or from source."
weight: 20
---

## Prebuilt binaries

Every [release](https://github.com/tamnd/regex101-cli/releases) carries archives for Linux, macOS,
and Windows on amd64 and arm64, plus deb, rpm, and apk packages for Linux.
Download, unpack, put `regex101` on your `PATH`, done. The `checksums.txt`
on each release is signed with keyless [cosign](https://docs.sigstore.dev/) if
you want to verify before running.

## With Go

```bash
go install github.com/tamnd/regex101-cli/cmd/regex101@latest
```

That puts `regex101` in `$(go env GOPATH)/bin`, which is `~/go/bin` unless
you moved it. Make sure that directory is on your `PATH`.

## From source

```bash
git clone https://github.com/tamnd/regex101-cli
cd regex101-cli
make build        # produces ./bin/regex101
./bin/regex101 version
```

## Container image

```bash
docker run --rm ghcr.io/tamnd/regex101:latest --help
```

## Checking the install

```bash
regex101 version
```

prints the version and exits.
