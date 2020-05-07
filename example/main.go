package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/writer"
	"gitlab.com/gomidi/portmididrv"
)

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// This example expects the first input and output port to be connected
// somehow (are either virtual MIDI through ports or physically connected).
// We write to the out port and listen to the in port.
func main() {
	drv, err := portmididrv.New()
	must(err)

	// make sure to close all open ports at the end
	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	outs, err := drv.Outs()
	must(err)

	if len(os.Args) == 2 && os.Args[1] == "list" {
		printInPorts(ins)
		printOutPorts(outs)
		return
	}

	in, out := ins[0], outs[0]

	must(in.Open())
	must(out.Open())

	wr := writer.New(out)

	// listen for MIDI
	rd := reader.New()
	go rd.ListenTo(in)

	{ // write MIDI to out that passes it to in on which we listen.
		err := writer.NoteOn(wr, 60, 100)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Nanosecond)
		writer.NoteOff(wr, 60)
		time.Sleep(time.Nanosecond)

		wr.SetChannel(1)

		writer.NoteOn(wr, 70, 100)
		time.Sleep(time.Nanosecond)
		writer.NoteOff(wr, 70)
		time.Sleep(time.Second * 1)
	}
}

func printPort(port midi.Port) {
	fmt.Printf("[%v] %s\n", port.Number(), port.String())
}

func printInPorts(ports []midi.In) {
	fmt.Printf("MIDI IN Ports\n")
	for _, port := range ports {
		printPort(port)
	}
	fmt.Printf("\n\n")
}

func printOutPorts(ports []midi.Out) {
	fmt.Printf("MIDI OUT Ports\n")
	for _, port := range ports {
		printPort(port)
	}
	fmt.Printf("\n\n")
}
