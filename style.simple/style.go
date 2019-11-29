package style_simple

import (
	"fmt"
	"github.com/Tnze/svg2gcode/v2/convsvg"
	"io"
)

func init() {
	convsvg.Register("simple", simpleStyle{})
}

type simpleStyle struct{}

func (simpleStyle) NewEncoder(dst io.Writer) convsvg.Encoder {
	return simpleEncoder{dst}
}

type simpleEncoder struct {
	dst io.Writer
}

func (simpleEncoder) StartPath() error { return nil }
func (simpleEncoder) ClosePath() error { return nil }

func (s simpleEncoder) MoveTo(x, y float64) (err error) {
	_, err = fmt.Fprintf(s.dst, "G0 X%f Y%f\n", x, y)
	return
}
func (s simpleEncoder) LineTo(x, y float64) (err error) {
	_, err = fmt.Fprintf(s.dst, "G1 X%f Y%f\n", x, y)
	return
}
