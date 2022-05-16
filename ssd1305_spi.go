package ssd1305

import (
	"machine"
	"time"

	"tinygo.org/x/drivers"
)

type SPIBus struct {
	wire     drivers.SPI
	dcPin    machine.Pin
	resetPin machine.Pin
	csPin    machine.Pin
}

// NewSPI creates a new SSD1305 connection. The SPI wire must already be configured.
func NewSPI(bus drivers.SPI, dcPin, resetPin, csPin machine.Pin) Device {
	dcPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	resetPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return Device{
		bus: &SPIBus{
			wire:     bus,
			dcPin:    dcPin,
			resetPin: resetPin,
			csPin:    csPin,
		},
	}
}

//

// setAddress does nothing, but it's required to avoid reflection
func (b *SPIBus) setAddress(address uint16) {
	// do nothing
	println("trying to Configure an address on a SPI device")
}

// configure configures some pins with the SPI bus ( RESET )
func (b *SPIBus) configure() {
	b.csPin.Low()
	b.dcPin.Low()
	b.resetPin.Low()

	b.resetPin.High()
	time.Sleep(100 * time.Millisecond)
	b.resetPin.Low()
	time.Sleep(100 * time.Millisecond)
	b.csPin.High()
	b.dcPin.Low()
	b.resetPin.High()
}

// tx sends data to the display (SPIBus implementation)
func (b *SPIBus) tx(data []byte, isCommand bool) {
	if isCommand {
		b.csPin.High()
		b.dcPin.Low()
		b.csPin.Low()

		b.wire.Tx(data, nil)
		b.csPin.High()
	} else {
		//time.Sleep(1 * time.Millisecond)
		b.csPin.High()
		b.dcPin.High()
		b.csPin.Low()

		b.wire.Tx(data, nil)
		b.csPin.High()
	}
}

func (b *SPIBus) dc(hilo bool) {
	if hilo {
		b.dcPin.High()
	} else {
		b.dcPin.Low()
	}
}
