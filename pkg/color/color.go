package color

type ColorCode int

const (
	Green  ColorCode = 0x008000
	Red    ColorCode = 0xFF0000
	Yellow ColorCode = 0xFFFF00
	Black  ColorCode = 0x000000
	Pactus ColorCode = 0x052D5A
)

func (ColorCode) String() string {
	return "todo"
}

func (c ColorCode) ToInt() int {
	return int(c)
}
