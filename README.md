# hackattic

My solutions to [hackattic challenges](https://hackattic.com/changelog).

Here is my profile - https://hackattic.com/u/PDmatrix

## Overview

All solutions are in the `pkg` folder. Each solution has the same folder name as the hackattic challenge name.

## Usage

You can build the binary locally and run it

```bash
go build -o hackattic ./cmd/hackattic
HACKATTIC_ACCESS_TOKEN=[secret_access_token] ./hackattic --challenge [challenge_name]
```

Or you can open this repository in `Visual Studio Code` and use `Run and Debug`. But you need to create an `.env` file with `HACKATTIC_ACCESS_TOKEN=secret` in it.

## About the challenges

Some challenges require you to have `docker` installed on your system. There are also some challenges that use other external tools such as `backup_restore` which requires you to have `psql` installed. You can inspect the code to check if external tools are required.

This also means that these challenges won't work on `Windows OS`.

