package plugin

import (
	"os/exec"

	"github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
)

type process struct {
	cmd *exec.Cmd
}

func (p *process) Kill() {
	if proc := p.cmd.Process; proc != nil {
		proc.Kill()
	}
}

func (p *Plugin) apiProcessNew(l *lua.State) int {
	callback := luar.NewLuaObject(l, 1)
	command := l.ToString(2)

	args := make([]string, l.GetTop()-2)
	for i := 3; i <= l.GetTop(); i++ {
		args[i-3] = l.ToString(i)
	}

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
			p.callValue(callback, proc.cmd.ProcessState.Success(), str)
		} else {
			p.callValue(callback, false, "")
		}
		callback.Close()
	}()

	obj := luar.NewLuaObjectFromValue(l, proc)
	obj.Push()
	obj.Close()
	return 1
}
