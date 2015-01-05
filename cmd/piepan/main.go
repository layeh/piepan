package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleutil"
	"github.com/layeh/piepan"
	"github.com/robertkrimen/otto"
)

func main() {
	// Flags
	username := flag.String("username", "piepan-bot", "username of the bot")
	password := flag.String("password", "", "user password")
	server := flag.String("server", "localhost:64738", "address of the server")
	certificateFile := flag.String("certificate", "", "user certificate file (PEM)")
	keyFile := flag.String("key", "", "user certificate key file (PEM)")
	insecure := flag.Bool("insecure", false, "skip certificate checking")
	lock := flag.String("lock", "", "server certificate lock file")
	serverName := flag.String("servername", "", "override server name used in TLS handshake")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] [scripts...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "a scriptable bot framework for Mumble\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Configuration
	config := gumble.Config{
		Username: *username,
		Password: *password,
		Address:  *server,
	}

	client := gumble.NewClient(&config)

	if *insecure {
		config.TLSConfig.InsecureSkipVerify = true
	}
	if *serverName != "" {
		config.TLSConfig.ServerName = *serverName
	}
	if *lock != "" {
		gumbleutil.CertificateLockFile(client, &config, *lock)
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

	// piepan
	piepan := piepan.New(client)
	piepan.ErrFunc = func(err error) {
		if ottoErr, ok := err.(*otto.Error); ok {
			fmt.Fprintf(os.Stderr, "%s\n", ottoErr.String())
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}

	for _, script := range flag.Args() {
		if err := piepan.LoadScriptFile(script); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}

	client.Attach(gumbleutil.AutoBitrate)
	client.Attach(piepan)

	keepAlive := make(chan bool)
	client.Attach(gumbleutil.Listener{
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
