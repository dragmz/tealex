package tealex

type Error struct {
	l  int
	li int

	msg string
}

func (e Error) Line() int {
	return e.l
}

func (e Error) Begin() int {
	return e.li
}

func (e Error) End() int {
	return e.li
}

func (e Error) String() string {
	return e.msg
}
