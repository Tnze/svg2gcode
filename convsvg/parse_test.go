package convsvg

import (
	"testing"
)

func TestTokenizer(t *testing.T) {
	d := `m0.53191,902.65957c14.91971,-5.53805 76.20114,23.18369 104.63804,2.04334c0,0 8.9185,-8.16047 27.27686,-1.51143c29.27465,10.6027 1.38168,-28.26932 34.57447,-13.29786Z`
	tk := newPathTokenizer(d)
	for {
		v, tok := tk.get()
		switch v {
		case scanError:
			t.Fatalf("[%d] %s", v, tok)
		case scanEOF:
			for i := 0; i < 5; i++ {
				tk.get() // Test if it will break
			}
			return
		default:
			t.Logf("[%d] %s", v, tok)
		}
	}
}
