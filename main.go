package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type AppState struct {
	currentPath string
	window      fyne.Window
	entries     []fs.DirEntry
	table       *widget.Table
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

	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}

func (s *AppState) walkDir() []fs.DirEntry {
	fullPath := filepath.Join(Root, s.currentPath)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		fmt.Println("Error reading dir:", err)
		return nil
	}
	sort.Slice(entries, func(i int, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

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
	headerPathInput := widget.NewEntry()
	headerPathInput.SetText(s.currentPath)

	headerPathInput.Resize(fyne.NewSize(2000, headerPathInput.Size().Height))

	goBtn := widget.NewButton("Go!", func() {
		path := headerPathInput.Text
		s.currentPath = path
		s.render()
	})

	headerBox := container.NewVBox(headerPathInput, goBtn)

	s.entries = s.walkDir()

	s.table = widget.NewTable(
		func() (int, int) {
			return len(s.entries) + 2, 3
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {

			label := cell.(*widget.Label)

			if id.Row == 0 {
				switch id.Col {
				case 0:
					label.SetText("Name")
				case 1:
					label.SetText("Type")
				case 2:
					label.SetText("Size (bytes)")
				}
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else if id.Row == 1 {
				switch id.Col {
				case 0:
					label.SetText("..")
				case 1:
					label.SetText("Folder")
				case 2:
					label.SetText("")
				}

			} else {
				entry := s.entries[id.Row-2]
				info, _ := entry.Info()

				switch id.Col {
				case 0:
					label.SetText(entry.Name())
				case 1:
					if entry.IsDir() {
						label.SetText("Folder")
					} else {
						label.SetText("File")
					}
				case 2:
					if info != nil && !entry.IsDir() {
						label.SetText(strconv.FormatInt(info.Size(), 10))
					} else {
						label.SetText("-")
					}
				}
				label.TextStyle = fyne.TextStyle{}
			}
		},
	)

	s.table.SetColumnWidth(0, 400)
	s.table.SetColumnWidth(1, 100)
	s.table.SetColumnWidth(2, 100)

	s.table.OnSelected = func(id widget.TableCellID) {
		fmt.Printf("OnSelected: %v", id)

		s.table.Unselect(id)

		if id.Row == 1 {
			pathSplit := strings.Split(s.currentPath, string(os.PathSeparator))
			if len(pathSplit) > 1 {
				backPath := strings.Join(pathSplit[0:len(pathSplit)-1], string(os.PathSeparator))
				if backPath == "" {
					backPath = string(os.PathSeparator)
				}
				s.currentPath = backPath
				s.render()
			}
			return
		} else if id.Row == 2 {
			return
		}
		entry := s.entries[id.Row-1]
		fmt.Printf("entry: %v", entry)
		if entry.IsDir() {
			s.changeDir(entry)
		} else {
			fullPath := filepath.Join(s.currentPath, entry.Name())
			openFileDefault(fullPath)
		}

	}

	content := container.NewBorder(headerBox, nil, nil, nil, s.table)
	s.window.SetContent(content)
}
