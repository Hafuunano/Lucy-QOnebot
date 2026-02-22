// Package core (plugins.go): blank imports to load plugins. Each plugin must call protocol.Register(Plugin) in its init().
package core

import (
	_ "github.com/Hafuunano/Plugin-Collections/plugins/plugin-echo"
	_ "github.com/Hafuunano/Plugin-Collections/plugins/plugin-hello"
	_ "github.com/Hafuunano/Plugin-Collections/plugins/plugin-ping"
)
