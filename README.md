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
go build -o dotnet-service-manager.exe
```

This will create an executable file named `dotnet-service-manager.exe` (or `dotnet-service-manager` on Linux/Mac).

## 🚀 Usage

### First Run

```bash
./dotnet-service-manager.exe
```

On first run, no services will be defined. The application will prompt you to import services.

### Adding Services

You can add services in two ways:

#### Method 1: Create a JSON File

Create a JSON file with your service definitions:

**services.json**
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

**Field Descriptions**:
- `Category`: Category name for grouping services
- `Name`: Service display name
- `Path`: Path to the service directory or `.csproj` file
- `Port`: (Optional) Service URL - if not specified, will be auto-detected from logs

> **Note**: On Windows, use double backslashes `\\` in paths

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
├── main.go              # Application entry point
├── config.go            # Configuration management
├── process.go           # Process management and log capture
├── ui.go                # Terminal UI (using tview)
│
├── go.mod               # Go module definition
├── go.sum               # Go dependency checksums
│
├── examples/            # Example JSON files
│   └── services.json
│
└── README.md            # This file
```

## 📝 Configuration File Location

Your service definitions are stored at:
```
%APPDATA%\dotnet-service-manager\config.json
```

On Windows, this is typically: `C:\Users\USERNAME\AppData\Roaming\dotnet-service-manager\config.json`

## 🎯 Usage Example

1. **Start the application**:
   ```bash
   ./dotnet-service-manager.exe
   ```

2. **Services are displayed** grouped by category

3. **To start a service**:
   - Select it with arrow keys
   - Press Enter

4. **To build a service**:
   - Select it with arrow keys
   - Press `B`
   - Check logs with `L`

5. **To view logs**:
   - Select service
   - Press `L`
   - Logs update automatically in real-time
   - Press ESC to return

6. **To switch categories**:
   - Press `Tab`

7. **To start all services in current category**:
   - Press `A`

8. **To exit**:
   - Press `Q`
   - All services are automatically stopped ✅

## 🌐 URL/Port Detection

The application supports two ways to display service URLs:

1. **Manual**: Specify `"Port"` in your JSON file
2. **Automatic**: Detects URLs from service logs (e.g., "Now listening on: http://localhost:5293")

The URL column will:
- Show **green** when a URL is detected or configured
- Show **gray "-"** when no URL is available yet

## 🔨 Build Operations

You can perform common .NET operations on any service:

- **Build** (`B` key): Runs `dotnet build`
- **Clean** (`C` key): Runs `dotnet clean`
- **Restore** (`R` key): Runs `dotnet restore`

All operation output is captured and displayed in the service logs. Press `L` to view logs and see build results.

## 🐛 Troubleshooting

### "go: command not found"
Go is not installed. See installation section.

### "dotnet: command not found"
.NET SDK is not installed or not in PATH.

### Service won't start
- Verify the path is correct
- Check that the directory contains a `.csproj` file
- Try running `dotnet run` manually in that directory

### Path with spaces
Make sure to use double backslashes and the full path:
```json
"Path": "C:\\Program Files\\My App\\Service"
```

## 🔧 Development

### Update dependencies
```bash
go get -u ./...
go mod tidy
```

### Build
```bash
go build -o dotnet-service-manager.exe
```

### Run without building
```bash
go run .
```

## 📚 Code Overview

- **main.go**: Entry point - displays banner, loads config, starts UI
- **config.go**: Handles JSON file reading/writing, service management
- **process.go**: Executes `dotnet` commands, captures logs, detects URLs
- **ui.go**: Creates terminal interface using tview library

## 🙏 Built With

- [tview](https://github.com/rivo/tview) - Terminal UI library (MIT License)
- [tcell](https://github.com/gdamore/tcell) - Terminal cell-based view (Apache 2.0)

## 📄 License

MIT License - See LICENSE file for details

## 🤝 Contributing

Contributions are welcome! Feel free to:
- Report bugs
- Suggest new features
- Submit pull requests

## 🎉 Happy Coding!

If you have any questions or suggestions, feel free to open an issue!
