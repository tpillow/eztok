package eztok

// The type of the callback function required by some Context utils.
type UntilRuneCallback func(r rune) bool

// Consumes and concatenates runes into a string given the current Context state
// until Context.PeekRune(0) returns NilRune or callback(Context.PeekRune(0))
// returns true.
func ReadRunesUntil(ctx Context, callback UntilRuneCallback) string {
	str := ""
	for r := ctx.PeekRune(0); r != NilRune && !callback(r); r = ctx.PeekRune(0) {
		r = ctx.NextRune()
		str += string(r)
	}
	return str
}

// The inverse of ReadRunesUntil; the reading of runes will stop if
// callback(Context.PeekRune(0)) returns false instead of true.
func ReadRunesUntilNot(ctx Context, callback UntilRuneCallback) string {
	return ReadRunesUntil(ctx, func(r rune) bool { return !callback(r) })
}
