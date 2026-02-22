package core

import (
	skillcore "github.com/Hafuunano/Core-SkillAction/core"
	"github.com/Hafuunano/Plugin-Collections/middlewares/whitelist"
	"github.com/Hafuunano/Plugin-Collections/plugins/plugin-agent"
	"github.com/Hafuunano/Protocol-ConvertTool/protocol"
	"github.com/Hafuunano/Protocol-ConvertTool/protocol/zerobot"
)

// Init registers message handlers via the protocol layer (zerobot). Plugins are loaded by plugins.go imports.
// Plugins that need a store (e.g. orderCard, plugin-agent) use lazy-loaded DefaultCache() on first use if the host did not call SetStore.
func Init() {
	cache := skillcore.DefaultCache()
	zerobot.InstallWithMiddlewares([]protocol.Middleware{whitelist.New(cache)})
	pluginagent.SetStore(cache) // LLM config (URL/Key/Model) persisted in same cache/database
}
