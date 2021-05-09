package cmd

import (
	"github.com/spf13/pflag"

	"github.com/vitalik-malkin/go-labs/backoop/internal"
)

type makeCmdConfig struct {
	PlanFile string
}

func (cfg *makeCmdConfig) BindFlags(prefix string) *pflag.FlagSet {
	if prefix != "" {
		prefix = internal.CmdConfigFlagNameDelimiter + prefix
	}

	flagSet := pflag.NewFlagSet("", pflag.PanicOnError)

	flagSet.StringVar(&cfg.PlanFile, prefix+"plan-file", "", "file path to the backup plan")

	return flagSet
}
