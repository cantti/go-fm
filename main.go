package main

import (
	"log"
	"os/user"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	mainView := newMainView()
	user, _ := user.Current()
	mainView.dir0.pathInput.SetText(user.HomeDir)
	mainView.dir1.pathInput.SetText(user.HomeDir)
	mainView.dir0.readDir(user.HomeDir)
	mainView.dir1.readDir(user.HomeDir)
	mainView.app.Run()
}
