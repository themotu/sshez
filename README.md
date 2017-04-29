# sshez
A simple ssh config and key automation tool.

This isn't by any means a replacement for ssh. It still uses the standard ssh client. It simply checks to see if an alias and key exist in your config. If it doesn't it creates it and generates you a key and sets up the alias for you.

The benefit is you have a clean concise way os generating a new key per server. I find it's far too easy to just re-use a key and this allows me to have a new key per server with very little effort.


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
sudo bash -c "curl -L https://github.com/themotu/sshez/releases/download/latest/sshez-amd64 -o /usr/local/bin/sshez" && sudo chmod +x /usr/local/bin/sshez

Install it for your user (bash or zsh):
`sshez install`

The above will require you to log out and log back in our to "source" your rc file until you log out: `source .bashrc` or `source .zshrc`

## Settings
TODO: key types, folder name scheme, collision detecion for alias, check for alias in install, check for sshez in alias, prompt for config file creation and ssh folder, promt for remove password auth on server, bug with aliases

## Technical details

By default we are using the ed25519 key type introduced in ssh 6.5. You can change this in the settings file. Elliptic-curve Diffie Hellman in Daniel Bernstein's Curve25519offers better security than ECDSA and DSA as well as good performance. You may change this in the settings file in .ssh/sshez.conf