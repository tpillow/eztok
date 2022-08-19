package eztok

import (
	"bufio"
	"os"
	"strings"
)

// The string used as the Origin.Name of a string that has been tokenized.
const TokenizeStringOriginName = "<string>"

// Tokenize content using the provided Tokenizer. The Origin.Name of the Token objects
// generated will be TokenizeStringOrigin.
func TokenizeString(tokenizer Tokenizer, content string) ([]*Token, error) {
	return tokenizer.Tokenize(NewReaderContext(
		bufio.NewReader(strings.NewReader(content)), TokenizeStringOriginName))
}

// Tokenize the contents of the file at path using the provided Tokenizer. The Origin.Name
// of the Token objects generated will be path.
func TokenizeFile(tokenizer Tokenizer, path string) ([]*Token, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return tokenizer.Tokenize(NewReaderContext(bufio.NewReader(file), path))
}
