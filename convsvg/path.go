package convsvg

type pathtokenizer struct {
	data string
	i, v int
	stat func(pt *pathtokenizer, c byte) int
}

func newPathTokenizer(data string) (p pathtokenizer) {
	p.data = data
	p.stat = statStart
	p.v = p.stat(&p, data[p.i])
	p.i++
	return
}

func (p *pathtokenizer) get() (v int, token string) {
	if p.v == scanEOF {
		return scanEOF, ""
	}
	for p.v == scanSpace {
		p.v = p.stat(p, p.data[p.i])
		p.i++
	}
	start := p.i - 1
	v = p.v
	for p.v == v {
		switch {
		case p.i == len(p.data):
			p.i++
			fallthrough
		case p.i > len(p.data):
			p.v = scanEOF
		default:
			p.v = p.stat(p, p.data[p.i])
			p.i++
		}
	}
	token = p.data[start : p.i-1]
	return
}

const (
	scanCommand = iota
	scanSpace
	scanNum
	scanError
	scanEOF
)

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n'
}

func statStart(pt *pathtokenizer, c byte) int {
	if isSpace(c) {
		pt.stat = statStart
		return scanSpace
	}
	switch c {
	case 'M', 'm', // MoveTo
		'L', 'l', 'H', 'h', 'V', 'v', // LineTo
		'C', 'c', 'S', 's', // Cubic Bézier Curve
		'Q', 'q', 'T', 't', // Quadratic Bézier Curve
		'A', 'a', // Elliptical Arc Curve
		'Z', 'z': // ClosePath
		pt.stat = statStartNum
		return scanCommand
	default:
		return statStartNum(pt, c)
	}
}

func statStartNum(pt *pathtokenizer, c byte) int {
	if isSpace(c) {
		pt.stat = statStartNum
		return scanSpace
	}
	switch c {
	case '-':
		pt.stat = statNum
		return scanNum
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		pt.stat = statNum
		return statNum(pt, c)
	case '.':
		pt.stat = statNumDot
		return scanNum
	}
	return statErr(pt, c)
}

func statNum(pt *pathtokenizer, c byte) int {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return scanNum
	case '.':
		pt.stat = statNumDot
		return scanNum
	case 'e':
		pt.stat = statNumStartE
		return scanNum
	}
	return statEndNum(pt, c)
}

func statNumDot(pt *pathtokenizer, c byte) int {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return scanNum
	case 'e':
		pt.stat = statNumStartE
		return scanNum
	}
	return statEndNum(pt, c)
}

func statNumStartE(pt *pathtokenizer, c byte) int {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return statNumE(pt, c)
	case '-':
		pt.stat = statNumE
		return scanNum
	}
	return statErr(pt, c)
}

func statNumE(pt *pathtokenizer, c byte) int {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		pt.stat = statNumE
		return scanNum
	}
	return statEndNum(pt, c)
}

func statEndNum(pt *pathtokenizer, c byte) int {
	if isSpace(c) || c == ',' {
		pt.stat = statEndNum
		return scanSpace
	}
	return statStart(pt, c)
}

func statErr(pt *pathtokenizer, c byte) int {
	pt.stat = statErr
	return scanError
}
