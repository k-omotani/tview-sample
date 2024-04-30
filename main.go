package main

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
	"time"
)

func createCommandList() *tview.List {
	commandList := tview.NewList()
	commandList.SetBorder(true).SetTitle("Command")
	return commandList
}

func createInfoPanel(app *tview.Application, u *ui) *tview.Flex {

	infoTable := tview.NewTable()
	infoTable.SetBorder(true).SetTitle("Information")

	cnt := 0
	infoTable.SetCellSimple(cnt, 0, "Data1:")
	infoTable.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	info1 := tview.NewTableCell("aaa")
	infoTable.SetCell(cnt, 1, info1)
	cnt++

	infoTable.SetCellSimple(cnt, 0, "Data2:")
	infoTable.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	info2 := tview.NewTableCell("bbb")
	infoTable.SetCell(cnt, 1, info2)
	cnt++

	infoTable.SetCellSimple(cnt, 0, "Time:")
	infoTable.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	info3 := tview.NewTableCell("0")
	infoTable.SetCell(cnt, 1, info3)
	cnt++

	outputTable := tview.NewTable()
	outputTable.SetBorder(true).SetTitle("Information")

	cnt = 0
	outputTable.SetCellSimple(cnt, 0, "Output:")
	outputTable.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	output := tview.NewTableCell("123")
	outputTable.SetCell(cnt, 1, output)

	infoPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(infoTable, 0, 1, false).
		AddItem(outputTable, 0, 1, false)

	infoTable.SetCellSimple(cnt, 0, "Time:")
	infoTable.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	u.curTime = tview.NewTableCell("0")
	infoTable.SetCell(cnt, 1, u.curTime)

	return infoPanel
}

func createLayout(cList tview.Primitive, recvPanel tview.Primitive) *tview.Flex {
	bodyLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(cList, 20, 1, true).
		AddItem(recvPanel, 0, 1, false)

	header := tview.NewTextView()
	header.SetBorder(true)
	header.SetText("tview study")
	header.SetTextAlign(tview.AlignCenter)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 1, false).
		AddItem(bodyLayout, 0, 1, true)

	return layout
}

func createApplication() *tview.Application {
	app := tview.NewApplication()
	pages := tview.NewPages()

	ui := &ui{}
	ui.app = app

	infoPanel := createInfoPanel(app, ui)

	commandList := createCommandList()
	commandList.AddItem("Test", "test", 'p', testCommand(pages))
	commandList.AddItem("Quit", "quit", 'q', func() {
		app.Stop()
	})
	layout := createLayout(commandList, infoPanel)
	pages.AddPage("main", layout, true, true)

	go updateTime(ui)

	app.SetRoot(pages, true)
	return app
}

func createModalForm(pages *tview.Pages, form tview.Primitive, height int, width int) tview.Primitive {
	modal := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
	return modal
}

func testCommand(pages *tview.Pages) func() {
	return func() {
		cancelFunc := func() {
			pages.SwitchToPage("main")
			pages.RemovePage("modal")
		}

		onFunc := func() {
			pages.SwitchToPage("main")
			pages.RemovePage("modal")
		}

		form := tview.NewForm()
		form.AddButton("ON", onFunc)
		form.AddButton("Cancel", cancelFunc)
		form.SetCancelFunc(cancelFunc)
		form.SetButtonsAlign(tview.AlignCenter)
		form.SetBorder(true).SetTitle("Test")
		modal := createModalForm(pages, form, 13, 55)
		pages.AddPage("modal", modal, true, true)
	}
}

type ui struct {
	app     *tview.Application
	curTime *tview.TableCell
}

const refreshInterval = 500 * time.Millisecond

func currentTimeString() string {
	t := time.Now()
	return fmt.Sprintf(t.Format("15:04:05"))
}

func updateTime(u *ui) {
	for {
		time.Sleep(refreshInterval)
		u.app.QueueUpdateDraw(func() {
			u.curTime.SetText(currentTimeString())
		})
	}
}

func main() {
	runewidth.DefaultCondition = &runewidth.Condition{EastAsianWidth: false}

	app := createApplication()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
