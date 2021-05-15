package models

type ParamValDecoder string

const (
	None    = ParamValDecoder("")
	Command = ParamValDecoder("command")
	Bash    = ParamValDecoder("bash")
	JSON    = ParamValDecoder("json")
)
