package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
)

// ServiceProcess represents a running service process
type ServiceProcess struct {
	Service     *Service     // Service information
	Cmd         *exec.Cmd    // Running command
	Status      string       // Status: "Stopped", "Running", "Starting", "Error"
	Logs        []string     // Log messages
	LogLock     sync.RWMutex // Log read/write lock
	Running     bool         // Is it running?
	DetectedURL string       // URL detected from logs (e.g., http://localhost:5060)
}

// ProcessManager manages all service processes
type ProcessManager struct {
	Processes map[string]*ServiceProcess // Service name -> ServiceProcess
	Lock      sync.RWMutex               // Process map lock
}

// NewProcessManager creates a new ProcessManager
func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		Processes: make(map[string]*ServiceProcess),
	}
}

// StartService starts a service
func (pm *ProcessManager) StartService(service *Service) error {
	pm.Lock.Lock()
	defer pm.Lock.Unlock()

	// Check if already running
	if proc, exists := pm.Processes[service.Name]; exists && proc.Running {
		return fmt.Errorf("service is already running")
	}

	// Check if path exists
	if _, err := os.Stat(service.Path); os.IsNotExist(err) {
		return fmt.Errorf("service directory/file not found: %s\nPlease check the Path value in your JSON file", service.Path)
	}

	// Check if path is a .csproj file or directory
	var workDir string
	var cmdArgs []string

	if strings.HasSuffix(strings.ToLower(service.Path), ".csproj") {
		// .csproj file provided, get project directory
		workDir = filepath.Dir(service.Path)
		// Run with dotnet run --project parameter
		cmdArgs = []string{"run", "--project", service.Path}
	} else {
		// Directory provided, use normal dotnet run
		workDir = service.Path
		cmdArgs = []string{"run"}
	}

	// Prepare dotnet run command
	cmd := exec.Command("dotnet", cmdArgs...)
	cmd.Dir = workDir

	// Environment variables for UTF-8 encoding
	cmd.Env = append(cmd.Environ(), "DOTNET_CLI_UI_LANGUAGE=en-US")

	// UTF-8 console encoding for Windows
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    false,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}

	// Create stdout and stderr pipes
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Create ServiceProcess
	proc := &ServiceProcess{
		Service: service,
		Cmd:     cmd,
		Status:  "Starting",
		Logs:    []string{fmt.Sprintf("[%s] Starting service...", service.Name)},
		Running: false,
	}

	pm.Processes[service.Name] = proc

	// Start command
	if err := cmd.Start(); err != nil {
		proc.Status = "Error"
		proc.Logs = append(proc.Logs, fmt.Sprintf("[%s] Error: %v", service.Name, err))
		return fmt.Errorf("failed to start service: %v", err)
	}

	proc.Running = true
	proc.Status = "Running"
	proc.Logs = append(proc.Logs, fmt.Sprintf("[%s] Service started (PID: %d)", service.Name, cmd.Process.Pid))

	// Read logs (goroutines)
	go proc.readLogs(stdout, "OUT")
	go proc.readLogs(stderr, "ERR")

	// Wait for process to finish (goroutine)
	go func() {
		err := cmd.Wait()
		pm.Lock.Lock()
		defer pm.Lock.Unlock()

		proc.Running = false
		if err != nil {
			proc.Status = "Error"
			proc.addLog(fmt.Sprintf("[%s] Service terminated with error: %v", service.Name, err))
		} else {
			proc.Status = "Stopped"
			proc.addLog(fmt.Sprintf("[%s] Service stopped", service.Name))
		}
	}()

	return nil
}

// StopService stops a service
func (pm *ProcessManager) StopService(serviceName string) error {
	pm.Lock.Lock()
	defer pm.Lock.Unlock()

	proc, exists := pm.Processes[serviceName]
	if !exists || !proc.Running {
		return fmt.Errorf("service is not running")
	}

	// Kill the process
	if proc.Cmd != nil && proc.Cmd.Process != nil {
		if err := proc.Cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to stop service: %v", err)
		}
		proc.addLog(fmt.Sprintf("[%s] Stop command sent", serviceName))
	}

	return nil
}

// StopAllServices stops all running services
func (pm *ProcessManager) StopAllServices() {
	pm.Lock.Lock()
	defer pm.Lock.Unlock()

	for _, proc := range pm.Processes {
		if proc.Running && proc.Cmd != nil && proc.Cmd.Process != nil {
			proc.Cmd.Process.Kill()
		}
	}
}

