// driverstation implements core DriverStation functionality, such as
// enabling the robot and sending basic state information.
//
// This package is intended to provide the underlying DriverStation
// protocol implementation, so that various CLI and GUI interfaces can
// be built on top of it.
//
// TODO:
//     - Allow IPs to be specified
//     - Handle received messages
//     - Reset
//     - Better error handling
//
// UNSUPPORTED:
//
//     - Joystick support
//     - UNSUPPORTED: E-Stop
//     - UNSUPPORTED: Logging
package driverstation
