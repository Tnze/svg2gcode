package convsvg

import (
	"encoding/xml"
	"log"
)

func parsePath(encoder Encoder, element xml.StartElement) error {
	log.Println(element)
	return nil
}
