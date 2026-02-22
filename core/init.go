package core

import (
	"github.com/Hafuunano/Protocol-ConvertTool/protocol/zerobot"
)

// Init registers message handlers via the protocol layer (zerobot). Plugins are loaded by plugins.go imports.
func Init() {
	zerobot.Install()
}
