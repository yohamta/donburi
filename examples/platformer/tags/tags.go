package tags

import "github.com/yohamta/donburi"

var (
	Player           = donburi.NewTag().SetName("Player")
	Platform         = donburi.NewTag().SetName("Platform")
	FloatingPlatform = donburi.NewTag().SetName("FloatingPlatform")
	Wall             = donburi.NewTag().SetName("Wall")
	Ramp             = donburi.NewTag().SetName("Ramp")
)