// ExecuteDotnetCommand executes a dotnet command on a service (build, clean, restore, etc.)
func (pm *ProcessManager) ExecuteDotnetCommand(service *Service, command string) error {
	// Check if path exists
	if _, err := os.Stat(service.Path); os.IsNotExist(err) {
		return fmt.Errorf("service directory/file not found: %s", service.Path)
	}

	// Check if path is a .csproj file or directory
	var workDir string
	var cmdArgs []string

	if strings.HasSuffix(strings.ToLower(service.Path), ".csproj") {
		// .csproj file provided
		workDir = filepath.Dir(service.Path)
		cmdArgs = []string{command, service.Path}
	} else {
		// Directory provided
		workDir = service.Path
		cmdArgs = []string{command}
	}

	// Prepare command
	cmd := exec.Command("dotnet", cmdArgs...)
	cmd.Dir = workDir

	// Create or use existing process
	pm.Lock.Lock()
	proc, exists := pm.Processes[service.Name]
	if !exists {
		proc = &ServiceProcess{
			Service: service,
			Status:  "Stopped",
			Logs:    []string{},
			Running: false,
		}
		pm.Processes[service.Name] = proc
	}
	pm.Lock.Unlock()

	// Add log
	proc.addLog(fmt.Sprintf("[%s] Running dotnet %s command...", service.Name, command))

	// Capture output
	output, err := cmd.CombinedOutput()

	// Add output to logs
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if line != "" {
			proc.addLog(fmt.Sprintf("[%s][%s] %s", service.Name, strings.ToUpper(command), line))
		}
	}

	if err != nil {
		proc.addLog(fmt.Sprintf("[%s] Error: %v", service.Name, err))
		return fmt.Errorf("dotnet %s failed: %v", command, err)
	}

	proc.addLog(fmt.Sprintf("[%s] dotnet %s completed ✓", service.Name, command))
	return nil
}

// GetServiceStatus returns a service's status
func (pm *ProcessManager) GetServiceStatus(serviceName string) string {
	pm.Lock.RLock()
	defer pm.Lock.RUnlock()

	if proc, exists := pm.Processes[serviceName]; exists {
		return proc.Status
	}
	return "Stopped"
}

// GetServiceURL returns a service's URL (from JSON or detected from logs)
func (pm *ProcessManager) GetServiceURL(service *Service) string {
	pm.Lock.RLock()
	defer pm.Lock.RUnlock()

	// First check for manually entered port from JSON
	if service.Port != "" {
		return service.Port
	}

	// Then check for URL detected from logs
	if proc, exists := pm.Processes[service.Name]; exists {
		if proc.DetectedURL != "" {
			return proc.DetectedURL
		}
	}

	return "-"
}

// GetServiceLogs returns a service's logs
func (pm *ProcessManager) GetServiceLogs(serviceName string) []string {
	pm.Lock.RLock()
	defer pm.Lock.RUnlock()

	if proc, exists := pm.Processes[serviceName]; exists {
		proc.LogLock.RLock()
		defer proc.LogLock.RUnlock()

		// Return a copy of the logs
		logsCopy := make([]string, len(proc.Logs))
		copy(logsCopy, proc.Logs)
		return logsCopy
	}
	return []string{}
}

// readLogs reads and stores log stream (with UTF-8 support)
func (proc *ServiceProcess) readLogs(reader io.Reader, prefix string) {
	// Direct UTF-8 reading (.NET applications output UTF-8)
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024) // Large buffer (1MB)

	// Regex for URL detection ("Now listening on: http://localhost:5060" etc.)
	urlPattern := regexp.MustCompile(`(?i)now listening on[:\s]+(https?://[^\s]+)`)

	for scanner.Scan() {
		line := scanner.Text()

		// Try to detect URL
		if matches := urlPattern.FindStringSubmatch(line); len(matches) > 1 {
			proc.DetectedURL = matches[1]
		}

		logLine := fmt.Sprintf("[%s][%s] %s", proc.Service.Name, prefix, line)
		proc.addLog(logLine)
	}
}

// addLog adds a log in a thread-safe manner
func (proc *ServiceProcess) addLog(log string) {
	proc.LogLock.Lock()
	defer proc.LogLock.Unlock()

	proc.Logs = append(proc.Logs, log)

	// Limit log count (last 1000 lines)
	if len(proc.Logs) > 1000 {
		proc.Logs = proc.Logs[len(proc.Logs)-1000:]
	}
}
