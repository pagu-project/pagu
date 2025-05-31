package markdown

type TelegramRenderer struct{}

func NewTelegramRenderer() *TelegramRenderer {
	return &TelegramRenderer{}
}

func (*TelegramRenderer) Render(input string) string {
	return input
}
