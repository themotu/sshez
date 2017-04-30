package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
	"text/template"
)

func install() error {
	user := getUser()

	rcfile := fmt.Sprintf("%s/.%src", user["home"], user["shell"])
	f, err := os.OpenFile(rcfile, os.O_APPEND|os.O_WRONLY, 0700)
	if err != nil {
		panic(err)
	}

	installText := "alias ssh=\"sshez\"\n"

	defer f.Close()

	_, err = f.WriteString(installText)
	if err != nil {
		return err
	}

	return nil

}

func readLine(fn string, n int) (string, error) {
	if n < 1 {
		return "", fmt.Errorf("invalid request: line %d", n)
	}
	f, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()
	bf := bufio.NewReader(f)
	var line string
	for lnum := 0; lnum < n; lnum++ {
		line, err = bf.ReadString('\n')
		if err == io.EOF {
			switch lnum {
			case 0:
				return "", errors.New("no lines in file")
			case 1:
				return "", errors.New("only 1 line")
			default:
				return "", fmt.Errorf("only %d lines", lnum)
			}
		}
		if err != nil {
			return "", err
		}
	}
	if line == "" {
		return "", fmt.Errorf("line %d empty", n)
	}
	return line, nil
}

func yN(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func runCommand(main string, sub string) {
	subArgs := strings.Split(sub, " ")

	cmd := exec.Command(main, subArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		// fmt.Println(err)
	}
}

func getUser() map[string]string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	configFile := fmt.Sprintf("%s/.ssh/config", dir)

	fullShell := strings.Split(os.Getenv("SHELL"), "/")
	shell := fullShell[len(fullShell)-1]

	return map[string]string{"home": dir, "config": configFile, "shell": shell}
}

func checkAlias(alias string) bool {
	user := getUser()
	if _, err := os.Stat(user["config"]); os.IsNotExist(err) {
		//prompt to make directory and/or file
	}

	alias = strings.ToLower(alias)

	f, err := ioutil.ReadFile(user["config"])
	if err != nil {
		//mkfile here
	}

	s := string(f)
	s = strings.ToLower(s)
	return strings.Contains(s, alias)
}

func checkNick(nick string) bool {
	user := getUser()
	if _, err := os.Stat(user["config"]); os.IsNotExist(err) {
		//prompt to make directory and/or file
	}

	f, err := ioutil.ReadFile(user["config"])
	if err != nil {
		//mkfile here
	}

	regexstr := fmt.Sprintf(`(?m)^(g?)host\s%s$`, nick)

	re, err := regexp.Compile(regexstr)
	if err != nil {
		fmt.Println("There was an issue with that nickname, try again.")
		return true
	}

	match := re.MatchString(string(f))

	if match {
		return true
	}
	return false

}

func grabNick(host string) string {
	var nickname string
	//nicknameQ while true
	fmt.Println("\nAn ssh 'Hostname' AKA nickname is used for helping to remember a host. You will only need to run: 'ssh nickname' and will be able to connect")
	for {
		nick := yN("\nDo you want an ssh nickname?")
		if nick {
			fmt.Print("Enter nickname E.g mycoolazurehost, uberbox, tehcloud, etc. : ")
			fmt.Scanln(&nickname)

			if nickname == "" {
				fmt.Println(fmt.Sprintf("\nApparently you're not with the whole \"nickname program\" so we are using %s", host))
				nickname = host
			}
			if !checkNick(nickname) {
				return nickname
			} else {
				fmt.Println("That nickname is already in use.")
			}
		} else {
			nickname = host
		}
	}
}

type configUpdate struct {
	Host                string
	Hostname            string
	User                string
	Port                uint16
	Identityfile        string
	IdentitiesOnly      string
	TCPKeepAlive        string
	ServerAliveInterval uint32
}

const configTemplate = `
host {{.Host}}
	Hostname {{.Hostname}}
	User {{.User}}
	Port {{.Port}}
	Identityfile {{.Identityfile}}
	IdentitiesOnly {{.IdentitiesOnly}}
	TCPKeepAlive {{.TCPKeepAlive}}
	ServerAliveInterval {{.ServerAliveInterval}}
`

