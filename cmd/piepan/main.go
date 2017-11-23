package main // import "layeh.com/piepan/cmd/piepan"

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"

	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleutil"
	"layeh.com/piepan"
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
	ffmpeg := flag.String("ffmpeg", "ffmpeg", "ffmpeg-capable executable for media streaming")
	var accessTokens strFlagSlice
	flag.Var(&accessTokens, "access-token", "server access token (can be defined multiple times)")
	var scriptArgs strFlagSlice
	flag.Var(&scriptArgs, "script-args", "script arguments for lua from the command line")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "piepan v0.10.0\n")
		fmt.Fprintf(os.Stderr, "usage: %s [options] [script files]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "an easy to use framework for writing Mumble bots using Lua\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Configuration
	config := gumble.NewConfig()
	config.Username = *username
	config.Password = *password
	config.Tokens = gumble.AccessTokens(accessTokens)

	instance := piepan.New(scriptArgs)
	instance.AudioCommand = *ffmpeg

	var tlsConfig tls.Config

	if *insecure {
		tlsConfig.InsecureSkipVerify = true
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
		tlsConfig.Certificates = append(tlsConfig.Certificates, certificate)
	}

	config.Attach(instance)
	config.Attach(gumbleutil.AutoBitrate)

	// Load scripts
	for _, script := range flag.Args() {
		if err := instance.LoadFile(script); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", script, err)
		}
	}

	keepAlive := make(chan bool)
	exitStatus := 0
	config.Attach(gumbleutil.Listener{
		Disconnect: func(e *gumble.DisconnectEvent) {
			if e.Type != gumble.DisconnectUser {
				exitStatus = int(e.Type) + 1
			}
			keepAlive <- true
		},
	})

	_, err := gumble.DialWithDialer(new(net.Dialer), *server, config, &tlsConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		if reject, ok := err.(*gumble.RejectError); ok {
			os.Exit(100 + int(reject.Type))
		}
		os.Exit(99)
	}

	<-keepAlive
	os.Exit(exitStatus)
}
