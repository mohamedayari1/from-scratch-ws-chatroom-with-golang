package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)
type parserState string
const (
	StateInit parserState = "init"
	StateDone parserState = "done"
	StateError parserState = "error"
)

type RequestLine struct {
	HttpVersion  string
	RequestTarget string
	Method       string
}


type Request struct {
	RequestLine RequestLine
	state parserState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}


var ErrorMalformedRequestLine = fmt.Errorf("malformed request-line")
var ErrorUnsupportedHTTPVersion = fmt.Errorf("unsupported HTTP version")
var ErrorRequestInErrorState = fmt.Errorf("request in error state")
var SEPARATOR = []byte("\r\n")




func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))

	httpParts := strings.Split(string(parts[2]), "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, 0, ErrorMalformedRequestLine
	}

	
	if len(parts) != 3 {
		return nil, 0, ErrorMalformedRequestLine
	}
	rl := &RequestLine{
		Method:       string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:  string(httpParts[1]),
	}


	return rl, read, nil

}

func (r *Request) parse(data []byte) (int, error) {

	read := 0
outer: 
	for {
		switch r.state {
		case StateError:
			return 0, ErrorRequestInErrorState
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err !=  nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}

				r.RequestLine = *rl
				read += n

				r.state = StateDone

			case StateDone:
				break outer
		}
	}
	return read, nil
}

func (kachta *Request) done() bool {
	return kachta.state == StateDone || kachta.state == StateError
}

func (kachta *Request) error() bool {
	return kachta.state == StateError
}


func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil  {
			return nil, err
		}
		bufLen += n
		readN, err := request.parse(buf[:bufLen])

		if err != nil  {
			return nil, err
		}
		copy(buf, buf[readN: bufLen])
		bufLen -= readN

	}
	return request, nil



}