func createConfig(path string, host string, userHome string) {
	keyType := "ed25519"
	//do regex here to see if host is an ipv4 or ipv6
	var user string
	var port uint16
	identityfile := fmt.Sprintf("~/.ssh/%s/%s", host, keyType)
	IdentitiesOnly := "yes"
	TCPKeepAlive := "yes"
	ServerAliveInterval := 120

	port = 22

	nickname := grabNick(host)

	fmt.Print("\nWhat user will you ssh with?: ")
	fmt.Scanln(&user)

	fmt.Print("\nWhat Port?[22]: ")
	fmt.Scanln(&port)

	HostToAdd := configUpdate{nickname, host, user, uint16(port), identityfile, IdentitiesOnly, TCPKeepAlive, uint32(ServerAliveInterval)}
	t := template.Must(template.New("config").Parse(configTemplate))

	configLocation := fmt.Sprintf("%s/.ssh/config", userHome)

	f, err := os.OpenFile(configLocation, os.O_APPEND|os.O_WRONLY, 0700)
	if err != nil {
		panic(err)
	}
	t.Execute(f, HostToAdd)
	f.Close()
	fmt.Println("\nUpdated config file with the following:")
	t.Execute(os.Stdout, HostToAdd)
}

func createHost(host string) {
	user := getUser()
	keyPath := fmt.Sprintf("%s/.ssh/%s", user["home"], host)
	fmt.Println(fmt.Sprintf("created folder %s", keyPath))
	os.MkdirAll(keyPath, 0700)

	//todo: get from config type
	keyType := "ed25519"

	params := fmt.Sprintf("-t %s -f %s/%s", keyType, keyPath, keyType)

	runCommand("ssh-keygen", params)
	//update host here with nickname if desired
	fmt.Println(fmt.Sprintf("%s added", host))
	createConfig(keyPath, host, user["home"])
	copyKey(host)
	os.Exit(0)
}

func updateHost() {
	fmt.Println("I can't do this yet :/")
}

func getPubkey(host string) (string, string) {

	var pubkey string

	user := getUser()

	//todo: get from config type
	keyType := "ed25519"

	fileName := fmt.Sprintf("%s/.ssh/%s/%s.pub", user["home"], host, keyType)

	_, err := os.Stat(fileName)
	if err != nil {
		fmt.Println("That host doesn't exist")
	} else {
		pubkey, _ = readLine(fileName, 1)
	}
	return host, pubkey
}

func copyKey(host string) {
	keyType := "ed25519"
	params := fmt.Sprintf("cat ~/.ssh/%s/%s.pub | ssh %s 'mkdir ~/.ssh; cat >> ~/.ssh/authorized_keys' 2>/dev/null \n", host, keyType, host)
	fmt.Print("Run the following command and your key will be copied to the remote server. I suggest you then disable password login.\n\n")
	fmt.Println(params)
	fmt.Println(fmt.Sprintf("After running this command you will be able to ssh without a password to %s", host))
	os.Exit(0)
}

func main() {
	var fA string

	if len(os.Args) > 1 {
		fA = os.Args[1]
	} else {
		// if they're just typing ssh let them, poor souls...
		cmd := exec.Command("ssh")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.Run()
		os.Exit(0)
	}

	if len(os.Args) == 2 && fA == "install" {
		install()
		os.Exit(0)
	}

	if len(os.Args) == 3 { // ssh update hostname == 3
		switch fA {
		case "update":
			updateHost()
			os.Exit(0)
		case "pub":
			host, key := getPubkey(os.Args[2])
			fmt.Println(fmt.Sprintf("\nPublic key for %s:\n\n%s", host, key))
			os.Exit(0)
		case "copy":
			copyKey(os.Args[2])
			os.Exit(0)
		case "nick":
			checkNick(os.Args[2])
			os.Exit(0)
		}
	}

	if checkAlias(fA) || checkNick(fA) {
		// os.Exit(0)
		if checkAlias(fA) && !checkNick(fA) {
			// check if we are copying the key for the first time and let it do it's thing
			if strings.Contains(strings.Join(os.Args[2:], ","), `mkdir ~/.ssh; cat >> ~/.ssh/authorized_keys`) {
				command := strings.Join(os.Args[1:], " ")
				runCommand("ssh", command)
				os.Exit(0)
			}
			fmt.Println("That Hostname is in your config but using a different alias AKA host.")
			fmt.Println("This will be added in 1.1 meanwhile you should modify the config file or remember the alias you used.")
			os.Exit(0)
		}
		command := strings.Join(os.Args[1:], " ")
		runCommand("ssh", command)
	} else {
		createHost(fA)
		copyKey(fA)
	}

}
