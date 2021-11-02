package filesaver

import (
	"2021_2_MAMBa/internal/pkg/utils/log"
	"2021_2_MAMBa/internal/pkg/utils/randomizer"
	"fmt"
	"io"
	"os"
)

func createFile(root, dir, name string) (*os.File, error) {
	_, err := os.ReadDir(root + dir)
	if err != nil {
		err = os.MkdirAll(root+dir, 0777)
		if err != nil {
			return nil, err
		}
	}
	file, err := os.Create(root + dir + name)
	return file, err
}

func UploadFile(reader io.Reader, root, path, ext string) (string, error) {
	randString, err := randomizer.GenerateRandomString(6)
	if err != nil {
		return "", err
	}
	filename := randString + ext
	log.Info("Created file with name " + filename)
	file, err := createFile(root, path, filename)
	if err != nil {
		return "", fmt.Errorf("file creating error: %s", err)
	}
	defer file.Close()

	filename = path + filename
	_, err = io.Copy(file, reader)
	if err != nil {
		return "", fmt.Errorf("copy error: %s", err)
	}
	return filename, nil
}
