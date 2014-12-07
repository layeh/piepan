package piepan

func (in *Instance) disconnect() {
	if client := in.client; client != nil {
		client.Disconnect()
	}
}
