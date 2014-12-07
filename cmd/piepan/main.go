package main

import (
	"crypto/tls"
	"flag"
	"os"
	"fmt"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleutil"
	"github.com/layeh/piepan"
)

func main() {
	// Flags
	username := flag.String("username", "piepan-bot", "username of the bot")
	password := flag.String("password", "", "user password")
	server := flag.String("server", "localhost:64738", "address of the server")
	certificateFile := flag.String("certificate", "", "user certificate file (PEM)")
	keyFile := flag.String("key", "", "user certificate key file (PEM)")
	insecure := flag.Bool("insecure", false, "skip certificate checking")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] [scripts...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "a bot framework for Mumble\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Configuration
	mux := gumbleutil.EventMultiplexer{}
	keepAlive := make(chan bool)

	config := gumble.Config{
		Username: *username,
		Password: *password,
		Address:  *server,
		Listener: &mux,
	}

	if *insecure {
		config.TLSConfig.InsecureSkipVerify = true
	}
	if *certificateFile != "" {
		if *keyFile == "" {
			keyFile = certificateFile
		}
		if certificate, err := tls.LoadX509KeyPair(*certificateFile, *keyFile); err != nil {
			panic(err)
		} else {
			config.TLSConfig.Certificates = append(config.TLSConfig.Certificates, certificate)
		}
	}

	client := gumble.NewClient(&config)

	// piepan
	piepan := piepan.New(client)
	defer piepan.Destroy()

	for _, script := range flag.Args() {
		if err := piepan.LoadScriptFile(script); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}

	// Event multiplexer
	mux.Attach(piepan)
	mux.Attach(gumbleutil.Listener{
		Disconnect: func(e *gumble.DisconnectEvent) {
			keepAlive <- true
		},
	})

	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	<-keepAlive
}
