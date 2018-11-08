package filesystem

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/logger"
)

func GetHashedNameForFile(uID uint, filename string) string {
	hasher := sha256.New()
	hasher.Write([]byte(time.Now().String() + fmt.Sprintf("%v", uID) + filename))
	hash := hex.EncodeToString(hasher.Sum(nil))

	return hash[:16] + path.Ext(filename)
}

func SaveFile(file io.Reader, dir, filename string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	f, err := os.OpenFile(dir+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, file)
	logger.Infow("saved file",
		"path", dir,
		"filename", filename)

	return nil
}
