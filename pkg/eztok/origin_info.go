package eztok

import "fmt"

const (
	// The origin used for tokenized Go strings.
	StringOrigin = "<string>"
)

// Holds information to determine where a token came from.
type OriginInfo struct {
	// Origin name (i.e. a filename, path, StringOrigin, etc.).
	Origin string
	// Line number from the Origin.
	LineNum int
	// Column number from the Origin.
	ColNum int
}

// Creates a new OriginInfo with the given parameters.
func NewOriginInfo(origin string, lineNum int, colNum int) *OriginInfo {
	return &OriginInfo{origin, lineNum, colNum}
}

// Obtain a string representation of the OriginInfo, used mostly in error messages.
func (info *OriginInfo) ToString() string {
	return fmt.Sprintf("%v:%v:%v", info.Origin, info.LineNum, info.ColNum)
}
