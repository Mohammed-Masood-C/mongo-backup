# MongoDB Backup & Restore CLI

A simple, lightweight command-line tool for backing up and restoring MongoDB databases.

---

## 🚀 Features

* Backup MongoDB databases to local files
* Restore databases from backups

---

## 📦 Installation

Download the latest release from GitHub and place the binary somewhere on your system.

---

## ⚙️ First Run Setup

Run the binary once:

```bash
backup.exe
```

This will automatically create a configuration folder at:

```text
%APPDATA%\mongo-backup
```

---

## 🧩 Configuration

Inside the `mongo-backup` folder, create a JSON file for each database you want to manage.

### 📁 File Naming

It is reccomended for your configuration files to be named like:

```text
databasename.json
```

In the case of tracking same database name across different mongo clusters, use names like below:
```text
cluster1-databasename.json
```

### 📝 Example Configuration

```json
{
  "databaseName": "<database name>",
  "uri": "<mongo db uri where database is present>",
  "backupFolderPath": "<folder path to store backups>"
}
```

You can add multiple files to track multiple databases.

---

## ▶️ Usage

After adding your configuration files, run:

```bash
backup.exe
```

---

## 🧭 Notes

* Ensure your MongoDB URI is valid and accessible
* Make sure the backup folder paths exist or are writable

---
