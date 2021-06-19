package convsvg

import (
	"encoding/xml"
	"fmt"
	"io"
	"sync"
)

func Convert(dst io.Writer, r io.Reader, styleName string, offsetX, offsetY, k float64) error {
	// find target style
	stylesMu.Lock()
	style, ok := styles[styleName]
	stylesMu.Unlock()
	if !ok {
		return fmt.Errorf("convsvg: unknown style %q (forgotten import?)", styleName)
	}

	// parse svg as xml
	decoder := xml.NewDecoder(r)
	encoder := NormMatrix(style.NewEncoder(dst))
	encoder.translation(offsetX, offsetY)
	encoder.scale(k, k)
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("convsvg: decode token error: %v", err)
		}

		if se, ok := token.(xml.StartElement); ok {
			switch se.Name.Local {
			case "path": // path标签
				err = parsePath(encoder, se)
			}
			if err != nil {
				return fmt.Errorf("convsvg: <%s> parse error: %v", se.Name.Local, err)
			}
		}
	}
	return nil
}

var (
	stylesMu sync.Mutex
	styles   = make(map[string]Style)
)

// Register make a g-code available by the provided name.
// If Register is called twice with the same name or if style is nil,
// it panics.
func Register(name string, s Style) {
	stylesMu.Lock()
	defer stylesMu.Unlock()
	if s == nil {
		panic("convsvg: Register style is nil")
	}
	if _, dup := styles[name]; dup {
		panic("convsvg: Register called twice for driver " + name)
	}
	styles[name] = s
}

type Style interface {
	NewEncoder(dst io.Writer) Encoder
}

type Encoder interface {
	StartPath() error
	MoveTo(x, y float64) error
	LineTo(x, y float64) error
	ClosePath() error
}
