# sshez
A simple ssh config and key automation tool.

This isn't by any means a replacement for ssh. It still uses the standard ssh client. It simply checks to see if an alias exists. If it doesn't it creates it and generates you a key.


## Usage

With domain:
```
ssh mysite.mydomain.com

It looks like mysite.mydomain.com isn't in your config file. Would you like to add it now and generate a new key for it? (Y/n): y

Passphrase? (Y/n): n

A new key has been generated at ~/.ssh/mysite.mydomain.com/id_rsa
```

With IP:
```
ssh 172.17.0.1

It looks like 172.17.0.1 isn't in your config file. Would you like to add it now and generate a new key for it? (Y/n/advanced): y

Passphrase? (Y/n): n

What alias would you like to use for 172.17.0.1?: myserver

A new pub/private key has been generated at ~/.ssh/myserver/id_rsa
```

```
sshez pub mysite.mydomain.com
Public key:
<public key>
```

## Installation
Install the package
`go get github.com/themotu/sshez`

Install it for your user:
`sshez install ~/.bashrc` or for zsh `sshez install ~/.zshrc`

## Settings
TODO: key types, folder name scheme, collision detecion for alias, check for alias in install, check for sshez in alias, prompt for config file creation and ssh folder, promt for remove password auth on server

## Technical details

