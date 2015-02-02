package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"

	"github.com/layeh/bconf"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/layeh/gumble/gumbleutil"
	"github.com/layeh/piepan"

	_ "github.com/layeh/piepan/plugins/autobitrate"
	_ "github.com/layeh/piepan/plugins/javascript"
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
		fmt.Fprintf(os.Stderr, "usage: %s [options] [configuration file]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "a bot framework for Mumble\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nDefault configuration file name: piepan.conf\n")
		for _, pluginName := range piepan.PluginNames {
			plugin := piepan.Plugins[pluginName]
			fmt.Fprintf(os.Stderr, "\nPlugin: %s\n%s\n", pluginName, plugin.Help)
		}
	}

	flag.Parse()

	// Configuration
	config := gumble.Config{
		Username: *username,
		Password: *password,
		Address:  *server,
	}

	client := gumble.NewClient(&config)
	instance := piepan.Instance{
		Client: client,
	}
	audio, _ := gumble_ffmpeg.New(client)
	instance.FFmpeg = audio

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

	// Load configuration file
	configurationFileName := "piepan.conf"
	if len(flag.Args()) >= 1 {
		configurationFileName = flag.Arg(0)
	}
	configFile, err := bconf.DecodeFile(configurationFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	for _, block := range configFile.Blocks["plugin"] {
		pluginName := block.Tag.String(0)
		plugin := piepan.Plugins[pluginName]
		if plugin == nil {
			fmt.Fprintf(os.Stderr, "unknown plugin: `%s`\n", pluginName)
			os.Exit(1)
		}
		if err := plugin.Init(&instance, block); err != nil {
			fmt.Fprintf(os.Stderr, "%s plugin error: %s\n", pluginName, err)
			os.Exit(1)
		}
	}

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
