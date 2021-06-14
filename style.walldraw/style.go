package style_walldraw

import (
	"fmt"
	"io"

	"github.com/Tnze/svg2gcode/v2/convsvg"
)

const PEN_DOWN = 0.9
const PEN_UP = 0.3

func init() {
	convsvg.Register("walldraw", walldrawStyle{})
}

type walldrawStyle struct{}

func (walldrawStyle) NewEncoder(dst io.Writer) convsvg.Encoder {
	return &walldrawEncoder{dst, false}
}

type walldrawEncoder struct {
	dst       io.Writer
	penStatus bool
}

func (walldrawEncoder) StartPath() error { return nil }
func (walldrawEncoder) ClosePath() error { return nil }

func (s *walldrawEncoder) MoveTo(x, y float64) (err error) {
	if s.penStatus {
		if _, err = fmt.Fprintf(s.dst, "G2 S%f\n", PEN_UP); err != nil {
			return err
		}
		s.penStatus = false
	}
	_, err = fmt.Fprintf(s.dst, "G0 X%f Y%f\n", x, y)
	return
}
func (s *walldrawEncoder) LineTo(x, y float64) (err error) {
	if !s.penStatus {
		if _, err = fmt.Fprintf(s.dst, "G2 S%f\n", PEN_DOWN); err != nil {
			return err
		}
		s.penStatus = true
	}
	_, err = fmt.Fprintf(s.dst, "G0 X%f Y%f\n", x, y)
	return
}
