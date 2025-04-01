package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"PZBackup/lib" // Import the backup package (replace "your_project" with the module name)
)

type Config struct {
	SaveFolder     string `json:"save_folder"`
	BackupFolder   string `json:"backup_folder"`
	BackupInterval int    `json:"backup_interval"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the Project Zomboid save folder path: ")
	saveFolder, _ := reader.ReadString('\n')
	saveFolder = strings.TrimSpace(saveFolder)

	fmt.Print("Enter the backup folder path: ")
	backupFolder, _ := reader.ReadString('\n')
	backupFolder = strings.TrimSpace(backupFolder)

	fmt.Print("Enter backup interval in minutes: ")
	var interval int
	_, err := fmt.Scan(&interval)
	if err != nil || interval <= 5 {
		fmt.Println("Invalid input, defaulting to 5 minutes.")
		interval = 5
	}

	// Get max backups to keep
	fmt.Print("Enter the maximum number of backups to keep (default is 3): ")
	var maxBackup int
	_, err = fmt.Scan(&maxBackup)
	if err != nil || maxBackup <= 0 {
		fmt.Println("Invalid input, defaulting to 3 backups.")
		maxBackup = 3
	}

	for {
		lib.ClearOldBackups(backupFolder, maxBackup) // Keep last 5 backups
		lib.Backup(saveFolder, backupFolder)

		fmt.Printf("Waiting %v for the next backup...\n", time.Duration(interval)*time.Minute)
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
