package core

import (
	"github.com/Hafuunano/Protocol-ConvertTool/protocol/zerobot"
)

// Init registers message handlers via the protocol layer (zerobot). Plugins are loaded by plugins.go imports.
// Plugins that need a store (e.g. orderCard) use lazy-loaded DefaultCache() on first use if the host did not call SetStore.
func Init() {
	zerobot.Install()
}
