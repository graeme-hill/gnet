package main

import (
	"github.com/graeme-hill/gnet/test/fakeuploader"
)

func main() {
	fustop := make(chan struct{})
	fustopped := make(chan struct{})
	fakeuploader.Run(fustop, fustopped, "http://localhost:10000/upload")
}