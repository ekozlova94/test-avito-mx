package testgetter

import (
	"io"
	"os"
)

type GetterImpl struct {
}

func NewGetter() *GetterImpl {
	return &GetterImpl{}
}

func (s *GetterImpl) DownloadFile(url string, file *os.File) error {
	xlsxFile, err := os.Open(url)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, xlsxFile)
	if err != nil {
		return err
	}
	return nil
}
