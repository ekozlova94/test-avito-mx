package prodgetter

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type GetterImpl struct {
	log *zap.Logger
}

func NewGetter(log *zap.Logger) *GetterImpl {
	return &GetterImpl{
		log: log,
	}
}

func (s *GetterImpl) DownloadFile(url string, file *os.File) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	//noinspection GoUnhandledErrorResult
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("received a response code not 200")
	}

	written, err := io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	s.log.Info("File downloaded successfully", zap.Int64("written byte", written))
	return nil
}
