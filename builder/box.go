package builder

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"bytes"
	"compress/gzip"
)

type box struct {
	Name  string
	Files []file
}

func (b *box) Walk(root string, compress bool) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		name := strings.Replace(path, root+string(os.PathSeparator), "", 1)
		name = strings.Replace(name, "\\", "/", -1)
		f := file{
			Name: name,
		}

		bb, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.WithStack(err)
		}
		if compress {
			bb, err = compressFile(bb)
			if err != nil {
				return errors.WithStack(err)
			}
		}
		bb, err = json.Marshal(bb)
		if err != nil {
			return errors.WithStack(err)
		}
		f.Contents = strings.Replace(string(bb), "\"", "\\\"", -1)

		b.Files = append(b.Files, f)
		return nil
	})
}

func compressFile(bb []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(bb)
	if err != nil {
		return bb, errors.WithStack(err)
	}
	err = writer.Close()
	if err != nil {
		return bb, errors.WithStack(err)
	}
	return buf.Bytes(), nil
}