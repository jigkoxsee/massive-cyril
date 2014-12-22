package main

import (
	//"encoding/json"
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	//"strings"
)

/*
func tunnelOutput(stdout io.ReadCloser,) {
	status := false
	r := bufio.NewReader(stdout)
	line, _, err := r.ReadLine()
	for ;  ; {
		line, _, _ = r.ReadLine()
		status = true
	}
}
*/

func argConfig() (string, string, string) {
	flag.Parse()
	username, password, port := flag.Arg(0), flag.Arg(1), flag.Arg(2)

	reader := bufio.NewReader(os.Stdin)
	if username == "" {
		fmt.Println("Username : ")
		username, _ = reader.ReadString('\n')
		username = username[:len(username)-2]
	}
	if password == "" {
		fmt.Println("Password : ")
		password, _ = reader.ReadString('\n')
		password = password[:len(password)-2]
	}
	if port == "" {
		port = "443"
	}

	return username, password, port
}

/**
app.exe username password port
*/

func main() {

	var config struct {
		username string
		password string
		port     string
	}
	fmt.Println(flag.Arg(0))
	config.username, config.password, config.port = argConfig()

	remote := config.username + "@vps.jigko.me"
	fmt.Println("----------------")

	// SSH Tunnel
	tunnelExec := exec.Command("plink.exe", "-ssh", remote, "-pw", config.password, "-C", "-T", "-L", "127.0.0.1:"+config.port+":127.0.0.1:"+config.port)
	tunnelExec.Stdout = os.Stdout
	//tunnelOut, _ := tunnelExec.StdoutPipe()
	tunnelIn, _ := tunnelExec.StdinPipe()

	// OpenVPN
	openvpnExec := exec.Command("openvpn/openvpn.exe", "openvpn/config/client.ovpn")
	openvpnExec.Stdout = os.Stdout
	//openvpnOut, _ := openvpnExec.StdoutPipe()
	openvpnIn, _ := openvpnExec.StdinPipe()

	if err := tunnelExec.Start(); err != nil {
		log.Fatal(err)
	}

	if err := openvpnExec.Start(); err != nil {
		log.Fatal(err)
	}

	if err := tunnelExec.Wait(); err != nil {
		log.Fatal(err)
	}

	io.WriteString(tunnelIn, "tunnel\n")
	io.WriteString(openvpnIn, "openvpn\n")

	fmt.Printf("END")

}
