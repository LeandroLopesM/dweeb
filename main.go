package main

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Message struct {
	author string
	content string
}

type AppCtx struct {
	mu sync.Mutex

	messages []Message
	peers []string

	msgBox *tview.TextView
	msgInput *tview.TextArea
	userList *tview.TextView
}

func main() {
	rootFrame := tview.NewGrid().
		SetRows(0, 4).
		SetColumns(0, 20, 5)
	messageBuffer := tview.NewTextView().
		SetDynamicColors(true)
	messageInput := tview.NewTextArea().
		SetPlaceholder("Write your message here!")
	sidebar := tview.NewTextView()
	submitBtn := tview.NewButton("Send")

	appCtx := AppCtx {
		msgBox: messageBuffer,
		msgInput: messageInput,
		userList: sidebar,
	}
	appCtx.peers = append(appCtx.peers, "invalid:-1")

	submitBtn.
		SetSelectedFunc(func() {
			submitBtn.SetDisabled(true)
			if len(messageInput.GetText()) != 0 {
				go appCtx.submitMessage(messageInput.GetText())
				messageInput.SetText("", true)
			}
		}).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorBlack)


	messageBuffer.SetBorder(true).
		SetTitle("Messages")

	messageInput.SetBorder(true).
		SetTitle("Message")

	sidebar.SetBorder(true).
		SetTitle("Sidebar")

	rootFrame.AddItem(messageBuffer, 0, 0, 1, 1, 0, 0, false)
	rootFrame.AddItem(sidebar, 0, 1, 1, 1, 0, 0, false)
	rootFrame.AddItem(messageInput, 1, 0, 1, 1, 0, 0, true)

	rootFrame.AddItem(submitBtn, 1, 1, 1, 1, 0, 0, false)

	go receiverThread(&appCtx)

	if err := tview.NewApplication().SetRoot(rootFrame, true).EnablePaste(true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
