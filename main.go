package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Banner
	printBanner()

	// Load configuration
	fmt.Println("Loading configuration...")
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("ERROR: Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// If no services exist, inform the user
	if len(config.Services) == 0 {
		fmt.Println("\n[INFO] No services defined yet!")
		fmt.Println("To add services, press 'I' key inside the application to import a JSON file.")
		fmt.Println("\nExample JSON format:")
		fmt.Println(`[
  {
    "Category": "Category01",
    "Name": "ProductService",
    "Path": "C:\\Projects\\Services\\ProductService"
  },
  {
    "Category": "Category02",
    "Name": "UserService",
    "Path": "C:\\Projects\\Services\\UserService"
  }
]`)
		fmt.Println("\nPress Enter to continue...")
		fmt.Scanln()
	} else {
		fmt.Printf("✓ %d services loaded\n", len(config.Services))
	}

	// Start Process Manager
	pm := NewProcessManager()

	// Cleanup: stop all processes when application closes
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\n\nClosing application, stopping all services...")
		pm.StopAllServices()
		os.Exit(0)
	}()

	// Start UI
	fmt.Println("Starting UI...\n")
	ui := NewUI(config, pm)

	if err := ui.Start(); err != nil {
		fmt.Printf("ERROR: Failed to start UI: %v\n", err)
		pm.StopAllServices()
		os.Exit(1)
	}

	// Cleanup when application closes
	fmt.Println("\nStopping all services...")
	pm.StopAllServices()
	fmt.Println("✓ Cleanup completed. See you!")
}

// printBanner prints application banner
func printBanner() {
	banner := `
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║        .NET Service Manager                               ║
║        Microservice Management Tool                      ║
║                                                           ║
║        v1.0.0                                             ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}
