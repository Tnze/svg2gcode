package convsvg

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
)

func parsePath(encoder Encoder, element xml.StartElement) error {
	for _, v := range element.Attr {
		if v.Name.Local == "d" {
			if err := encoder.StartPath(); err != nil {
				return err
			}

			tokenizer := newPathTokenizer(v.Value)
			tokenType, token := tokenizer.get()

			var x, y float64
			var cmd byte
			for tokenType != scanEOF {
				// update command, or just use the last one
				if tokenType == scanCommand {
					cmd = token[0]
					tokenType, token = tokenizer.get()
				}
				switch cmd {
				case 'Z', 'z':
					if err := encoder.ClosePath(); err != nil {
						return err
					}

				case 'M', 'm', 'L', 'l':
					var lx, ly float64
					var err error
					if lx, err = readNumber(&tokenizer, &tokenType, &token); err != nil {
						return err
					}
					if ly, err = readNumber(&tokenizer, &tokenType, &token); err != nil {
						return err
					}

					if cmd == 'm' || cmd == 'l' { // 相对坐标
						x += lx
						y += ly
					} else { //绝对坐标
						x, y = lx, ly
					}

					if cmd == 'm' || cmd == 'M' {
						err = encoder.MoveTo(x, y)
					} else {
						err = encoder.LineTo(x, y)
					}
					if err != nil {
						return err
					}

				case 'V', 'v', 'H', 'h':
					var v float64
					var err error
					if v, err = readNumber(&tokenizer, &tokenType, &token); err != nil {
						return err
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
					if err := encoder.LineTo(x, y); err != nil {
						return err
					}

				case 'C', 'c':
					var err error
					points := [4][2]float64{{x, y}}
					for j := 0; j < 3; j++ {
						lit := &points[j+1]
						if lit[0], err = readNumber(&tokenizer, &tokenType, &token); err != nil {
							return fmt.Errorf("parse x error: %v", err)
						}
						if lit[1], err = readNumber(&tokenizer, &tokenType, &token); err != nil {
							return fmt.Errorf("parse y error: %v", err)
						}

						if cmd == 'c' {
							lit[0] += x
							lit[1] += y
						}
					}
					x = points[3][0]
					y = points[3][1]

					for _, bi := range bezierInsert(points) {
						if err = encoder.LineTo(bi[0], bi[1]); err != nil {
							return err
						}
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

// readNumber "reads out" a number
func readNumber(pt *pathtokenizer, tokenType *int, token *string) (v float64, err error) {
	if *tokenType != scanNum {
		return 0, fmt.Errorf("this token is not a Number but %d", *tokenType)
	}
	v, err = strconv.ParseFloat(*token, 64) // parse float64
	*tokenType, *token = pt.get()           // read next token
	return
}
