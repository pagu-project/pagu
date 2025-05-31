package markdown

type WhatsAppRenderer struct{}

func NewWhatsAppRenderer() *WhatsAppRenderer {
	return &WhatsAppRenderer{}
}

func (*WhatsAppRenderer) Render(input string) string {
	return input
}
