package plugin

import (
	"os"
	"strings"
)

func createEnvVars() map[string]string {
	vars := make(map[string]string)
	for _, val := range os.Environ() {
		split := strings.SplitN(val, "=", 2)
		if len(split) != 2 {
			continue
		}
		vars[split[0]] = split[1]
	}
	return vars
}
