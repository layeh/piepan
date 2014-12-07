package piepan

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
)

type process struct {
	cmd *exec.Cmd
}

func (p *process) Kill() {
	p.cmd.Process.Kill()
}

func (in *Instance) processNew(l *lua.State) int {
	callback := luar.NewLuaObject(l, 1)
	command := l.CheckString(2)

	var args []string
	argCount := l.GetTop()
	for i := 3; i <= argCount; i++ {
		value := l.CheckString(i)
		args = append(args, value)
	}

	p := &process{
		cmd: exec.Command(command, args...),
	}

	go func() {
		defer callback.Close()

		var str string
		bytes, _ := p.cmd.Output()
		if bytes != nil {
			str = string(bytes)
		}
		in.stateLock.Lock()
		defer in.stateLock.Unlock()

		if _, err := callback.Call(p.cmd.ProcessState.Success(), str); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}()

	obj := luar.NewLuaObjectFromValue(l, p)
	obj.Push()
	return 1
}
