package main

import (
	"os/user"
)

func main() {
	mainView := newMainView()
	user, _ := user.Current()
	mainView.dir0.pathInput.SetText(user.HomeDir)
	mainView.dir1.pathInput.SetText(user.HomeDir)
	mainView.dir0.readDir(user.HomeDir)
	mainView.dir1.readDir(user.HomeDir)
	mainView.app.Run()
}
