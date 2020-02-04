package mod

import (
	"os/exec"
	"time"
)

type CommandsHistory struct {
	Commands       []Command `json:"coms"`
	CommandsNumber int       `json:"comnumber"`
}
type Command struct {
	Com      exec.Cmd
	ComID    string
	Category string
	Out      string
	Time     time.Time
}

type CommandsNumber struct {
	CommandsNumber int `json:"comnumber"`
}
