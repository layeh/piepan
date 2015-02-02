package plugin

import (
	"os/exec"

	"github.com/robertkrimen/otto"
)

type process struct {
	cmd *exec.Cmd
}

func (p *process) Kill() {
	p.cmd.Process.Kill()
}

func (in *Instance) apiProcessNew(call otto.FunctionCall) otto.Value {
	callback := call.Argument(0)
	command := call.Argument(1).String()

	args := make([]string, len(call.ArgumentList)-2)
	for i, arg := range call.ArgumentList[2:] {
		args[i] = arg.String()
	}

	p := &process{
		cmd: exec.Command(command, args...),
	}

	go func() {
		var str string
		bytes, _ := p.cmd.Output()
		if bytes != nil {
			str = string(bytes)
		}
		in.callValue(callback, p.cmd.ProcessState.Success(), str)
	}()

	ret, _ := in.state.ToValue(p)
	return ret
}
