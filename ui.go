package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI manages the application's user interface
type UI struct {
	App               *tview.Application
	Pages             *tview.Pages
	Config            *Config
	ProcessManager    *ProcessManager
	CurrentCategory   string
	ServiceList       *tview.Table
	LogView           *tview.TextView
	StatusBar         *tview.TextView
	HelpText          *tview.TextView
	CurrentLogService string    // Currently displayed service logs
	LogUpdateStop     chan bool // Channel to stop log updates
}

// NewUI creates a new UI
func NewUI(config *Config, pm *ProcessManager) *UI {
	return &UI{
		App:            tview.NewApplication(),
		Config:         config,
		ProcessManager: pm,
	}
}

// Start starts the UI
func (ui *UI) Start() error {
	// Main layout
	ui.Pages = tview.NewPages()

	// Category tabs
	categories := ui.Config.GetCategories()
	if len(categories) == 0 {
		categories = []string{"Default"}
	}

	// Select first category
	ui.CurrentCategory = categories[0]

	// Service list table
	ui.ServiceList = tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetFixed(1, 0)

	ui.ServiceList.SetTitle(" Services ").SetBorder(true)

	// Log view
	ui.LogView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			ui.App.Draw()
		})

	ui.LogView.SetTitle(" Logs ").SetBorder(true)

	// Status bar
	ui.StatusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	// Help text
	ui.HelpText = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("[yellow]Keys: [white]↑/↓[gray]: Select | [white]Enter[gray]: Start/Stop | [white]B[gray]: Build | [white]C[gray]: Clean | [white]R[gray]: Restore | [white]A[gray]: Start All | [white]S[gray]: Stop All | [white]L[gray]: Log | [white]I[gray]: Import | [white]D[gray]: Clear All | [white]Tab[gray]: Category | [white]Q[gray]: Quit")

	// Main layout (grid)
	grid := tview.NewGrid().
		SetRows(0, 1, 1).
		SetColumns(0).
		AddItem(ui.ServiceList, 0, 0, 1, 1, 0, 0, true).
		AddItem(ui.StatusBar, 1, 0, 1, 1, 0, 0, false).
		AddItem(ui.HelpText, 2, 0, 1, 1, 0, 0, false)

	// Separate page for log view
	logPage := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(0).
		AddItem(ui.LogView, 0, 0, 1, 1, 0, 0, true).
		AddItem(tview.NewTextView().SetText("[yellow]ESC[white]: Go Back").SetTextAlign(tview.AlignCenter), 1, 0, 1, 1, 0, 0, false)

	// Add pages
	ui.Pages.AddPage("main", grid, true, true)
	ui.Pages.AddPage("logs", logPage, true, false)

	// Setup key bindings
	ui.setupKeyBindings()

	// Refresh service list
	ui.refreshServiceList()

	// Update status bar
	go ui.updateStatusBar()

	// Start the application
	ui.App.SetRoot(ui.Pages, true).SetFocus(ui.ServiceList)

	return ui.App.Run()
}

// setupKeyBindings sets up keyboard shortcuts
func (ui *UI) setupKeyBindings() {
	ui.ServiceList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := ui.ServiceList.GetSelection()

		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyCtrlC:
			ui.App.Stop()
			return nil

		case tcell.KeyTab:
			// Switch category
			ui.switchCategory()
			return nil

		case tcell.KeyEnter:
			// Start/stop selected service
			if row > 0 {
				ui.toggleService(row - 1)
			}
			return nil
		}

		switch event.Rune() {
		case 'q', 'Q':
			ui.App.Stop()
			return nil

		case 'a', 'A':
			// Start all services
			ui.startAllServices()
			return nil

		case 's', 'S':
			// Stop all services
			ui.stopAllServices()
			return nil

		case 'l', 'L':
			// Switch to log view
			if row > 0 {
				ui.showLogs(row - 1)
			}
			return nil

		case 'b', 'B':
			// Build command
			if row > 0 {
				ui.executeDotnetCommand(row-1, "build")
			}
			return nil

		case 'c', 'C':
			// Clean command
			if row > 0 {
				ui.executeDotnetCommand(row-1, "clean")
			}
			return nil

		case 'r', 'R':
			// Restore command
			if row > 0 {
				ui.executeDotnetCommand(row-1, "restore")
			}
			return nil

		case 'i', 'I':
			// JSON import dialog
			ui.showImportDialog()
			return nil

		case 'd', 'D':
			// Clear all services (Clear All)
			ui.showClearAllDialog()
			return nil
		}

		return event
	})

	ui.LogView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			// Stop log updates
			if ui.LogUpdateStop != nil {
				ui.LogUpdateStop <- true
			}
			ui.Pages.SwitchToPage("main")
			ui.App.SetFocus(ui.ServiceList)
			return nil
		}
		return event
	})
}

