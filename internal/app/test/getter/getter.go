package getter

import "os"

type Getter interface {
	DownloadFile(url string, file *os.File) error
}
