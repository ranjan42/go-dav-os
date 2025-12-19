package keyboard

func inb(port uint16) byte
func outb(port uint16, value byte)

const (
	portData   uint16 = 0x60
	portStatus uint16 = 0x64

	statusOutputBuffer = 1 // bit 0 => output buffer full
)

func readScancode() byte {
	for {
		status := inb(portStatus)
		if (status & statusOutputBuffer) != 0 {
			return inb(portData)
		}
	}
}

func ReadKey() rune {
	for {
		sc := readScancode()

		if sc&0x80 != 0 {
			continue
		}

		if sc < byte(len(LayoutIT)) {
			r := LayoutIT[sc]
			if r != 0 {
				return r
			}
		}
	}
}
