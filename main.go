package main

import (
	"fmt"
	"image/color"
	"io/fs"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	runApp()
	go func() {

		os.Exit(0)
	}()
}

var root = "/Users/mac/study/go_file_manager"
var fileSystem = os.DirFS(root)

func walkDir(currentPath string) []fs.DirEntry {
	fileList := make([]fs.DirEntry, 0)

	fs.WalkDir(fileSystem, currentPath, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		fileList = append(fileList, dirEntry)
		return nil
	})

	return fileList
}

type UserState struct {
	selectedEntry string
}

func runApp() {
	userState := new(UserState)
	userState.selectedEntry = "."

	myApp := app.New()
	myWindow := myApp.NewWindow("Container")
	green := color.NRGBA{0, 180, 0, 255}

	text1 := canvas.NewText("Hello", green)

	renderButtons := func(fsEntries []fs.DirEntry, createBtnCallBack func(fsEntry fs.DirEntry) func()) []fyne.CanvasObject {
		buttons := make([]fyne.CanvasObject, 0, len(fsEntries))

		for _, fsEntry := range fsEntries {
			btnLabel := fsEntry.Name()
			btnCallback := createBtnCallBack(fsEntry)

			btn := widget.NewButton(btnLabel, func() {
				fmt.Printf("Clicked on the %s\n", btnLabel)
				btnCallback()
			})

			buttons = append(buttons, btn)
		}

		return buttons
	}

	fsEntries := walkDir(userState.selectedEntry)

	buttons := renderButtons(fsEntries, func(fsEntry fs.DirEntry) func() {
		var createBtnCb func(fsEntry fs.DirEntry) func()

		createBtnCb = func(fsEntry fs.DirEntry) func() {
			return func() {
				if fsEntry.IsDir() {
					userState.selectedEntry = fsEntry.Name()
					childFsEntries := walkDir(userState.selectedEntry)

					btns := renderButtons(childFsEntries, createBtnCb)

					containerContent := []fyne.CanvasObject{text1}
					containerContent = append(containerContent, btns...)
					content := container.New(
						layout.NewVBoxLayout(),
						containerContent...,
					)

					myWindow.SetContent(content)
				}
			}
		}

		return func() {
			createBtnCb(fsEntry)()
		}

	})

	containerContent := []fyne.CanvasObject{text1}
	containerContent = append(containerContent, buttons...)

	content := container.New(
		layout.NewVBoxLayout(),
		containerContent...,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
