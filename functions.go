package piepan

import (
	"os"
)

func (in *Instance) disconnect() {
	if client := in.client; client != nil {
		client.Disconnect()
		os.Exit(0)
	}
}
