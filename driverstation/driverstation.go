package driverstation

import (
	"fmt"
	"hash/crc32"
	"log"
	"net"
	"sync"
	"time"
)

// Constants for the driverstation
const (
	version_number = "10020800"
	send_port      = 1110
	recv_port      = 1150
)

// DS is a DriverStation to talk to FRC robots.
type DS struct {
	// Communications
	send net.Conn
	recv *net.UDPConn

	// Control
	m sync.Mutex

	// Data
	team     int32
	loop     int32
	enabled  bool
	state    State
	alliance Alliance
	station  Station
	sync     uint32
}

// New creates a new DriverStation that connects to the robot immediately.
func New(team int32) *DS {
	return &DS{
		team:     team,
		state:    Teleop,
		alliance: Red,
		station:  Station1,
		sync:     2,
	}
}

// Connect sets up the send and receive UDP connections.
func (ds *DS) Connect() error {
	var err error
	ds.send, err = net.Dial("udp", fmt.Sprintf("10.%d.%d.2:%d", ds.team/100, ds.team%100, send_port))
	if err != nil {
		return err
	}

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", recv_port))
	if err != nil {
		return err
	}
	ds.recv, err = net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	return nil
}

// Run sends and receives data with the robot to allow control.
func (ds *DS) Run() {
	send_time := time.Tick(20 * time.Millisecond)
	reads := ds.receive()
	for {
		select {
		case <-send_time:
			ds.m.Lock()
			ds.loop += 1
			ds.m.Unlock()

			_, err := ds.send.Write(ds.packData())
			if err != nil {
				log.Fatal(err)
			}

		case r := <-reads:
			if r.err != nil {
				log.Fatal(r.err)
			}
			if r.n != 1024 {
				log.Fatal("Didn't receive full packet.")
			}
			// TODO: log.Println("Received message.")
		}
	}
}

// SetEnabled enables and disables the robot.
func (ds *DS) SetEnabled(enabled bool) *DS {
	ds.m.Lock()
	defer ds.m.Unlock()
	ds.enabled = enabled
	return ds
}

// SetState switches the robot between Teleop, Auto and Test modes.
func (ds *DS) SetState(state State) *DS {
	ds.m.Lock()
	defer ds.m.Unlock()
	ds.state = state
	return ds
}

// SetAlliance sets the alliance to red or blue.
func (ds *DS) SetAlliance(alliance Alliance) *DS {
	ds.m.Lock()
	defer ds.m.Unlock()
	ds.alliance = alliance
	return ds
}

// SetStation sets the station too 1, 2 or 3.
func (ds *DS) SetStation(station Station) *DS {
	ds.m.Lock()
	defer ds.m.Unlock()
	ds.station = station
	return ds
}

// Resync synchronizes the driverstation with the robot.
func (ds *DS) Resync() *DS {
	ds.m.Lock()
	defer ds.m.Unlock()
	ds.sync = 2
	return ds
}

type read struct {
	buff []byte
	n    int
	err  error
}

func (ds *DS) receive() <-chan *read {
	ch := make(chan *read)
	go func() {
		for {
			buff := make([]byte, 1024)
			n, err := ds.recv.Read(buff)
			ch <- &read{buff, n, err}
		}
	}()
	return ch
}

func (ds *DS) packData() []byte {
	buff := make([]byte, 1024)
	ds.m.Lock()
	defer ds.m.Unlock()

	// Add loops (2 bytes)
	buff[0] = byte(ds.loop >> 8)
	buff[1] = byte(ds.loop)

	// Add Status (4 bytes)
	buff[2] = ds.status()
	buff[3] = 0xFF // Digital Inputs
	buff[4] = byte(ds.team >> 8)
	buff[5] = byte(ds.team)

	// TODO: Alliance R/B, 1/2/3 (2 bytes)
	if ds.alliance == Red {
		buff[6] = 0x52
	} else {
		buff[6] = 0x42
	}
	buff[7] = 0x30 + byte(ds.station) // Station 1, 2 or 3

	// TODO: Joystick data (???)

	// DS Version
	for i, b := range []byte(version_number) {
		buff[72+i] = b
	}

	// crc32 (last 4 bytes)
	crc := crc32.ChecksumIEEE(buff) // IEEE style?
	buff[1020] = byte(crc >> 24)
	buff[1021] = byte(crc >> 16)
	buff[1022] = byte(crc >> 8)
	buff[1023] = byte(crc)
	return buff
}

// Bit flags for the status byte
//
// Bits
// 0: FPGA Checksum
// 1: Test Mode
// 2: Resynch
// 3: FMS Attached
// 4: Auto
// 5: Enabled
// 6: Not E-Stopped
// 7: Reset
const (
	flagFPGAChecksum byte = 1 << iota
	flagTest
	flagResync
	flagFMSAttached
	flagAuto
	flagEnabled
	flagNotEStopped
	flagResetFlag
)

// status returns the status byte.
//
// Note: FPGA Checksum, FMS Attached are never used.
func (ds *DS) status() byte {
	var b byte = flagNotEStopped

	if ds.state == Auto {
		b |= flagAuto
	} else if ds.state == Test {
		b |= flagTest
	}

	if ds.enabled {
		b |= flagEnabled
	}

	if ds.sync > 0 {
		b |= flagResync
		ds.sync -= 1
	}

	return b
}
