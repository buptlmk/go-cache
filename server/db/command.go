package db

import (
	"fmt"
	"go-cache/internal"
	"strings"
)

var cmdTable = make(map[string]*command)

type ExecFunc func(db *DB, key string, value interface{}) *internal.Payload

type command struct {
	executor ExecFunc
}

func RegisterCommand(name string, executor ExecFunc) {
	name = strings.ToLower(name)

	cmdTable[name] = &command{
		executor: executor,
	}
}

func ExecCmd(d *DB, command string, key string, value interface{}) *internal.Payload {
	cmd, ok := cmdTable[command]
	if !ok {
		return &internal.Payload{
			Value: fmt.Sprintf("the command:%s is invalid", command),
		}
	}

	return cmd.executor(d, key, value)
}
