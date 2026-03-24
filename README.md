# 🚀 .NET Service Manager

**Manage Your .NET Microservices from the Terminal!**

A terminal-based UI application built with Go that allows you to manage .NET microservices with category-based organization, start/stop services, build/restore projects, and monitor logs in real-time.

## ✨ Features

- 🎯 **Category-Based Organization**: Group your services into custom categories
- ⚡ **Quick Start/Stop**: Manage services individually or in bulk
- 🔨 **Build Operations**: Build, clean, and restore .NET projects directly
- 📊 **Live Log Monitoring**: Monitor service logs in real-time with auto-refresh
- 🌐 **URL Detection**: Automatically detects and displays service URLs/ports
- 💾 **Persistent Configuration**: Service definitions are saved and loaded automatically
- 🧹 **Auto Cleanup**: All services are stopped automatically when the app exits
- 🎨 **User-Friendly Terminal UI**: Color-coded interface with simple keyboard shortcuts
- 🌍 **UTF-8 Support**: Full Unicode support for international characters

## 📋 Requirements

1. **Go (Golang)** - Version 1.21 or higher
   - Download: https://go.dev/dl/
   
2. **.NET SDK** - To run your microservices
   - Download: https://dotnet.microsoft.com/download

3. **Supported Operating Systems**:
   - ✅ Windows (7 and later)
   - ✅ Linux
   - ✅ macOS

## 🛠️ Installation

### 1. Install Go (if not already installed)

Download the installer from [Go's official website](https://go.dev/dl/) and install it.

Verify the installation:
```bash
go version
```

### 2. Clone and Build

```bash
git clone <your-repo-url>
cd dotnet-service-manager
go mod download
go build -o dotnet-service-manager
```

This will create an executable file named `dotnet-service-manager` (on Windows, use `dotnet-service-manager.exe`).

**Cross-Platform Building:**

You can also build for a different OS/architecture:

```bash
# Build for Linux (from Windows)
GOOS=linux GOARCH=amd64 go build -o dotnet-service-manager

# Build for Windows (from Linux/macOS)
GOOS=windows GOARCH=amd64 go build -o dotnet-service-manager.exe

# Build for macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o dotnet-service-manager

# Build for macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o dotnet-service-manager
```

## 🚀 Usage

### First Run

**Windows:**
```bash
dotnet-service-manager.exe
```

**Linux/macOS:**
```bash
./dotnet-service-manager
```

On first run, no services will be defined. The application will prompt you to import services.

### Adding Services

You can add services in two ways:

#### Method 1: Create a JSON File
```json
[
  {
    "Category": "Category1",
    "Name": "ServiceA",
    "Path": "C:\\Projects\\MyApp\\ServiceA",
    "Port": "http://localhost:5001"
  },
  {
    "Category": "Category1",
    "Name": "ServiceB",
    "Path": "C:\\Projects\\MyApp\\ServiceB\\ServiceB.csproj"
  },
  {
    "Category": "Category2",
    "Name": "ServiceC",
    "Path": "C:\\Projects\\MyApp\\ServiceC"
  }
]
```

**Path Examples:**

Windows (use double backslashes):
```json
{
  "Category": "Category1",
  "Name": "ServiceA",
  "Path": "C:\\Projects\\MyApp\\ServiceA\\ServiceA.csproj"
}
```

Linux/macOS (use forward slashes):
```json
{
  "Category": "Category1",
  "Name": "ServiceA",
  "Path": "/home/user/projects/myapp/ServiceA/ServiceA.csproj"
}
```

Or use relative paths (relative to where the app runs):
```json
[
  {
    "Category": "Category1",
    "Name": "ServiceA",
    "Path": "../projects/ServiceA"
  },
  {
    "Category": "Category2",
    "Name": "ServiceC",
    "Path": "C:\\Projects\\MyApp\\ServiceC"
  }
]
```

**Field Descriptions**:
- `Category`: Category name for grouping services
- `Name`: Service display name
- `Path`: Path to the service directory or `.csproj` file
#### Method 2: Import via Application

1. Run the application
2. Press `I` key
3. Enter the full path to your JSON file
4. Press Enter

Your services are now saved and will be loaded automatically on each startup!

## ⌨️ Keyboard Shortcuts

| Key | Description |
|-----|-------------|
| **↑/↓** | Navigate between services |
| **Enter** | Start/Stop selected service |
| **B** | Build selected service (`dotnet build`) |
| **C** | Clean selected service (`dotnet clean`) |
| **R** | Restore selected service (`dotnet restore`) |
| **A** | Start all services |
| **S** | Stop all services |
| **L** | Show logs for selected service |
| **I** | Import services from JSON file |
| **D** | Clear all services (with confirmation) |
| **Tab** | Switch between categories |
| **Q** or **Ctrl+C** | Exit application |
| **ESC** (in log view) | Return to main screen |

## 📁 Project Structure

```
dotnet-service-manager/
│
├── main.go                  # Application entry point
├── config.go                # Configuration management (shared)
├── config_windows.go        # Windows-specific config path
├── config_unix.go           # Linux/macOS-specific config path
├── process.go               # Process management (shared)
├── process_windows.go       # Windows-specific process attributes
├── process_unix.go          # Linux/macOS-specific process attributes
├── ui.go                    # Terminal UI (using tview)
│
├── go.mod                   # Go module definition
├── go.sum                   # Go dependency checksums
│
├── examples/                # Example JSON files
│   └── services.json
│
└── README.md                # This file
```

**Platform-Specific Files:**
- `config_windows.go` / `config_unix.go`: Handle platform-specific config directory locations
- `process_windows.go` / `process_unix.go`: Handle platform-specific process attributes
- Files are automatically selected at compile time based on the target OS

## 📝 Configuration File Location

Your service definitions are stored at:

**Windows:**
```
%APPDATA%\dotnet-service-manager\config.json
```
Example: `C:\Users\USERNAME\AppData\Roaming\dotnet-service-manager\config.json`

**Linux:**
```
~/.config/dotnet-service-manager/config.json
```
Example: `/home/username/.config/dotnet-service-manager/config.json`

**macOS:**
```
~/.config/dotnet-service-manager/config.json
```
Example: `/Users/username/.config/dotnet-service-manager/

On Windows, this is typically: `C:\Users\USERNAME\AppData\Roaming\dotnet-service-manager\config.json`

## 🤝 Contributing

Contributions are welcome! Feel free to submit issues or pull requests.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
