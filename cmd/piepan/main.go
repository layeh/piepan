package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleutil"
	"github.com/layeh/piepan"
)

type strFlagSlice []string

func (s *strFlagSlice) Set(str string) error {
	*s = append(*s, str)
	return nil
}

func (s *strFlagSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func main() {
	// Flags
	username := flag.String("username", "piepan-bot", "username of the bot")
	password := flag.String("password", "", "user password")
	server := flag.String("server", "localhost:64738", "address of the server")
	certificateFile := flag.String("certificate", "", "user certificate file (PEM)")
	keyFile := flag.String("key", "", "user certificate key file (PEM)")
	insecure := flag.Bool("insecure", false, "skip certificate checking")
	lock := flag.String("lock", "", "server certificate lock file")
	ffmpeg := flag.String("ffmpeg", "ffmpeg", "ffmpeg-capable executable for media streaming")
	var accessTokens strFlagSlice
	flag.Var(&accessTokens, "access-token", "server access token (can be defined multiple times)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "piepan v0.9.0\n")
		fmt.Fprintf(os.Stderr, "usage: %s [options] [script files]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "an easy to use framework for writing Mumble bots using Lua\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Configuration
	config := gumble.NewConfig()
	config.Username = *username
	config.Password = *password
	config.Address = *server
	config.Tokens = gumble.AccessTokens(accessTokens)

	client := gumble.NewClient(config)
	instance := piepan.New(client)
	instance.AudioCommand = *ffmpeg

	if *insecure {
		config.TLSConfig.InsecureSkipVerify = true
	}
	if *lock != "" {
		gumbleutil.CertificateLockFile(client, *lock)
	}
	if *certificateFile != "" {
		if *keyFile == "" {
			keyFile = certificateFile
		}
		certificate, err := tls.LoadX509KeyPair(*certificateFile, *keyFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		config.TLSConfig.Certificates = append(config.TLSConfig.Certificates, certificate)
	}

	client.Attach(gumbleutil.AutoBitrate)

	// Load scripts
	for _, script := range flag.Args() {
		if err := instance.LoadFile(script); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", script, err)
		}
	}

	keepAlive := make(chan bool)
	exitStatus := 0
	client.Attach(gumbleutil.Listener{
		Disconnect: func(e *gumble.DisconnectEvent) {
			if e.Type != gumble.DisconnectUser {
				exitStatus = int(e.Type) + 1
			}
			keepAlive <- true
		},
	})

	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	<-keepAlive
	os.Exit(exitStatus)
}
