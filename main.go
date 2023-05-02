package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/LinuxSploit/MediaScript/transcriber"
	"golang.design/x/clipboard"
)

var (
	Application fyne.App
	Window      fyne.Window
	// models short name
	Models []string = []string{}
	// current model full name
	CurrentModelName string
	// output entry
	outputScript *widget.Entry = widget.NewEntry()
)

func getGGMLFiles(dirPath string) []string {
	var ggmlFiles []string
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return ggmlFiles
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "ggml-") {
			ggmlFiles = append(ggmlFiles, file.Name())
		}
	}
	return ggmlFiles
}
func init() {
	Models = getGGMLFiles("./models")
}

func main() {
	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		log.Fatal(err)
	}
	//	app
	Application = app.NewWithID("com.github.linuxsploit")
	Window = Application.NewWindow("MediaScript")
	Window.Resize(fyne.Size{Height: 600, Width: 600})

	Window.SetContent(UIcanvas())
	Window.ShowAndRun()
}

func UIcanvas() fyne.CanvasObject {

	// Fefining UI components first to access every where
	var fileBrowseBtn *widget.Button
	var convertBtn *widget.Button
	var filecard *widget.Card

	// Wav File Browse Button
	fileBrowseBtn = widget.NewButton("Browse", func() {
		fileBrowseBtn.Disable()
		convertBtn.Disable()

		// Browse a file Dialog for .wav only
		fileBrowseDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			// if File open dialog cancelled, then return
			if reader == nil {
				return
			}
			filecard.Subtitle = reader.URI().Path()
			filecard.Refresh()
		}, Window)
		fileBrowseDialog.SetFilter(storage.NewExtensionFileFilter([]string{".wav", ".mp3", ".mp4"}))
		fileBrowseDialog.Show()

		fileBrowseBtn.Enable()
		convertBtn.Enable()
	})
	// File card with info and browse button
	filecard = widget.NewCard("Supported all Media Formats:", "Select a Media file", fileBrowseBtn)

	// Select widget to select a model
	modelSelect := widget.NewSelect(Models, func(s string) {
		CurrentModelName = s
		fmt.Println(CurrentModelName)
	})
	modelSelect.SetSelectedIndex(0)

	// multiline entry for media tab output
	outputScript := widget.NewMultiLineEntry()
	outputScript.Wrapping = fyne.TextWrapWord

	// Tab 2 : Script tab for Output
	scriptTab := container.NewTabItemWithIcon("Script", theme.DocumentSaveIcon(), outputScript)

	// Progress bar
	progressbar := widget.NewProgressBarInfinite()
	progressbar.Hide()
	progressbar.Stop()

	// Convert button for transcribing
	convertBtn = widget.NewButton("Convert", func() {
		//
		atranscriber := transcriber.NewTranScriber()
		log.Println(fileBrowseBtn.Text)
		if err := atranscriber.ReadWav(filecard.Subtitle); err != nil {
			log.Fatal(err)
		}
		progressbar.Show()
		progressbar.Start()
		if err := atranscriber.Transcribe("./models/" + CurrentModelName); err != nil {
			log.Fatal(err)
		}
		progressbar.Stop()
		progressbar.Hide()
		textOut := ""
		for i, line := range atranscriber.ScriptLines {
			textOut = textOut + strconv.Itoa(i+1) + "\n" +
				fmt.Sprintf("%02d:%02d:%02d,%03d", int(line.Start.Hours()), int(line.Start.Minutes())%60, int(line.Start.Seconds())%60, line.Start.Milliseconds()%1000) + " --> " +
				fmt.Sprintf("%02d:%02d:%02d,%03d", int(line.End.Hours()), int(line.End.Minutes())%60, int(line.End.Seconds())%60, line.End.Milliseconds()%1000) + "\n" +
				line.Text + "\n\n"
		}
		outputScript.SetText(textOut)
		//
	})

	// Tab 1 : Media tab for input
	mediaTab := container.NewTabItemWithIcon("Media", theme.MediaPlayIcon(), container.NewVBox(
		filecard,
		modelSelect,
		convertBtn,
		progressbar,
	))
	// about tab item
	fiverr_link, err := url.Parse("https://www.fiverr.com/bilalhameed4556")
	if err != nil {
		log.Fatal(err)
	}
	github_link, err := url.Parse("https://github.com/linuxsploit")
	if err != nil {
		log.Fatal(err)
	}
	// Tabs container which contains both: MediaTab, ScriptTab
	tabs := container.NewAppTabs(
		mediaTab,
		scriptTab,
		container.NewTabItemWithIcon("About", theme.AccountIcon(), container.NewVBox(
			widget.NewLabelWithStyle("Develop by:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Bilal Hameed", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
			widget.NewHyperlinkWithStyle("Fiverr", fiverr_link, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewHyperlinkWithStyle("GitHub", github_link, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		)),
	)
	// copy action toolbar
	Window.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("Action", fyne.NewMenuItem("export .srt", func() {
		fmt.Println("copy")
		Window.Clipboard().SetContent(outputScript.Text)
		///
		dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
			/* if File open dialog cancelled, then return
			if uc == nil {
				return
			}*/
			fmt.Println(uc.URI().Path())
			ioutil.WriteFile(uc.URI().Path(), []byte(outputScript.Text), 0644)
		}, Window).Show()
		///
	}))))

	tabs.SetTabLocation(container.TabLocationLeading)
	return tabs
}
