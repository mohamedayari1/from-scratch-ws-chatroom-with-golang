package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion  string
	RequestTarget string
	Method       string
}

func (r * RequestLine ) validHttp() bool {
	return r.HttpVersion == "1.1" 
}

type Request struct {
	RequestLine RequestLine
}

var ERROR_BAD_START_LINE = fmt.Errorf("malformed request-line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported HTTP version")
var SEPARATOR = "\r\n"

func parseRequestLine(b string) (*RequestLine, string, error) {
	idx := strings.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, b, nil
	}
	startLine := b[:idx]
	restOfMessage := b[idx+len(SEPARATOR):]
	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, restOfMessage, ERROR_BAD_START_LINE
	}
	rl := &RequestLine{
		Method:       parts[0],
		RequestTarget: parts[1],
		HttpVersion:  parts[2],
	}

	if !rl.validHttp() {
		return nil, restOfMessage, ERROR_UNSUPPORTED_HTTP_VERSION
	}
	return rl, restOfMessage, nil

}


func RequestFromReader(reader io.Reader) (*Request, error) {

	data, err := io.ReadAll(reader)
	if err != nil {
		errors.Join(
			fmt.Errorf("failed to read from reader"),
			err,
		)
	}
	rest := string(data)

	rl, _, err := parseRequestLine(rest)
	return &Request{
		RequestLine: *rl,
	}, err
		
}
