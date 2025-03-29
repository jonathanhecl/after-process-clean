# AfterProcessClean

AfterProcessClean is a lightweight Windows utility that monitors and tracks new processes that start on your system. It's particularly useful for identifying which applications are launched after installing new software or during system operations.

## 📚 Features

- **Real-time Process Monitoring**: Scans your system every 5 seconds to detect new processes
- **Process Tracking**: Keeps track of processes that were running before and after monitoring began
- **Detailed Information**: Captures process paths, runtime duration, and other relevant information
- **Low Resource Usage**: Optimized for minimal memory and CPU usage
- **Windows Support**: Compatible with Windows 7/10/11 (both 32-bit and 64-bit versions)

## 📚 How It Works

AfterProcessClean works by taking a snapshot of running processes when it starts, then continuously monitoring for new processes that appear afterward. When you exit the application (by pressing Ctrl+C), it will display a list of all new processes that were detected during the monitoring session.

## 📚 Usage

1. Simply run the executable:
   ```
   AfterProcessClean.exe
   ```

2. The application will start monitoring processes in the background
3. When you want to see the results, press Ctrl+C
4. A list of new processes will be displayed, showing:
   - Full path to the executable
   - How long the process has been running

## 🛠️ Building from Source

To build AfterProcessClean from source:

1. Ensure you have Go installed on your system
2. Clone this repository
3. Run:
   ```
   go build
   ```

## 🛠️ System Requirements

- Windows 7/10/11 (32-bit or 64-bit)
- Minimal system resources (uses less than 10MB of RAM)

## 🤝 Contributing

Contributions are welcome! Feel free to submit issues and pull requests.

## 🔗 Links

- [GitHub Repository](https://github.com/jonathanhecl/after-process-clean)
- [Report Issues](https://github.com/jonathanhecl/after-process-clean/issues)
- [Releases](https://github.com/jonathanhecl/after-process-clean/releases)