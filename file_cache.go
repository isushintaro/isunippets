package isunippets

import (
	"fmt"
	"io"
	"os"
	"path"
)

var (
	fileCacheRootDir = path.Join(os.TempDir(), "isunippets-fileCache")
	fileCacheBaseDir = path.Join(fileCacheRootDir, fmt.Sprintf("%d", os.Getpid()))
)

func CleanupFileCache() error {
	if _, err := os.Stat(fileCacheRootDir); os.IsNotExist(err) {
		return nil
	}
	return os.RemoveAll(fileCacheRootDir)
}

func PutFileCacheData(data []byte, fileName string) error {
	filePath := path.Join(fileCacheBaseDir, fileName)
	baseDir := path.Dir(filePath)

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
			return err
		}
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	if _, err = f.Write(data); err != nil {
		return err
	}

	return nil
}

func PutFileCacheStream(stream io.Reader, fileName string) error {
	baseDir := fileCacheBaseDir
	filePath := path.Join(baseDir, fileName)

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
			return err
		}
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	if _, err = io.Copy(f, stream); err != nil {
		return err
	}

	return nil
}

func GetFileCacheData(fileName string) ([]byte, error) {
	filePath := path.Join(fileCacheBaseDir, fileName)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetFileCacheStream(fileName string) (io.Reader, error) {
	filePath := path.Join(fileCacheBaseDir, fileName)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return f, nil
}
