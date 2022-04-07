package static

import (
	"os"
	"path"
)

var (
	home, _    = os.UserHomeDir()
	ConfigPath = path.Join(home, ".config/madamchuckle/madamchuckle.yaml")
	DBPath     = path.Join(home, ".madamchuckle.db")
)

var (
	FieldOfficesEventID = 0
	InvasionsEventID    = 1
	SillyMeterEventID   = 2
	StatusEventID       = 3
	DoodlesEventID      = 4
	PopulationEventID   = 5
)
