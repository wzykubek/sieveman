# Sieveman
Universal ManageSieve protocol client. It can be used as CLI tool or as library as well.

---

## Usage

### Authentication
When using any other command you need to pass credentials.
Password can be passed in similar way like in [mbsync](https://isync.sourceforge.io/mbsync.html), e.g. using GnuPG like below. These flags are global and can be used in any order, before the command or after.

```bash
sieveman -H imap.example.com -u jdoe@example.com -p $(gpg -qd password.asc)
```

Other option is to use something like `dmenu` or `rofi` to show the prompt.

### Listing scripts
To get list of available scripts run following command.
```bash
sieveman ls # ...and login data
```

### Reading script
To download a script use the following command and specify arguments for script name and output file. 
```bash
sieveman get Open-Xchange script.sieve # ...and login data
```

You can use `-` character to print the output to console instead.
```bash
sieveman get Open-Xchange - # ...and login data
```

## License

This project is licensed under ISC license.
