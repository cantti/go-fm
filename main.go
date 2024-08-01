package main

import (
	"os/user"
)

func main() {
	mainView := newMainView()
	user, _ := user.Current()
	mainView.dir0.pathInput.SetText(user.HomeDir)
	mainView.dir1.pathInput.SetText(user.HomeDir)
	mainView.readDir(mainView.dir0)
	mainView.readDir(mainView.dir1)
	mainView.app.Run()
}
