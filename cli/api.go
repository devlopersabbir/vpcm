package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Manage the local background API server",
}

func getPIDFile() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "vpsm", "api.pid")
}

var apiStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the local API server in the background",
	RunE: func(cmd *cobra.Command, args []string) error {
		pidFile := getPIDFile()

		// Check if already running
		if isAPIRunning() {
			fmt.Println("API server is already running.")
			return nil
		}

		// Find vpsm-api executable in path or use standard /usr/local/bin/vpsm-api
		execPath, err := exec.LookPath("vpsm-api")
		if err != nil {
			execPath = "/usr/local/bin/vpsm-api"
		}

		configHome := filepath.Join(os.Getenv("HOME"), ".config", "vpsm")
		logFile := filepath.Join(configHome, "api.log")
		lf, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		c := exec.Command(execPath)
		c.Stdout = lf
		c.Stderr = lf

		if err := c.Start(); err != nil {
			return fmt.Errorf("failed to start API server: %w", err)
		}

		// Write PID file
		pid := c.Process.Pid
		err = os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0600)
		if err != nil {
			return fmt.Errorf("failed to write pid file: %w", err)
		}

		fmt.Printf("Started local API server in background (PID: %d)\n", pid)
		return nil
	},
}

var apiStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the background API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		pidFile := getPIDFile()
		data, err := os.ReadFile(pidFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("API server is not running (no PID file found).")
				return nil
			}
			return err
		}

		pid, err := strconv.Atoi(string(data))
		if err != nil {
			return fmt.Errorf("invalid PID in file: %w", err)
		}

		proc, err := os.FindProcess(pid)
		if err != nil {
			return fmt.Errorf("failed to find process: %w", err)
		}

		// Send SIGTERM
		err = proc.Signal(syscall.SIGTERM)
		if err != nil {
			// Process might already be dead
			_ = os.Remove(pidFile)
			fmt.Println("API server stopped (process was not active).")
			return nil
		}

		_ = os.Remove(pidFile)
		fmt.Printf("Stopped background API server (PID: %d)\n", pid)
		return nil
	},
}

var apiRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the background API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiStopCmd.RunE(cmd, args); err != nil {
			return err
		}
		return apiStartCmd.RunE(cmd, args)
	},
}

var apiStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the background API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		pidFile := getPIDFile()
		data, err := os.ReadFile(pidFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("API server is: STOPPED")
				return nil
			}
			return err
		}

		pid, err := strconv.Atoi(string(data))
		if err != nil {
			fmt.Println("API server is: STOPPED (corrupt PID file)")
			return nil
		}

		// On Unix-like systems, sending signal 0 checks if process exists
		proc, err := os.FindProcess(pid)
		if err == nil {
			err = proc.Signal(syscall.Signal(0))
			if err == nil {
				fmt.Printf("API server is: RUNNING (PID: %d)\n", pid)
				return nil
			}
		}

		fmt.Println("API server is: STOPPED (inactive process)")
		_ = os.Remove(pidFile)
		return nil
	},
}

func isAPIRunning() bool {
	pidFile := getPIDFile()
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return false
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}

var followLogs bool

var apiLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View local background API server logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		configHome := filepath.Join(os.Getenv("HOME"), ".config", "vpsm")
		logFile := filepath.Join(configHome, "api.log")

		// Check if file exists
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			fmt.Println("No log file found. Starting the API server will generate logs.")
			return nil
		}

		if !followLogs {
			// Read the file and print it
			data, err := os.ReadFile(logFile)
			if err != nil {
				return fmt.Errorf("failed to read log file: %w", err)
			}
			fmt.Print(string(data))
			return nil
		}

		// Follow mode (tail -f)
		file, err := os.Open(logFile)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		defer file.Close()

		// Seek to end
		_, err = file.Seek(0, io.SeekEnd)
		if err != nil {
			return err
		}

		buf := make([]byte, 1024)
		for {
			n, err := file.Read(buf)
			if n > 0 {
				fmt.Print(string(buf[:n]))
			}
			if err == io.EOF {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if err != nil {
				return err
			}
		}
	},
}

func init() {
	apiLogsCmd.Flags().BoolVarP(&followLogs, "follow", "f", false, "Follow log output in real-time")
	apiCmd.AddCommand(apiStartCmd, apiStopCmd, apiRestartCmd, apiLogsCmd, apiStatusCmd)
}
