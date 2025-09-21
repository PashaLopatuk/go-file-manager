package main

import (
	"fmt"
	"image/color"
	"io/fs"
	"os/exec"
	"runtime"

	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
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

func openFileDefault(filepath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", filepath)
	case "darwin":
		cmd = exec.Command("open", filepath)
	case "linux":
		cmd = exec.Command("xdg-open", filepath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}

func (s *AppState) render() {
	green := color.NRGBA{0, 180, 0, 255}
	header := canvas.NewText("Current dir: "+s.currentPath, green)

	entries := s.walkDir()
	buttons := make([]fyne.CanvasObject, 0, len(entries))

	for _, entry := range entries {
		label := entry.Name()

		ent := entry
		btn := widget.NewButton(label, func() {
			fmt.Printf("Clicked on: %s\n", ent.Name())
			if ent.IsDir() {
				s.changeDir(ent)
			} else {
				openFileDefault(s.currentPath + "/" + label)
			}
		})

		if entry.IsDir() {
			btn.SetIcon(theme.FolderIcon())
		}

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

	scrollContainer := container.NewScroll(content)

	scrollContainer.SetMinSize(fyne.NewSize(600, 400))

	s.window.SetContent(scrollContainer)
}
