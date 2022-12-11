package cmd

import "flag"

type CmdLineOpts struct {
	ConfigurationProvider string
	ConfigurationSource   string
}

func ParseOptions() *CmdLineOpts {
	configProvider := flag.String("configuration-provider", "file", "Configuration provider")
	configSource := flag.String("configuration-source", "", "Configuration source")

	flag.Parse()

	return &CmdLineOpts{
		ConfigurationProvider: *configProvider,
		ConfigurationSource:   *configSource,
	}
}
