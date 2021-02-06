package ui

import (
	"fmt"

	"github.com/byxorna/regtest/pkg/version"
)

var (
	bootupMessage = fmt.Sprintf(`Launching %s\nVersion %s (%s)\nCompiled %s\n`, version.Name, version.Version, version.Commit, version.Date)
)
