package tealex

type TokenType string

const (
	TokenEol       = "EOL"
	TokenSemicolon = ";"

	TokenValue = "Value" // value

	TokenComment = "Comment" // value
)

type Token struct {
	v string // value

	l int // line
	b int // begin
	e int // end

	t TokenType
}

func (t Token) IsSeparator() bool {
	switch t.t {
	case TokenEol:
		return true
	case TokenSemicolon:
		return true
	}

	return false
}

func (t Token) StartLine() int {
	return t.l
}

func (t Token) StartCharacter() int {
	return t.b
}

func (t Token) EndLine() int {
	return t.l
}

func (t Token) EndCharacter() int {
	return t.e
}

func (t Token) String() string {
	return t.v
}

func (t Token) Line() int {
	return t.l
}

func (t Token) Begin() int {
	return t.b
}

func (t Token) End() int {
	return t.e
}

func (t Token) Type() TokenType {
	return t.t
}
