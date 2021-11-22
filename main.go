package main

import (
	"flag"
	"fmt"
	"os"

	"qumusabel/tools/wifi_mac/client"
)

func usage() {
	fmt.Printf("usage: %s [OPTIONS] ban/unban MAC\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	var (
		ip       = flag.String("i", "192.168.1.1", "Router IP")
		username = flag.String("u", "user", "Router username")
		password = flag.String("p", "user", "Router password")
	)
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
	}

	var (
		action = flag.Arg(0)
		mac    = flag.Arg(1)
	)
	if action != "ban" && action != "unban" {
		flag.Usage()
	}

	client, err := client.Login(*ip, *username, *password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("[+] client logged in")

	switch action {
	case "ban":
		err := client.BanMac(mac)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("banned %s\n", mac)
	case "unban":
		err := client.UnbanMac(mac)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("unbanned %s\n", mac)
	}
}
