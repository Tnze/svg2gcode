package convsvg

type matrix struct {
	Encoder
	m [3][3]float64
}

func NormMatrix(e Encoder) matrix {
	return matrix{
		Encoder: e,
		m: [3][3]float64{
			{1, 0, 0},
			{0, 1, 0},
			{0, 0, 1},
		},
	}
}

func (m matrix) vec(x, y float64) (mx, my float64) {
	mx = m.m[0][0]*x + m.m[0][1]*y + m.m[0][2]*1
	my = m.m[1][0]*x + m.m[1][1]*y + m.m[1][2]*1
	mz := m.m[2][0]*x + m.m[2][1]*y + m.m[2][2]*1
	return mx / mz, my / mz
}

func (m *matrix) scale(x, y float64) {
	m.m[0][0] *= x
	m.m[0][1] *= x
	m.m[0][2] *= x
	m.m[1][0] *= y
	m.m[1][1] *= y
	m.m[1][2] *= y
}

func (m *matrix) translation(x, y float64) {
	m.m[0][0] += x * m.m[2][0]
	m.m[0][1] += x * m.m[2][1]
	m.m[0][2] += x * m.m[2][2]
	m.m[1][0] += y * m.m[2][0]
	m.m[1][1] += y * m.m[2][1]
	m.m[1][2] += y * m.m[2][2]
}

func (m matrix) MoveTo(x, y float64) error { return m.Encoder.MoveTo(m.vec(x, y)) }
func (m matrix) LineTo(x, y float64) error { return m.Encoder.LineTo(m.vec(x, y)) }
