package eztok

// The type of the callback function required by Node.CanParseToken.
type CanParseTokenCallback func(ctx Context) bool

// The type of the callback function required by Node.ParseToken.
type ParseTokenCallback func(ctx Context) (*Token, error)

// Represents a Node whose CanParseToken and ParseToken functions are
// implemented via function callbacks.
type CallbackNode struct {
	// The function to call when CallbackNode.CanParseToken is called.
	canParseCallback CanParseTokenCallback
	// The function to call when CallbackNode.ParseToken is called.
	parseCallback ParseTokenCallback
}

// Returns a new CallbackNode with the provided parameters.
func NewCallbackNode(canParseCallback CanParseTokenCallback, parseCallback ParseTokenCallback) *CallbackNode {
	return &CallbackNode{canParseCallback, parseCallback}
}

// Calls the CanParseTokenCallback function held by the CallbackNode.
func (node *CallbackNode) CanParseToken(ctx Context) bool {
	return node.canParseCallback(ctx)
}

// Calls the ParseTokenCallback function held by the CallbackNode.
func (node *CallbackNode) ParseToken(ctx Context) (*Token, error) {
	return node.parseCallback(ctx)
}
