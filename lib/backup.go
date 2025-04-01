package lib

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

// Perform backup
func Backup(saveFolder, backupFolder string) {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := filepath.Join(backupFolder, "backup_"+timestamp)

	if err := os.MkdirAll(backupPath, os.ModePerm); err != nil {
		fmt.Println("Error creating backup directory:", err)
		return
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("xcopy", filepath.Join(saveFolder, "*"), backupPath, "/E", "/I", "/Y")
	} else {
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

// Function to keep only the last N backups
func ClearOldBackups(backupFolder string, maxBackups int) {
	err := os.MkdirAll(backupFolder, 0755)
	if err != nil {
		fmt.Println("Error creating backup folder:", err)
		return
	}

	entries, err := os.ReadDir(backupFolder)
	if err != nil {
		fmt.Println("Error reading backup folder:", err)
		return
	}

	var backups []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "backup_") {
			backups = append(backups, entry)
		}
	}

	if len(backups) > maxBackups {
		sort.Slice(backups, func(i, j int) bool {
			return backups[i].Name() < backups[j].Name()
		})

		toDelete := len(backups) - maxBackups
		for i := 0; i < toDelete; i++ {
			oldestBackup := filepath.Join(backupFolder, backups[i].Name())
			fmt.Println("Deleting old backup:", oldestBackup)
			err := os.RemoveAll(oldestBackup)
			if err != nil {
				fmt.Println("Error removing old backup:", err)
			}
		}
	}
}
