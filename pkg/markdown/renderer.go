package markdown

type Renderer interface {
	Render(input string) string
}
