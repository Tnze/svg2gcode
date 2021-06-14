package convsvg

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

func parsePath(encoder Encoder, element xml.StartElement, offsetX, offsetY, scale float64) error {
	for _, v := range element.Attr {
		if v.Name.Local == "d" {
			if err := encoder.StartPath(); err != nil {
				return err
			}

			fs, err := fieldsD(v.Value)
			if err != nil {
				return err
			}

			var x, y float64
			for i := 0; i < len(fs); {
				cmd, ok := fs[i].(rune)
				if !ok {
					return fmt.Errorf("%T(%v) cannot be a path command", fs[i], fs[i])
				}
				i++

				switch cmd {
				case 'Z', 'z':
					if err = encoder.ClosePath(); err != nil {
						return err
					}

				case 'M', 'm', 'L', 'l':
					for i+2-1 < len(fs) {
						var lx, ly float64
						if lx, ok = fs[i].(float64); !ok {
							break
						}
						if ly, ok = fs[i+1].(float64); !ok {
							return fmt.Errorf("%T(%v) is not a number", fs[i+1], fs[i+1])
						}

						if cmd == 'm' || cmd == 'l' { // 相对坐标
							x += lx
							y += ly
						} else { //绝对坐标
							x, y = lx, ly
						}

						if cmd == 'm' || cmd == 'M' {
							err = encoder.MoveTo(scale*(x+offsetX), scale*(y+offsetY))
						} else {
							err = encoder.LineTo(scale*(x+offsetX), scale*(y+offsetY))
						}
						if err != nil {
							return err
						}
						i += 2
					}

				case 'V', 'v', 'H', 'h':
					var v float64
					if v, ok = fs[i].(float64); !ok {
						return fmt.Errorf("%T(%v) is not a number", fs[i+1], fs[i+1])
					}
					switch cmd {
					case 'V': // 绝对坐标
						y = v
					case 'v': // 相对坐标
						y += v
					case 'H':
						x = v
					case 'h':
						x += v
					}
					if err := encoder.LineTo(scale*(x+offsetX), scale*(y+offsetY)); err != nil {
						return err
					}
					i++

				case 'C', 'c':
					for i+3*2-1 < len(fs) {
						if _, ok = fs[i].(float64); !ok {
							break
						}
						points := [4][2]float64{{x, y}}
						for j := 0; j < 3; j++ {
							var lit [2]float64
							if lit[0], ok = fs[i+j*2].(float64); !ok {
								return fmt.Errorf("%T(%q) is not a number for x", fs[i+j*2], fs[i+j*2])
							}
							if lit[1], ok = fs[i+j*2+1].(float64); !ok {
								return fmt.Errorf("%T(%q) is not a number for y", fs[i+j*2+1], fs[i+j*2+1])
							}

							if cmd == 'c' {
								lit[0] += x
								lit[1] += y
							}
							points[j+1] = lit
						}
						x = points[3][0]
						y = points[3][1]

						for _, bi := range bezierInsert(points) {
							if err = encoder.LineTo(scale*(bi[0]+offsetX), scale*(bi[1]+offsetY)); err != nil {
								return err
							}
						}

						i += 3 * 2
					}
				default:
					return fmt.Errorf("unknown cmd: '%c'", cmd)
				}
			}

			return nil
		}
	}
	return errors.New("didn't find 'd' attr for <path>")
}

// split d attr and return []interface{} of byte and float64
func fieldsD(d string) (v []interface{}, err error) {
	var (
		si int
		s  bool
	)
	for i, c := range d {
		if unicode.IsLetter(c) && c != 'e' {
			if s {
				f, err := strconv.ParseFloat(d[si:i], 64)
				if err != nil {
					return v, err
				}
				v = append(v, f)
			}
			v = append(v, c)
			s = false
		} else if !unicode.IsSpace(c) && c != ',' {
			if !s {
				s = true
				si = i
			}
		} else {
			if s {
				s = false
				f, err := strconv.ParseFloat(d[si:i], 64)
				if err != nil {
					return v, err
				}
				v = append(v, f)
			}
		}
	}

	return
}
