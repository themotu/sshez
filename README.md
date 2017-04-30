# sshez
A simple ssh config and key automation tool for individual server keys and passwordless ssh.

This isn't by any means a replacement for ssh. It still uses the standard ssh client. It simply checks to see if an alias and key exist in your config. If it doesn't it creates it and generates you a key and sets up the alias for you.

The benefit is you have a clean concise way of generating a new key per server. I find it's far too easy to just re-use a key and this allows me to have a new key per server with very little effort.


## Usage

### Setting up a domain
(after running install, otherwise you'll have to use `sshez` instead of `ssh`)
```
ssh mysite.mydomain.com

created folder /home/myuser/.ssh/mysite.mydomain.com
Generating public/private ed25519 key pair.
Enter passphrase (empty for no passphrase): 
Enter same passphrase again: 
Your identification has been saved in /home/myuser/.ssh/mysite.mydomain.com/ed25519.
Your public key has been saved in /home/myuser/.ssh/mysite.mydomain.com/ed25519.pub.
The key fingerprint is:
SHA256:Tl+Qn2yO4thplSZb51S02Q34s+W1+5SOOOfREinRVps myuser@blue-steel
The key's randomart image is:
+--[ED25519 256]--+
|             . . |
|           .o + o|
|          o. = E.|
|           +o.B =|
|        S  o*+ =o|
|       o..==+ +.o|
|        o*o+.o oo|
|       +oo .o.=o |
|      ..+  .+o .o|
+----[SHA256]-----+
mysite.mydomain.com added

An ssh 'Hostname' AKA nickname is used for helping to remember a host. You will only need to run: 'ssh nickname' and will be able to connect

Do you want an ssh nickname? [y/n]: 

What user will you ssh with?: myuser

What Port?[22]: 

Updated config file with the following:

host mysite.mydomain.com
	Hostname mysite.mydomain.com
	User myuser
	Port 22
	Identityfile ~/.ssh/mysite.mydomain.com/ed25519
	IdentitiesOnly yes
	TCPKeepAlive yes
	ServerAliveInterval 120
Run the following command and your key will be copied to the remote server. I suggest you then disable password login.

cat ~/.ssh/mysite.mydomain.com/ed25519.pub | ssh mysite.mydomain.com 'mkdir ~/.ssh; cat >> ~/.ssh/authorized_keys' 2>/dev/null 

After running this command you will be able to ssh without a password to mysite.mydomain.com
```
### View public key:
```
ssh pub mysite.mydomain.com

Public key for mysite.mydomain.com:

ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIE30fImlS0BliJ1xyAydhbXsPbG7vxpMNLh2g0to0FxW myuser@blue-steel
```

### View copy instruction:
Use this to view key instructions for passwordless ssh to host
```
ssh copy mysite.mydomain.com
Run the following command and your key will be copied to the remote server. I suggest you then disable password login.

cat ~/.ssh/mysite.mydomain.com/ed25519.pub | ssh mysite.mydomain.com 'mkdir ~/.ssh; cat >> ~/.ssh/authorized_keys' 2>/dev/null 

After running this command you will be able to ssh without a password to mysite.mydomain.com
```

## Installation

For 64-bit linux:

`sudo bash -c "curl -L https://github.com/themotu/sshez/releases/download/latest/sshez-amd64 -o /usr/local/bin/sshez" && sudo chmod +x /usr/local/bin/sshez`

Other OSes:

change the above sshez-amd64 to any of:

* `sshez-386` for 32-bit linux
* `sshez-arm` for raspberry pi or other arm devices
* `sshez-osx` for Mac OSX

Install it for your user (bash or zsh):
`sshez install`

The above will require you to log out and log back in our to "source" your rc file until you log out: `source .bashrc` or `source .zshrc`

## Settings
TODO: 
* key types (config)
* folder name scheme(config)
* check for alias in install
* add version
* add help
* prompt for config file creation and ssh folder

Post 1.0:
* actually copy key to server instead of giving user command
* prompt for remove password auth on server
* additional services (github accounts, aws api, etc.)

## Technical details

By default we are using the ed25519 key type introduced in ssh 6.5. You can change this in the settings file. Elliptic-curve Diffie Hellman in Daniel Bernstein's Curve25519 offers better security than ECDSA and DSA as well as good performance. You may change this in the settings file in .ssh/sshez.conf