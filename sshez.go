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
	var nickname string
	var user string
	var port uint16
	identityfile := fmt.Sprintf("~/.ssh/%s/%s", host, keyType)
	IdentitiesOnly := "yes"
	TCPKeepAlive := "yes"
	ServerAliveInterval := 120

	port = 22
	//nicknameQ while true
	fmt.Println("\nAn ssh 'Hostname' AKA nickname is used for helping to remember a host. You will only need to run: 'ssh nickname' and will be able to connect")
	nick := yN("\nDo you want an ssh nickname?")

	if nick {
		fmt.Print("Enter nickname E.g mycoolazurehost, uberbox, tehcloud, etc. : ")
		fmt.Scanln(&nickname)

		if nickname == "" {
			fmt.Println(fmt.Sprintf("\nApparently you're not with the whole \"nickname program\" so we are using %s", host))
			nickname = host
		}
	} else {
		nickname = host
	}

	fmt.Print("\nWhat user will you ssh with?: ")
	fmt.Scanln(&user)

	fmt.Println("\nWhat Port?[22]: ")
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
	fmt.Println(fmt.Sprintf("After running this command you will be able to ssh without a key to %s", host))
	os.Exit(0)
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "install" {
		install()
		os.Exit(0)
	}

	if len(os.Args) == 3 { // ssh update hostname == 3
		switch os.Args[1] {
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
		}
	}

	if checkAlias(os.Args[1]) {
		command := strings.Join(os.Args[1:], " ")
		runCommand("ssh", command)
		// fmt.Println("\nUmm... Dude? That host is in there...")
		// os.Exit(0)
	} else {
		createHost(os.Args[1])
		copyKey(os.Args[1])
	}

}
