# SSD1305 Driver

An adaption of the [SSD1306 driver](https://pkg.go.dev/tinygo.org/x/drivers@v0.20.0/ssd1306).

This driver is for screens using the SSD1305 chip (like the 2.23" OLED screen by [Waveshare](https://www.waveshare.com/wiki/2.23inch_OLED_HAT)).

**N.B.** SPI only

## Example usage
```go
oled := ssd1305.NewSPI(machine.SPI1, dcPin, resetPin, csPin)
oled.Configure(ssd1305.Config{
Width:    128,
Height:   32,
})

oled.ClearDisplay()
tinydraw.FilledRectangle(&oled, 100, 24, 28, 8, ssd1305.BLACK)
oled.Display()
```

