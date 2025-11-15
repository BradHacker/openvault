package constants

import (
	"path"

	"github.com/adrg/xdg"
)

var DATA_DIR = path.Join(xdg.DataHome, "openvault")
var CONFIG_DIR = path.Join(xdg.ConfigHome, "openvault")
var CACHE_DIR = path.Join(xdg.CacheHome, "openvault")

var PBKDF2_ROUNDS = 650000
