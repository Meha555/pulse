package utils

import (
	"runtime/debug"
)

func GetVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main == (debug.Module{}) {
		return "unknown"
	}
	return info.Main.Version
}
