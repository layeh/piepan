package plugin

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/robertkrimen/otto"
)

type file struct {
	f *os.File
}

func (f *file) Close() {
	f.f.Close()
}

func (f *file) Write(call otto.FunctionCall) otto.Value {
	val, _ := f.f.WriteString(call.Argument(0).String())
	ret, _ := call.Otto.ToValue(val)
	return ret
}

func (f *file) Seek(call otto.FunctionCall) otto.Value {
	offset, _ := call.Argument(0).ToInteger()
	whence := call.Argument(1).String()

	if whence == "" {
		whence = "cur"
	}
	var iWhence int
	switch whence {
	case "set":
		iWhence = 0
	case "cur":
		iWhence = 1
	case "end":
		iWhence = 2
	default:
		return otto.UndefinedValue()
	}

	val, _ := f.f.Seek(offset, iWhence)
	ret, _ := call.Otto.ToValue(val)
	return ret
}

func (f *file) Read(call otto.FunctionCall) otto.Value {
	bytes, err := call.Argument(0).ToInteger()
	var data []byte
	if bytes <= 0 {
		data, err = ioutil.ReadAll(f.f)
		if err != nil {
			return otto.UndefinedValue()
		}
	} else {
		data = make([]byte, bytes)
		_, err = io.ReadFull(f.f, data)
		if err != nil {
			panic(err)
			return otto.UndefinedValue()
		}
	}

	ret, _ := call.Otto.ToValue(string(data))
	return ret
}

func (p *Plugin) apiFileOpen(call otto.FunctionCall) otto.Value {
	var filename, mode string
	switch len(call.ArgumentList) {
	case 1:
		filename = call.Argument(0).String()
		mode = "r"
	case 2:
		filename = call.Argument(0).String()
		mode = call.Argument(1).String()
	default:
		return otto.UndefinedValue()
	}

	var iMode int

	switch mode {
	case "r":
		iMode = os.O_RDONLY
	case "r+":
		iMode = os.O_RDWR
	case "w":
		iMode = os.O_WRONLY | os.O_CREATE
	case "w+":
		iMode = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	case "a":
		iMode = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	case "a+":
		iMode = os.O_RDWR | os.O_CREATE | os.O_APPEND
	default:
		return otto.UndefinedValue()
	}

	osFile, err := os.OpenFile(filename, iMode, 0644)
	if err != nil {
		return otto.UndefinedValue()
	}

	f := &file{
		f: osFile,
	}

	ret, _ := p.state.ToValue(f)
	return ret
}
