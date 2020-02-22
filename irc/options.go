package irc

import (
	"flag"
)

// Options specifies the parameters for a client, can be created by Flags which
// parses a set of flags, or by other means.
type Options struct {
	Server                 string
	Password               string
	Nick                   string
	OperUser, OperPassword string
}

// Flags adds the relevant flags from this package, with their names prefixed
// with the given prefix.
func Flags(prefix string, options *Options) {
	flag.StringVar(&options.Server, prefix+"server", "localhost:6667", "Server to connect to, host:port")
	flag.StringVar(&options.Password, prefix+"password", "", "Password to use when connecting to the server (optional)")
	flag.StringVar(&options.Nick, prefix+"nick", "promexp", "Nickname to use")
	flag.StringVar(&options.OperUser, prefix+"oper", "", "Username to use for /OPER (optional)")
	flag.StringVar(&options.OperPassword, prefix+"oper-password", "", "Password to use for /OPER (optional)")
}
