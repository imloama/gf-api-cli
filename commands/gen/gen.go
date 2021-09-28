package gen

import (
	"github.com/gogf/gf-cli/library/mlog"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/text/gstr"
)

func Help() {
	switch gcmd.GetArg(2) {
	case "api":
		HelpApi()
	case "dao":
		HelpDao()
	default:
		mlog.Print(gstr.TrimLeft(`
USAGE 
    gali gen TYPE [OPTION]

TYPE
    api        generate api service dao and model files
    dao        generate dao and model files

DESCRIPTION
    gen api
    gen dao
`))
	}
}

func Run() {
	genType := gcmd.GetArg(2)
	if genType == "" {
		mlog.Print("generating type cannot be empty")
		return
	}
	switch genType {
	case "api":
		doGenApi()
	case "dao":
		doGenDao()
	default:
		panic("未知命令")

	}
}