// refreshServiceList refreshes the service list
func (ui *UI) refreshServiceList() {
	ui.ServiceList.Clear()

	// Header row
	ui.ServiceList.SetCell(0, 0, tview.NewTableCell("Status").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	ui.ServiceList.SetCell(0, 1, tview.NewTableCell("Service Name").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
	ui.ServiceList.SetCell(0, 2, tview.NewTableCell("URL").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
	ui.ServiceList.SetCell(0, 3, tview.NewTableCell("Path").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))

	// Add services
	services := ui.Config.GetServicesByCategory(ui.CurrentCategory)
	for i, service := range services {
		row := i + 1

		// Status
		status := ui.ProcessManager.GetServiceStatus(service.Name)
		statusCell := tview.NewTableCell(ui.getStatusSymbol(status)).
			SetTextColor(ui.getStatusColor(status)).
			SetAlign(tview.AlignCenter)

		// Service name
		nameCell := tview.NewTableCell(service.Name).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft)

		// URL (from JSON or detected from logs)
		url := ui.ProcessManager.GetServiceURL(&service)
		urlColor := tcell.ColorGray
		if url != "-" {
			urlColor = tcell.ColorGreen
		}
		urlCell := tview.NewTableCell(url).
			SetTextColor(urlColor).
			SetAlign(tview.AlignLeft)

		// Path
		pathCell := tview.NewTableCell(service.Path).
			SetTextColor(tcell.ColorOrange.TrueColor()).
			SetAlign(tview.AlignLeft)

		ui.ServiceList.SetCell(row, 0, statusCell)
		ui.ServiceList.SetCell(row, 1, nameCell)
		ui.ServiceList.SetCell(row, 2, urlCell)
		ui.ServiceList.SetCell(row, 3, pathCell)
	}

	ui.ServiceList.SetTitle(fmt.Sprintf(" Services - %s ", ui.CurrentCategory)).SetBorder(true)
}

// getStatusSymbol returns symbol based on status
func (ui *UI) getStatusSymbol(status string) string {
	switch status {
	case "Running":
		return "●"
	case "Starting":
		return "◐"
	case "Error":
		return "✖"
	default:
		return "○"
	}
}

// getStatusColor returns color based on status
func (ui *UI) getStatusColor(status string) tcell.Color {
	switch status {
	case "Running":
		return tcell.ColorGreen
	case "Starting":
		return tcell.ColorYellow
	case "Error":
		return tcell.ColorRed
	default:
		return tcell.ColorGray
	}
}

// toggleService starts/stops a service
func (ui *UI) toggleService(index int) {
	services := ui.Config.GetServicesByCategory(ui.CurrentCategory)
	if index >= len(services) {
		return
	}

	service := services[index]
	status := ui.ProcessManager.GetServiceStatus(service.Name)

	if status == "Running" || status == "Starting" {
		ui.ProcessManager.StopService(service.Name)
	} else {
		ui.ProcessManager.StartService(&service)
	}

	ui.refreshServiceList()
}

// startAllServices starts all services
func (ui *UI) startAllServices() {
	services := ui.Config.GetServicesByCategory(ui.CurrentCategory)
	for _, service := range services {
		status := ui.ProcessManager.GetServiceStatus(service.Name)
		if status != "Running" && status != "Starting" {
			svc := service
			ui.ProcessManager.StartService(&svc)
		}
	}
	ui.refreshServiceList()
}

// stopAllServices stops all services
func (ui *UI) stopAllServices() {
	services := ui.Config.GetServicesByCategory(ui.CurrentCategory)
	for _, service := range services {
		status := ui.ProcessManager.GetServiceStatus(service.Name)
		if status == "Running" || status == "Starting" {
			ui.ProcessManager.StopService(service.Name)
		}
	}
	ui.refreshServiceList()
}

// switchCategory switches to next category
func (ui *UI) switchCategory() {
	categories := ui.Config.GetCategories()
	if len(categories) == 0 {
		return
	}

	// Switch to next category
	currentIndex := 0
	for i, cat := range categories {
		if cat == ui.CurrentCategory {
			currentIndex = i
			break
		}
	}

	nextIndex := (currentIndex + 1) % len(categories)
	ui.CurrentCategory = categories[nextIndex]
	ui.refreshServiceList()
}

// showLogs displays a service's logs (with live updates)
func (ui *UI) showLogs(index int) {
	services := ui.Config.GetServicesByCategory(ui.CurrentCategory)
	if index >= len(services) {
		return
	}

	service := services[index]
	ui.CurrentLogService = service.Name

	// Stop previous log update
	if ui.LogUpdateStop != nil {
		select {
		case ui.LogUpdateStop <- true:
		default:
		}
		time.Sleep(100 * time.Millisecond) // Wait for old goroutine to stop
	}

	// Create new stop channel
	ui.LogUpdateStop = make(chan bool, 1)

	// Set initial logs directly (don't use QueueUpdateDraw)
	logs := ui.ProcessManager.GetServiceLogs(ui.CurrentLogService)
	ui.LogView.SetTitle(fmt.Sprintf(" %s - Logs (Live) ", ui.CurrentLogService)).SetBorder(true)

	var content string
	if len(logs) == 0 {
		content = "[gray]No logs yet...[white]"
	} else {
		content = strings.Join(logs, "\n")
	}
	ui.LogView.SetText(content)

	// Switch to page
	ui.Pages.SwitchToPage("logs")
	ui.App.SetFocus(ui.LogView)

	// Start live log updates
	go ui.liveLogUpdate()
}

// updateLogView updates the log view
func (ui *UI) updateLogView() {
	if ui.CurrentLogService == "" {
		return
	}

	logs := ui.ProcessManager.GetServiceLogs(ui.CurrentLogService)

	var content string
	if len(logs) == 0 {
		content = "[gray]No logs yet...[white]"
	} else {
		content = strings.Join(logs, "\n")
	}

	ui.App.QueueUpdateDraw(func() {
		ui.LogView.SetText(content)
		ui.LogView.ScrollToEnd()
	})
}

// liveLogUpdate updates logs in real-time
func (ui *UI) liveLogUpdate() {
	ticker := time.NewTicker(500 * time.Millisecond) // Update every 500ms
	defer ticker.Stop()

	for {
		select {
		case <-ui.LogUpdateStop:
			// Stop updating
			return
		case <-ticker.C:
			// Update logs
			ui.updateLogView()
		}
	}
}

// executeDotnetCommand executes a dotnet command on a service (build, clean, restore)
func (ui *UI) executeDotnetCommand(index int, command string) {
	services := ui.Config.GetServicesByCategory(ui.CurrentCategory)
	if index >= len(services) {
		return
	}

	service := services[index]

	// Run command in background
	go func() {
		err := ui.ProcessManager.ExecuteDotnetCommand(&service, command)
		if err != nil {
			// Add log on error
			ui.App.QueueUpdateDraw(func() {
				ui.showError(fmt.Sprintf("dotnet %s error: %v", command, err))
			})
		}
		// Update service list
		ui.App.QueueUpdateDraw(func() {
			ui.refreshServiceList()
		})
	}()

	// Inform user
	ui.showInfo(fmt.Sprintf("Running dotnet %s command...", command))
}

// showInfo displays an information message
func (ui *UI) showInfo(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.Pages.RemovePage("info")
		})

	ui.Pages.AddPage("info", modal, true, true)

	// Auto-close after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		ui.App.QueueUpdateDraw(func() {
			ui.Pages.RemovePage("info")
		})
	}()
}

