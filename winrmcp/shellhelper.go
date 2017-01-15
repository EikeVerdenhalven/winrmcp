package winrmcp

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/masterzen/winrm"
)

// SendShellCommand sends commandString on the passed in WinRM shell
func SendShellCommand(shell *winrm.Shell, commandString string) error {
	cmd, err := shell.Execute(commandString)
	if err != nil {
		return err
	}
	defer cmd.Close()

	var wg sync.WaitGroup
	copyFunc := func(w io.Writer, r io.Reader) {
		defer wg.Done()
		io.Copy(w, r)
	}

	wg.Add(2)
	go copyFunc(os.Stdout, cmd.Stdout)
	go copyFunc(os.Stderr, cmd.Stderr)

	cmd.Wait()
	wg.Wait()

	if cmd.ExitCode() != 0 {
		return fmt.Errorf("command returned code=%d", cmd.ExitCode())
	}
	return nil
}

// ExecuteRemoteCommand runs commandString on the WinRM remote client
func ExecuteRemoteCommand(client *winrm.Client, commandString string) error {
	shell, err := client.CreateShell()
	if err != nil {
		return err
	}
	defer shell.Close()
	return SendShellCommand(shell, commandString)
}
