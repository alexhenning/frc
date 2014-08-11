package driverstation

// State indicates whether or not the robot is in teleoperated
// control, autonomous mode or test mode.
type State int

// The three states.
const (
	Teleop State = iota
	Auto
	Test
)

// Alliance indicates red or blue alliance.
type Alliance int

// The two alliances.
const (
	Red Alliance = iota
	Blue
)

// Station indicates the station number
type Station byte

// The three stations.
const (
	Station1 Station = iota + 1
	Station2
	Station3
)
