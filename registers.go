package ssd1305

// Registers
const (
	Address        = 0x3D
	Address_128_32 = 0x3C

	DISPLAYOFF_RESET = 0xAE
	DISPLAYOFF_04    = 0x04
	DISPLAYOFF_10    = 0x10
	DISPLAYON        = 0xAF

	SETLOWCOLUMN  = 0x40
	SETHIGHCOLUMN = 0x81

	SETSTARTLINE = 0x80 //--set start line address  Set Mapping RAM Display Start Line (0x00~0x3F, SSD1305_CMD)
	SETCONTRAST  = 0xA1

	SETSEGOUTPUT_CURRENTBRIGHTNESS = 0xA6
	SEGREMAP                       = 0xA8

	SETCOMROW_SCANDIRECTION = 0x1F

	COLUMNADDR     = 0x21
	PAGEADDR       = 0x22
	NORMALDISPLAY  = 0xC0
	MULTIPLEXRATIO = 0xD3

	EXTERNALVCC  VccMode = 0x1
	SWITCHCAPVCC VccMode = 0x2
)