// showImportDialog displays JSON import dialog
func (ui *UI) showImportDialog() {
	inputField := tview.NewInputField().
		SetLabel("JSON File Path: ").
		SetFieldWidth(50)

	form := tview.NewForm().
		AddFormItem(inputField).
		AddButton("Import", func() {
			path := inputField.GetText()
			if path != "" {
				if err := ui.Config.ImportServicesFromJSON(path); err != nil {
					ui.showError(fmt.Sprintf("Import error: %v", err))
				} else {
					ui.Config.SaveConfig()
					ui.refreshServiceList()
					ui.Pages.SwitchToPage("main")
					ui.App.SetFocus(ui.ServiceList)
				}
			}
		}).
		AddButton("Cancel", func() {
			ui.Pages.SwitchToPage("main")
			ui.App.SetFocus(ui.ServiceList)
		})

	form.SetBorder(true).SetTitle(" JSON Import ").SetTitleAlign(tview.AlignLeft)

	ui.Pages.AddPage("import", form, true, true)
}

// showClearAllDialog displays confirmation dialog for clearing all services
func (ui *UI) showClearAllDialog() {
	serviceCount := ui.Config.GetServiceCount()

	message := fmt.Sprintf(
		"CLEAR ALL SERVICES\n\n"+
			"Total %d services will be deleted!\n"+
			"This action cannot be undone.\n\n"+
			"Do you want to continue?",
		serviceCount,
	)

	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Yes, Clear", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 {
				// Clear all services
				ui.Config.ClearAllServices()
				ui.Config.SaveConfig()
				ui.refreshServiceList()
				ui.showInfo("All services cleared!")
			}
			ui.Pages.RemovePage("clearall")
			ui.Pages.SwitchToPage("main")
			ui.App.SetFocus(ui.ServiceList)
		})

	ui.Pages.AddPage("clearall", modal, true, true)
}

// showError displays an error message
func (ui *UI) showError(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.Pages.SwitchToPage("main")
			ui.App.SetFocus(ui.ServiceList)
		})

	ui.Pages.AddPage("error", modal, true, true)
}

// updateStatusBar updates the status bar
func (ui *UI) updateStatusBar() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		services := ui.Config.GetServicesByCategory(ui.CurrentCategory)
		running := 0
		stopped := 0
		errors := 0

		for _, service := range services {
			status := ui.ProcessManager.GetServiceStatus(service.Name)
			switch status {
			case "Running", "Starting":
				running++
			case "Error":
				errors++
			default:
				stopped++
			}
		}

		statusText := fmt.Sprintf(
			" [green]Running: %d[white] | [gray]Stopped: %d[white] | [red]Error: %d[white] | Category: [yellow]%s[white] ",
			running, stopped, errors, ui.CurrentCategory,
		)

		ui.App.QueueUpdateDraw(func() {
			ui.StatusBar.SetText(statusText)
			ui.refreshServiceList()
		})
	}
}
