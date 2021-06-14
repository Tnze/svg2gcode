package main

import (
	"flag"
	"log"
	"os"

	"github.com/Tnze/svg2gcode/v2/convsvg"
	_ "github.com/Tnze/svg2gcode/v2/style.simple"
	_ "github.com/Tnze/svg2gcode/v2/style.walldraw"
)

var (
	scale            float64 // 缩放
	offsetX, offsetY float64 // 移动

	outputFileName string // 输出文件名
	styleName      string // 输出gcode风格
)

func main() {
	flag.Float64Var(&scale, "s", 1, "Scale of final G-code")
	flag.Float64Var(&offsetX, "dx", 0, "Offset in X direction")
	flag.Float64Var(&offsetY, "dy", 0, "Offset in Y direction")
	flag.StringVar(&outputFileName, "o", "out.gcode", "Output file name")
	flag.StringVar(&styleName, "style", "simple", "G-code style")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("no input files")
	}

	svgFile, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer svgFile.Close()

	gcodeFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("convsvg: convert start")
	if err := convsvg.Convert(gcodeFile, svgFile, styleName, offsetX, offsetY, scale); err != nil {
		log.Fatal(err)
	}
}
