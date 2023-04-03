package main

import (
	_ "xyhelper-web/internal/packed"

	_ "github.com/cool-team-official/cool-admin-go/contrib/drivers/sqlite"

	_ "xyhelper-web/modules"

	"github.com/gogf/gf/v2/os/gctx"

	"xyhelper-web/internal/cmd"
)

func main() {
	// gres.Dump()
	cmd.Main.Run(gctx.New())
}
