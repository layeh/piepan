package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
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
		fmt.Fprintf(os.Stderr, "piepan v0.7.0\n")
		fmt.Fprintf(os.Stderr, "usage: %s [options] [script files]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "an easy to use framework for writing scriptable Mumble bots\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nScript files are defined in the following way:\n")
		fmt.Fprintf(os.Stderr, "  [type%c[environment%c]]filename\n", os.PathListSeparator, os.PathListSeparator)
		fmt.Fprintf(os.Stderr, "    filename: path to script file\n")
		fmt.Fprintf(os.Stderr, "    type: type of script file (default: file extension)\n")
		fmt.Fprintf(os.Stderr, "    environment: name of environment where script will be executed (default: type)\n\n")
		fmt.Fprintf(os.Stderr, "Enabled script types:\n")
		fmt.Fprintf(os.Stderr, "  Type         Name\n")
		for _, ext := range piepan.PluginExtensions {
			fmt.Fprintf(os.Stderr, "  %-12s %s\n", ext, piepan.Plugins[ext].Name)
		}
	}

	flag.Parse()

	// Configuration
	config := gumble.Config{
		Username: *username,
		Password: *password,
		Address:  *server,
		Tokens:   gumble.AccessTokens(accessTokens),
	}

	client := gumble.NewClient(&config)
	instance := piepan.New(client)
	audio, _ := gumble_ffmpeg.New(client)
	audio.Command = *ffmpeg
	instance.Audio = audio

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
		if certificate, err := tls.LoadX509KeyPair(*certificateFile, *keyFile); err != nil {
			panic(err)
		} else {
			config.TLSConfig.Certificates = append(config.TLSConfig.Certificates, certificate)
		}
	}

	client.Attach(gumbleutil.AutoBitrate)

	// Load scripts
	for _, script := range flag.Args() {
		if err := instance.LoadScript(script); err != nil {
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
