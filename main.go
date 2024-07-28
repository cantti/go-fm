package main

import (
	"os/user"
)

func main() {
	mainView := newMainView()
	user, _ := user.Current()
	mainView.dir0.dirPath = user.HomeDir
	mainView.dir1.dirPath = user.HomeDir
	mainView.readDir(mainView.dir0)
	mainView.readDir(mainView.dir1)
	mainView.tvapp.Run()
}
