# PZBackup

PZBackup is a simple backup utility for Project Zomboid game saves. It allows users to automatically back up their game saves to a specified location at defined intervals.

## Table of Contents

- [Usage](#usage)
- [Configuration](#configuration)
- [How It Works](#how-it-works)
- [License](#license)

## Usage

Download run the binary from [release](https://github.com/angga2oioi/project-zomboid-auto-backup/releases) page. If no configuration file is found, the application will prompt you to enter:

1. Project Zomboid save folder path
2. Backup folder path
3. Backup interval in minutes

## Configuration

The application uses a JSON configuration file to store the backup settings. The configuration is automatically created in the `config` directory upon first run if it does not exist.

### Example Configuration Structure

```json
{
  "save_folder": "/path/to/save_folder",
  "backup_folder": "/path/to/backup_folder",
  "backup_interval": 30
}
```

## How It Works

Once started, PZBackup will start the backup process based on the specified interval. It performs the following actions:

1. Clearance of the backup folder.
2. Copies all files from the save folder to the backup folder.
3. Waits for the defined interval before repeating the process.

The application is capable of working on multiple platforms (Windows, Linux, macOS) and uses the appropriate command for file copying based on the operating system.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
