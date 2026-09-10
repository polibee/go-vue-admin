package media

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/support/path"
)

type LocalStorage struct {
	root string
	base string
}

func NewConfiguredLocalStorage() *LocalStorage {
	return NewLocalStorage(path.Storage("app/media"), "/api/media/preview")
}

func NewLocalStorage(root, baseURL string) *LocalStorage {
	return &LocalStorage{root: filepath.Clean(root), base: strings.TrimRight(baseURL, "/")}
}

func (s *LocalStorage) Put(_ context.Context, name string, content io.Reader) error {
	full, err := s.safePath(name)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(full, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o640)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, content)
	return err
}

func (s *LocalStorage) Delete(_ context.Context, name string) error {
	full, err := s.safePath(name)
	if err != nil {
		return err
	}
	if err := os.Remove(full); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *LocalStorage) Path(name string) string {
	full, err := s.safePath(name)
	if err != nil {
		return ""
	}
	return full
}

func (s *LocalStorage) URL(name string) string {
	if _, err := s.safePath(name); err != nil {
		return ""
	}
	return s.base + "/" + name
}

func (s *LocalStorage) safePath(name string) (string, error) {
	name = filepath.Clean(strings.TrimSpace(name))
	if name == "." || filepath.IsAbs(name) || name == ".." || strings.HasPrefix(name, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("invalid media path")
	}
	return filepath.Join(s.root, name), nil
}
