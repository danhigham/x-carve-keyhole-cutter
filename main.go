package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/tarm/serial"
)

func main() {

	port := flag.String("port", "/dev/ttyUSB0", "Serial port to connect to CNC")
	spacing := flag.Int("spacing", 200, "Distance between keyholes (mm)")
	length := flag.Int("length", 20, "Length of keyhole (mm)")
	depth := flag.Int("depth", 5, "Depth of keyhole (mm)")

	flag.Parse()

	c := &serial.Config{Name: *port, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	commands := []string{
		"G10 L20 P1 X0 Y0 Z0",                //home
		"G91",                                //inc mode
		"G0 Z10",                             //up 10mm
		fmt.Sprintf("G0 X-%d", *spacing/2),   //move to left hole
		"G90",                                //abs mode
		"G0 Z0",                              //move to zero
		"G91",                                //inc mode
		fmt.Sprintf("G1 F500 Z-%d", *depth),  //cut hole depth
		fmt.Sprintf("G1 F500 Y%d", *length),  //cut hole length
		fmt.Sprintf("G1 F500 Y-%d", *length), //move back
		"G90",                                //abs mode
		"G1 F500 Z10",                        //move out hole
		"G91",                                //inc mode
		fmt.Sprintf("G0 X%d", *spacing),      //move to right hole
		"G90",                                //abs mode
		fmt.Sprintf("G1 F500 Z-%d", *depth),  //cut hole depth
		"G91",                                //inc mode
		fmt.Sprintf("G1 F500 Y%d", *length),  //cut hole length
		fmt.Sprintf("G1 F500 Y-%d", *length), //move back
		"G90",                                //abs mode
		"G1 F500 Z10",                        //move out hole
	}

	for _, c := range commands {
		log.Printf("REQ: \"%s\"\n", c)
		err = send(s, fmt.Sprintf("%s\n", c))

		if err != nil {
			log.Fatal(err)
		}
	}
}

func send(s *serial.Port, m string) error {

	n, err := s.Write([]byte(m))
	if err != nil {
		return err
	}

	buf := make([]byte, 128)
	n, err = s.Read(buf)
	if err != nil {
		return err
	}

	log.Printf("RES: %s", buf[:n])
	return err
}
