package main

import (
	"github.com/gogf/gf-cli/library/allyes"
	"github.com/gogf/gf-cli/library/mlog"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/text/gstr"
	"github.com/imloama/gf-api-cli/commands/gen"
)


var (
	helpContent = gstr.TrimLeft(`
USAGE
    gali COMMAND [ARGUMENT] [OPTION]

COMMAND

    gen        automatically generate go files for ORM models...
    version    show current binary version info

OPTION
    -y         all yes for all command without prompt ask 
    -?,-h      show this help or detail for specified command
    -v,-i      show version information

ADDITIONAL
    Use 'gali help COMMAND' or 'gf COMMAND -h' for detail about a command, which has '...' 
    in the tail of their comments.
`)
)
func main() {
	allyes.Init()
	command := gcmd.GetArg(1)
	// Help information
	if gcmd.ContainsOpt("h") && command != "" {
		help(command)
		return
	}
	switch command {
	case "help":
		help(gcmd.GetArg(2))
	case "version":
		version()
	case "gen":
		gen.Run()
	default:
		for k := range gcmd.GetOptAll() {
			switch k {
			case "?", "h":
				mlog.Print(helpContent)
				return
			case "i", "v":
				version()
				return
			}
		}
		mlog.Print(helpContent)
	}
}

// help shows more information for specified command.
func help(command string) {
	switch command {
	case "gen":
		gen.Help()
	default:
		mlog.Print(helpContent)
	}
}

// version prints the version information of the cli tool.
func version() {
	mlog.Printf("gali version: 1.1.0")
}
