package interrupt

import (
	"os"
	"os/signal"
)

func OnExit(exit func()) {
	go onExit(exit)
}

func onExit(exit func()) {
	notify := make(chan os.Signal)
	signal.Notify(notify, os.Interrupt, os.Kill)

	// block until we get something, then call exit
	<-notify
	exit()
}
