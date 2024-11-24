package engin

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const prefix = "rscript"

type internalFunc func(*ssh.Client, ...string) (string, error)

var internalFuncMap = map[string]internalFunc{
	prefix + ".uploadFile":   UploadFileWithFs(nil),
	prefix + ".downloadFile": downloadFile,
	prefix + ".exec":         execLocalCommand,
}

func SetInternalFunc(name string, f internalFunc) {
	internalFuncMap[name] = f
}

// remoteFilePath, localPath string
func downloadFile(client *ssh.Client, args ...string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("args error")
	}
	var (
		remoteFilePath = args[0]
		localFilePath  = args[1]
	)
	if remoteFilePath == "" {
		return "", errors.New("remote file path is empty")
	}
	if localFilePath == "" {
		return "", errors.New("local file path is empty")
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return "", fmt.Errorf("sftp.NewClient: %w", err)
	}
	defer sftpClient.Close()

	remoteFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		return "", fmt.Errorf("sftpClient.Open: %w", err)
	}
	defer remoteFile.Close()

	localFile, err := os.Create(localFilePath)
	if err != nil {
		return "", fmt.Errorf("os.Create: %w", err)
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}

	return "", nil
}

func UploadFileWithFs(f fs.FS) internalFunc {
	return func(client *ssh.Client, args ...string) (string, error) {
		if len(args) != 2 {
			return "", errors.New("args error")
		}
		var (
			localFilePath  = args[0]
			remoteFilePath = args[1]
		)
		if localFilePath == "" {
			return "", errors.New("local file path is empty")
		}
		if remoteFilePath == "" {
			return "", errors.New("remote file path is empty")
		}

		sftpClient, err := sftp.NewClient(client)
		if err != nil {
			return "", fmt.Errorf("sftp.NewClient: %w", err)
		}
		defer sftpClient.Close()

		var localFile fs.File

		if f == nil {
			localFile, err = os.Open(localFilePath)
			if err != nil {
				return "", fmt.Errorf("os.Open: %w", err)
			}
		} else {
			localFile, err = f.Open(localFilePath)
			if err != nil {
				return "", fmt.Errorf("f.Open: %w", err)
			}
		}

		defer localFile.Close()

		remoteFile, err := sftpClient.Create(remoteFilePath)
		if err != nil {
			return "", fmt.Errorf("sftpClient.Create: %w", err)
		}
		defer remoteFile.Close()

		_, err = io.Copy(remoteFile, localFile)
		if err != nil {
			return "", fmt.Errorf("io.Copy: %w", err)
		}

		return "", nil
	}
}

func execLocalCommand(_ *ssh.Client, args ...string) (string, error) {
	var (
		name    string
		cmdArgs []string
	)

	switch len(args) {
	case 0:
		return "", errors.New("args error")
	case 1:
		name = args[0]
	default:
		name = args[0]
		cmdArgs = append(cmdArgs, args[1:]...)
	}

	cmd := exec.Command(name, cmdArgs...)
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("cmd.Run: %w", err)
	}
	return out.String(), nil
}
