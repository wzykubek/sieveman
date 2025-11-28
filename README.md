# Sieveman
Universal ManageSieve protocol client. It can be used as CLI tool or library.
It works in two modes: script and interactive.

[![Go Reference](https://pkg.go.dev/badge/go.wzykubek.xyz/sieveman.svg)](https://pkg.go.dev/go.wzykubek.xyz/sieveman)

---

## Usage
The command naming convention is taken from common *nix tools, so it should be intuitive. However, there is some breakdown of the commands with examples below.

If you run `sieveman` without any arguments, it will start interactive mode. Keep in mind that you need to pass at least credentials as shown below.

### Authentication
When using any command you need to pass credentials. You can create a recursive alias for the `sieveman` command with the following options.
```bash
sieveman --host imap.example.com --username jdoe@example.com --password $(gpg -qd encrypted_password.txt.gpg)
```

### Listing
```bash
sieveman ls
```
It prints a list of all scripts with `*` (star) indicator for the active one.
```
Open-Xchange*
test_script.sieve
```
Run the command with `--help` flag to see all available options.

### Downloading
```bash
sieveman get Open-Xchange script.sieve
```
It will download Open-Xchange script and save it to local file named script.sieve. If file with the given filename exists, you can use `--force` flag to overwrite it.

You can use `-` (minus) character instead of the file name to print the output to console instead.
```bash
sieveman get Open-Xchange -
```

### Uploading
```bash
sieveman put script.sieve Open-Xchange
```
It will put local script.sieve as Open-Xchange. If script with following name exists on the server it **will be overwritten**.
Run the command with `--help` flag to see all available options.

### Direct editing
```bash
sieveman edit Open-Xchange
```
It will open Open-Xchange script in the `$EDITOR` and upload it to the server on saving.

### Activating and deactivating
```bash
sieveman activate Open-Xchange
```
It will activate Open-Xchange script on the server. Keep in mind that only one script can be active at a time.

```bash
sieveman deactivate
```
It will deactivate **all** scripts.

### Removing and renaming
```bash
sieveman rm Open-Xchange
```
It will remove Open-Xchange script from server.

```bash
sieveman mv Open-Xchange test.sieve
```
It will rename Open-Xchange script to test.sieve without changing active status.

## Installation
### Build from source
This is universal method to build a binary for any system and architecture. You need to have Go installed.
```bash
git clone https://github.com/wzykubek/sieveman && cd sieveman
```
```bash
make
```

On Linux distributions you can also install it to specified PREFIX (default `/usr/local`) with ease.
```bash
sudo make install
```

### Install from repository
 + [Arch Linux (AUR)](https://aur.archlinux.org/packages/sieveman)

## License
This project is licensed under [ISC license](LICENSE).
