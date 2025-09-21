package main

import (
	"fmt"
	"image/color"
	"io/fs"

	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type AppState struct {
	currentPath string
	window      fyne.Window
}

const Root = "/"

func main() {
	home, _ := os.Getwd()

	myApp := app.New()
	myWindow := myApp.NewWindow("File Explorer")

	state := &AppState{
		currentPath: home,
		window:      myWindow,
	}

	state.render()

	myWindow.ShowAndRun()
}

func (s *AppState) walkDir() []fs.DirEntry {
	fullPath := filepath.Join(Root, s.currentPath)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		panic(err)
	}
	return entries
}

func (s *AppState) changeDir(entry fs.DirEntry) {
	if !entry.IsDir() {
		fmt.Printf("Not a folder: %s\n", entry.Name())
		return
	}

	s.currentPath = filepath.Join(s.currentPath, entry.Name())
	s.render()
}

func (s *AppState) render() {
	green := color.NRGBA{0, 180, 0, 255}
	header := canvas.NewText("Current dir: "+s.currentPath, green)

	entries := s.walkDir()
	buttons := make([]fyne.CanvasObject, 0, len(entries))

	for _, entry := range entries {
		label := entry.Name()
		if entry.IsDir() {
			label = "[Folder] " + label
		}

		ent := entry
		btn := widget.NewButton(label, func() {
			fmt.Printf("Clicked on: %s\n", ent.Name())
			if ent.IsDir() {
				s.changeDir(ent)
			}
		})
		buttons = append(buttons, btn)
	}

	backBtn := widget.NewButton("..", func() {
		pathSplit := strings.Split(s.currentPath, "/")

		backPath := strings.Join(pathSplit[0:len(pathSplit)-1], "/")
		s.currentPath = backPath
		s.render()
	})

	content := container.New(
		layout.NewVBoxLayout(),
		header,
		backBtn,
		container.NewVBox(buttons...),
	)

	s.window.SetContent(content)
}
