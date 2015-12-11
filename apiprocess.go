package piepan

import (
	"os/exec"

	"github.com/yuin/gopher-lua"
)

type process struct {
	cmd *exec.Cmd
}

func (p *process) Kill() {
	if proc := p.cmd.Process; proc != nil {
		proc.Kill()
	}
}

func (s *State) apiProcessNew(callback *lua.LFunction, command string, args ...string) *process {
	proc := &process{
		cmd: exec.Command(command, args...),
	}

	go func() {
		var str string
		bytes, err := proc.cmd.Output()
		if err == nil {
			if bytes != nil {
				str = string(bytes)
			}
			s.callValue(callback, proc.cmd.ProcessState.Success(), str)
		} else {
			s.callValue(callback, false, "")
		}
	}()

	return proc
}
