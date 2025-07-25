# Sieveman
Universal ManageSieve protocol client. It can be used as CLI tool or as library as well. Currently works in script mode, but interactive mode will be available in future.

---

## Usage
Command naming convention is taken from common *nix tools, so you will probably feel like in home. However there is some documentation with examples below.

### Authentication
When using any command you need to pass credentials. You can create recursive alias to `sieveman` command with following options.
```bash
sieveman --host imap.example.com --username jdoe@example.com --password $(gpg -qd encrypted_password.txt.gpg)
```

### Listing scripts
```bash
sieveman ls
```
It prints list of all scripts with '*' indicator for active one.
```
Open-Xchange*
test_script.sieve
```
Run command with `--help` flag to see all available options.

### Reading script content
```bash
sieveman get Open-Xchange script.sieve
```
It will download Open-Xchange script and save it to local file script.sieve. If file with given filename exists, you can use `--force` flag.

You can use '-' (minus) character instead of file name to print the output to console instead.
```bash
sieveman get Open-Xchange -
```

### Uploading script
```bash
sieveman put script.sieve Open-Xchange
```
It will upload script to Open-Xchange server. If script with following name exists on server it **will be overwritten**.
Run command with `--help` flag to see all available options.

### Directly editing script
```bash
sieveman edit Open-Xchange
```
It will open script in default editor and upload it to server after saving.

To change used editor run with `$EDITOR` environment variable.
```bash
EDITOR=nano sieveman edit Open-Xchange
```

### Activating and deactivating scripts
```bash
sieveman activate Open-Xchange
```
It will activate script on server. Keep in mind that only one script can be active at a time.

```bash
sieveman deactivate
```
It will deactivate **all** scripts.

### Removing and renaming scripts
```bash
sieveman rm Open-Xchange
```
It will remove script from server.

```bash
sieveman mv Open-Xchange test.sieve
```
It will rename script on server without changing active status.

## Installation
### From Source
This is universal method to build a binary for any system and architecture. You need to have Go installed.
```bash
git clone https://github.com/wzykubek/licensmith
cd licensmith
# This command will create dist directory with compiled binary and generated shell completions.
make
```

On Linux distributions you can also install it to `/usr/local` with ease.
```bash
sudo make install
```

### Using Go Module Proxy
You can install this package using `go install`. Ensure `$GOPATH/bin` is in your `$PATH`.
```bash
go install go.wzykubek.xyz/sieveman@latest
```

## License
This project is licensed under ISC license.
