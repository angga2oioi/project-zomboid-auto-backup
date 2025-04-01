package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"os/exec"
	"strings"
	"time"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

// Configuration struct
type Config struct {
	SaveFolder     string `json:"save_folder"`
	BackupFolder   string `json:"backup_folder"`
	BackupInterval int    `json:"backup_interval"`
}

// Load configuration from a file
func loadConfig() (Config, error) {
	var config Config
	data, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

// Save configuration to a file
func saveConfig(config Config) error {
	dir := "config"
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	filePath := dir + "/config.json"

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func main() {
	// Try to load config from the file
	config, err := loadConfig()
	if err != nil {
		// Fallback to user input if no config is found
		reader := bufio.NewReader(os.Stdin)

		// Ask for save folder location
		fmt.Print("Enter the Project Zomboid save folder path: ")
		saveFolder, _ := reader.ReadString('\n')
		saveFolder = strings.TrimSpace(saveFolder)

		// Ask for backup folder location
		fmt.Print("Enter the backup folder path: ")
		backupFolder, _ := reader.ReadString('\n')
		backupFolder = strings.TrimSpace(backupFolder)

		// Ensure backup is not inside save folder
		if strings.HasPrefix(backupFolder, saveFolder) {
			fmt.Println("Error: Backup location cannot be inside the save folder.")
			return
		}

		// Ask for backup interval
		fmt.Print("Enter backup interval in minutes: ")
		var interval int
		_, err := fmt.Scan(&interval)
		if err != nil {
			fmt.Println("Invalid input, defaulting to 30 minutes.")
			interval = 30
		}

		// Save config for future runs
		config = Config{
			SaveFolder:     saveFolder,
			BackupFolder:   backupFolder,
			BackupInterval: interval,
		}
		err = saveConfig(config)
		if err != nil {
			fmt.Println("Error saving config:", err)
			return
		}
	}

	// Start the backup process
	for {
		clearBackupFolder(config.BackupFolder)
		backup(config.SaveFolder, config.BackupFolder)
		fmt.Printf("Waiting %v for the next backup...\n", time.Duration(config.BackupInterval)*time.Minute)
		time.Sleep(time.Duration(config.BackupInterval) * time.Minute)
	}
}

func clearBackupFolder(backupFolder string) {
	// Ensure the backup directory exists, create it if necessary
	err := os.MkdirAll(backupFolder, 0755)
	if err != nil {
		fmt.Println("Error creating backup folder:", err)
		return
	}

	// Read the directory contents
	dir, err := os.ReadDir(backupFolder)
	if err != nil {
		fmt.Println("Error reading backup folder:", err)
		return
	}

	// Clear the directory contents
	for _, entry := range dir {
		entryPath := filepath.Join(backupFolder, entry.Name())
		err := os.RemoveAll(entryPath)
		if err != nil {
			fmt.Println("Error removing file or directory:", err)
		}
	}
	fmt.Println("Backup folder cleared.")
}



func backup(saveFolder, backupFolder string) {
	// Create timestamped backup folder inside user-defined backup location
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := filepath.Join(backupFolder, "backup_"+timestamp)

	if err := os.MkdirAll(backupPath, os.ModePerm); err != nil {
		fmt.Println("Error creating backup directory:", err)
		return
	}

	// Determine which command to use based on the operating system
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Use xcopy on Windows
		cmd = exec.Command("xcopy", filepath.Join(saveFolder, "*"), backupPath, "/E", "/I", "/Y")
	} else {
		// Use cp on Unix-like systems (Linux, macOS)
		cmd = exec.Command("cp", "-r", filepath.Join(saveFolder, "/"), backupPath)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Creating backup at:", backupPath)
	if err := cmd.Run(); err != nil {
		fmt.Println("Backup failed:", err)
	} else {
		fmt.Println("Backup completed successfully.")
	}
}