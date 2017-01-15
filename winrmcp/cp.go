package winrmcp

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/masterzen/winrm"
	"github.com/nu7hatch/gouuid"
)

func doCopy(client *winrm.Client, config *Config, in io.Reader, toPath string) error {
	tempFile, err := tempFileName()
	if err != nil {
		return fmt.Errorf("Error generating unique filename: %v", err)
	}
	tempPath := "$env:TEMP\\" + tempFile

	debugLog(fmt.Sprintf("Copying file to %s\n", tempPath))

	err = uploadContent(client, config.MaxOperationsPerShell, "%TEMP%\\"+tempFile, in)
	if err != nil {
		return fmt.Errorf("Error uploading file to %s: %v", tempPath, err)
	}

	debugLog(fmt.Sprintf("Moving file from %s to %s", tempPath, toPath))

	err = restoreContent(client, tempPath, toPath)
	if err != nil {
		return fmt.Errorf("Error restoring file from %s to %s: %v", tempPath, toPath, err)
	}

	debugLog(fmt.Sprintf("Removing temporary file %s", tempPath))

	err = cleanupContent(client, tempPath)
	if err != nil {
		return fmt.Errorf("Error removing temporary file %s: %v", tempPath, err)
	}

	return nil
}

func chunkSize(filePathLength int) int {
	// Upload the file in chunks to get around the Windows command line size limit.
	// Base64 encodes each set of three bytes into four bytes. In addition the output
	// is padded to always be a multiple of four.
	//
	//   ceil(n / 3) * 4 = m1 - m2
	//
	//   where:
	//     n  = bytes
	//     m1 = max (8192 character command limit.)
	//     m2 = len(filePath)
	return ((8000 - filePathLength) / 4) * 3
}

func uploadContent(client *winrm.Client, maxChunks int, filePath string, reader io.Reader) error {
	var err error
	done := false
	for !done {
		done, err = uploadChunks(client, filePath, maxChunks, reader)
		if err != nil {
			return err
		}
	}

	return nil
}

func uploadChunks(client *winrm.Client, filePath string, maxChunks int, reader io.Reader) (bool, error) {
	shell, err := client.CreateShell()
	if err != nil {
		return false, fmt.Errorf("Couldn't create shell: %v", err)
	}
	defer shell.Close()

	chunk := make([]byte, chunkSize(len(filePath)))

	if maxChunks == 0 {
		maxChunks = 1
	}

	for i := 0; i < maxChunks; i++ {
		n, err := reader.Read(chunk)

		if err != nil && err != io.EOF {
			return false, err
		}
		if n == 0 {
			return true, nil
		}

		content := base64.StdEncoding.EncodeToString(chunk[:n])
		if err = appendContent(shell, filePath, content); err != nil {
			return false, err
		}
	}

	return false, nil
}

func restoreContent(client *winrm.Client, fromPath, toPath string) error {
	script := fmt.Sprintf(`
		$tmp_file_path = [System.IO.Path]::GetFullPath("%s")
		$dest_file_path = [System.IO.Path]::GetFullPath("%s".Trim("'"))
		if (Test-Path $dest_file_path) {
			rm $dest_file_path
		}
		else {
			$dest_dir = ([System.IO.Path]::GetDirectoryName($dest_file_path))
			New-Item -ItemType directory -Force -ErrorAction SilentlyContinue -Path $dest_dir | Out-Null
		}

		if (Test-Path $tmp_file_path) {
			$base64_lines = Get-Content $tmp_file_path
			$base64_string = [string]::join("",$base64_lines)
			$bytes = [System.Convert]::FromBase64String($base64_string)
			[System.IO.File]::WriteAllBytes($dest_file_path, $bytes)
		} else {
			echo $null > $dest_file_path
		}
	`, fromPath, toPath)

	return ExecuteRemoteCommand(client, winrm.Powershell(script))
}

func cleanupContent(client *winrm.Client, filePath string) error {
	return ExecuteRemoteCommand(client, winrm.Powershell(fmt.Sprintf("Remove-Item %s -ErrorAction SilentlyContinue", filePath)))
}

func appendContent(shell *winrm.Shell, filePath, content string) error {
	return SendShellCommand(shell, fmt.Sprintf("echo %s >> \"%s\"", content, filePath))
}

func tempFileName() (string, error) {
	uniquePart, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("winrmcp-%s.tmp", uniquePart), nil
}
