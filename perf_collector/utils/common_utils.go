package utils

import "os"

func DirExists(path string) bool {
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return fi.IsDir()
}

func CreateDir(path string) error {
	if !DirExists(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveToFile(output []byte, filename string) error {
	return os.WriteFile(filename, output, 0644)
}
