package eztok

import "fmt"

// Represnts information about where something (often a rune or Token)
// came from. Useful for providing user-friendly error information.
type Origin struct {
	// The origin name. Often this is the filepath of where something was
	// parsed from.
	Name string
	// The line number of where something was parsed from.
	LineNum int
	// The column number of where something was parsed from.
	ColNum int
}

// Returns a new Origin object with the given parameters.
func NewOrigin(name string, lineNum int, colNum int) *Origin {
	return &Origin{name, lineNum, colNum}
}

// Returns an error-friendly String representation of the Origin.
func (origin Origin) ToString() string {
	return fmt.Sprintf("%v:%v:%v", origin.Name, origin.LineNum, origin.ColNum)
}
