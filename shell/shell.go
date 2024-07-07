package shell

import (
	"bytes"
	"os/exec"
	"os"
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
)

func Shellout(command string) (string, string, error) {
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd := exec.Command("bash", "-c", command)
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()
    return stdout.String(), stderr.String(), err
}

func ZipDirectory(directoryName string, archiveName string) () {
	file, err := os.Create(archiveName)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    w := zip.NewWriter(file)
    defer w.Close()

    walker := func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()

        // Ensure that `path` is not absolute; it should not start with "/".
        // This snippet happens to work because I don't use 
        // absolute paths, but ensure your real-world code 
        // transforms path into a zip-root relative path.
        f, err := w.Create(path)
        if err != nil {
            return err
        }

        _, err = io.Copy(f, file)
        if err != nil {
            return err
        }

        return nil
    }
    err = filepath.Walk(directoryName, walker)
    if err != nil {
        panic(err)
    }
	fmt.Println("Shell: Archive created.")
}
