package convsvg

import (
	"math"
)

func bezierInsert(b [4][2]float64) (bi [][2]float64) {
	step := math.Min(10, math.Max(3,
		distance(b[0], b[1])+distance(b[1], b[2])+distance(b[2], b[3]), // 近似长度
	))
	for t := 0.0; t <= 1; t += 1.0 / step {
		c1 := (1 - t) * (1 - t) * (1 - t)
		c2 := 3 * t * (1 - t) * (1 - t)
		c3 := 3 * t * t * (1 - t)
		c4 := t * t * t

		bi = append(bi, [2]float64{
			c1*b[0][0] + c2*(b[1][0]) + c3*(b[2][0]) + c4*(b[3][0]),
			c1*b[0][1] + c2*(b[1][1]) + c3*(b[2][1]) + c4*(b[3][1]),
		})
	}
	return
}

func distance(p1, p2 [2]float64) float64 {
	dx := p1[0] - p2[0]
	dy := p1[1] - p2[1]
	return math.Sqrt(dx*dx + dy*dy)
}
