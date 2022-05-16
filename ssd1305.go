package ssd1305

import (
	"errors"
	"image/color"
	"time"

	"tinygo.org/x/drivers"
)

type VccMode uint8

// Device wraps I2C or SPI connection.
type Device struct {
	bus        Buser
	buffer     []byte
	width      int16
	height     int16
	bufferSize int16
	vccState   VccMode
	canReset   bool
}

var _ drivers.Displayer = &Device{}

// Config is the configuration for the display
type Config struct {
	Width    int16
	Height   int16
	VccState VccMode
	Address  uint16
}

type Buser interface {
	configure()
	tx(data []byte, isCommand bool)
	setAddress(address uint16)
	dc(hilo bool)
}

// Configure initializes the display with default configuration
func (d *Device) Configure(cfg Config) {
	if cfg.Width != 0 {
		d.width = cfg.Width
	} else {
		d.width = 128
	}
	if cfg.Height != 0 {
		d.height = cfg.Height
	} else {
		d.height = 32
	}
	if cfg.Address != 0 {
		d.bus.setAddress(cfg.Address)
	}
	if cfg.VccState != 0 {
		d.vccState = cfg.VccState
	} else {
		d.vccState = SWITCHCAPVCC
	}
	d.bufferSize = d.width * d.height / 8
	d.buffer = make([]byte, d.bufferSize)
	d.canReset = cfg.Address != 0 || d.width != 128 || d.height != 64 // I2C or not 128x64

	d.bus.dc(true)
	// Hardware Reset
	d.bus.configure()
	time.Sleep(100 * time.Nanosecond)

	//Set the initialization register
	d.initReg()
	time.Sleep(time.Millisecond * 200)

	// Turn on the OLED display
	d.Command(DISPLAYON)
}

func (d *Device) initReg() {
	d.Command(DISPLAYOFF_RESET)
	d.Command(DISPLAYOFF_04)
	d.Command(DISPLAYOFF_10)
	d.Command(SETLOWCOLUMN)
	d.Command(SETHIGHCOLUMN)
	d.Command(SETSTARTLINE)
	d.Command(SETCONTRAST)
	d.Command(SETSEGOUTPUT_CURRENTBRIGHTNESS)
	d.Command(SEGREMAP)

	d.Command(SETCOMROW_SCANDIRECTION)
	d.Command(NORMALDISPLAY)  //--set normal display  ����ɨ�跴��
	d.Command(MULTIPLEXRATIO) //--set multiplex ratio(1 to 64)
	d.Command(0x00)           //--1/64 duty
	d.Command(0xD5)           //-set display offset	Shift Mapping RAM Counter (0x00~0x3F)
	d.Command(0xF0)           //-not offset
	d.Command(0xD8)           //--set display clock divide ratio/oscillator frequency
	d.Command(0x05)           //--set divide ratio, Set Clock as 100 Frames/Sec
	d.Command(0xD9)           //--set pre-charge period
	d.Command(0xC2)           //Set Pre-Charge as 15 Clocks & Discharge as 1 Clock
	d.Command(0xDA)           //--set com pins hardware configuration
	d.Command(0x12)
	d.Command(0xDB) /*set vcomh*/
	d.Command(0x08) //Set VCOM Deselect Level

	d.Command(0xAF) //-Set Page Addressing Mode (0x00/0x01/0x02)

	d.ClearDisplay()
}

// ClearBuffer clears the image buffer
func (d *Device) ClearBuffer(chFill byte) {
	limit := len(d.buffer)
	for j := 0; j < limit; j++ {
		d.buffer[j] = chFill
	}

}

// ClearDisplay clears the image buffer and clear the display
func (d *Device) ClearDisplay() {
	d.ClearBuffer(0x00)
	d.Display()
}

// Display sends the whole buffer to the screen
func (d *Device) Display() error {
	var row uint8
	for row = 0; row < 4; row++ {
		d.Command(0xB0 + row)
		d.Command(0x04)
		d.Command(0x10)
		d.bus.dc(true)
		for num := 0; num < 128; num++ {
			d.Data(d.buffer[int(row)*128+num])
		}
	}

	return nil
}

// Tx sends data/command to the display
func (d *Device) Tx(data []byte, isCommand bool) {
	d.bus.tx(data, isCommand)
}

// Command is a helper function to write a command to the display
func (d *Device) Command(command uint8) {
	d.Tx([]byte{command}, true)
}

// Data is a helper function to write data to the display
func (d *Device) Data(data byte) {
	d.Tx([]byte{data}, false)
}

// SetPixel enables or disables a pixel in the buffer
// color.RGBA{0, 0, 0, 255} is consider transparent, anything else
// with enable a pixel on the screen
func (d *Device) SetPixel(x int16, y int16, c color.RGBA) {
	if x < 0 || x >= d.width || y < 0 || y >= d.height {
		return
	}
	page := 3 - y/8
	chBx := y % 8
	var chTemp byte = 1 << (7 - chBx)
	if c.R > 0 || c.G > 0 || c.B > 0 {
		d.buffer[128*page+x] |= chTemp
	} else {
		d.buffer[128*page+x] &= ^chTemp
	}

}

// GetPixel returns if the specified pixel is on (true) or off (false)
func (d *Device) GetPixel(x int16, y int16) bool {
	if x < 0 || x >= d.width || y < 0 || y >= d.height {
		return false
	}
	page := 3 - y/8
	return (d.buffer[128*page+x] >> uint8(y%8) & 0x1) == 1
}

// SetBuffer changes the whole buffer at once
func (d *Device) SetBuffer(buffer []byte) error {
	if int16(len(buffer)) != d.bufferSize {
		//return ErrBuffer
		return errors.New("wrong size buffer")
	}
	for i := int16(0); i < d.bufferSize; i++ {
		d.buffer[i] = buffer[i]
	}
	return nil
}

// Size returns the current size of the display.
func (d *Device) Size() (w, h int16) {
	return d.width, d.height
}
