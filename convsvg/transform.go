package convsvg

type offset struct {
	Encoder
	ox, oy float64
}

func (o offset) MoveTo(x, y float64) error { return o.Encoder.MoveTo(x+o.ox, y+o.oy) }
func (o offset) LineTo(x, y float64) error { return o.Encoder.LineTo(x+o.ox, y+o.oy) }

type scale struct {
	Encoder
	k float64
}

func (s scale) MoveTo(x, y float64) error { return s.Encoder.MoveTo(x*s.k, y*s.k) }
func (s scale) LineTo(x, y float64) error { return s.Encoder.LineTo(x*s.k, y*s.k) }
