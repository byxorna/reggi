package ui

import (
	"io"
	"io/ioutil"
)

type inputFile struct {
	source   string
	contents string
}

func NewInputFile(source string, reader io.Reader) (*inputFile, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return &inputFile{
		source:   source,
		contents: string(data),
	}, nil
}
