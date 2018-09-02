# portmididrv

[![Build Status Travis/Linux](https://travis-ci.org/gomidi/portmididrv.svg?branch=master)](http://travis-ci.org/gomidi/portmididrv)

## Purpose

A driver for the unified MIDI driver interface at https://github.com/gomidi/connect .

This driver is based on the portmidi project (see https://github.com/rakyll/portmidi for more information).

For a driver based on rtmidi, see https://github.com/gomidi/rtmididrv .

## Installation

It is recommended to use Go 1.11 with module support (`$GO111MODULE=on`).

## Linux / Debian

```
// install the headers of portmidi somehow, e.g. sudo apt-get install libportmidi-dev
go get -d github.com/gomidi/portmididrv
```

## Documentation

[![portmididrv docs](http://godoc.org/github.com/gomidi/portmididrv?status.png)](http://godoc.org/github.com/gomidi/portmididrv)

## Example

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gomidi/connect"
	"github.com/gomidi/mid"
	driver "github.com/gomidi/portmididrv"
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
	drv, err := driver.New()
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

	wr := mid.WriteTo(out)

	// listen for MIDI
	go mid.NewReader().ReadFrom(in)

	{ // write MIDI to out that passes it to in on which we listen.
		err := wr.NoteOn(60, 100)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Nanosecond)
		wr.NoteOff(60)
		time.Sleep(time.Nanosecond)

		wr.SetChannel(1)

		wr.NoteOn(70, 100)
		time.Sleep(time.Nanosecond)
		wr.NoteOff(70)
		time.Sleep(time.Second * 1)
	}
}

func printPort(port connect.Port) {
	fmt.Printf("[%v] %s\n", port.Number(), port.String())
}

func printInPorts(ports []connect.In) {
	fmt.Printf("MIDI IN Ports\n")
	for _, port := range ports {
		printPort(port)
	}
	fmt.Printf("\n\n")
}

func printOutPorts(ports []connect.Out) {
	fmt.Printf("MIDI OUT Ports\n")
	for _, port := range ports {
		printPort(port)
	}
	fmt.Printf("\n\n")
}

```
