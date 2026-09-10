package media

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeStorage struct {
	puts    map[string][]byte
	deleted []string
}

func (s *fakeStorage) Put(_ context.Context, path string, content io.Reader) error {
	data, err := io.ReadAll(content)
	if err != nil {
		return err
	}
	if s.puts == nil {
		s.puts = make(map[string][]byte)
	}
	s.puts[path] = data
	return nil
}

func (s *fakeStorage) Delete(_ context.Context, path string) error {
	s.deleted = append(s.deleted, path)
	delete(s.puts, path)
	return nil
}

func (s *fakeStorage) Path(path string) string { return "/tmp/media/" + path }

func (s *fakeStorage) URL(path string) string { return "/api/media/preview/" + path }

func TestServiceUploadsMetadataAndNormalizesFileExtension(t *testing.T) {
	storage := &fakeStorage{}
	service := NewService(NewMemoryRepository(), storage)

	item, err := service.Upload(context.Background(), UploadInput{
		OriginalName: "avatar.PNG", MIMEType: "image/png", Size: 7,
		Content: bytes.NewReader([]byte("png-data")),
	})
	require.NoError(t, err)
	require.Equal(t, "avatar.PNG", item.OriginalName)
	require.Equal(t, "image/png", item.MIMEType)
	require.Equal(t, int64(7), item.Size)
	require.Contains(t, item.Path, ".png")
	require.Equal(t, []byte("png-data"), storage.puts[item.Path])
}

func TestServiceRejectsUnsafeAndOversizedUploads(t *testing.T) {
	service := NewService(NewMemoryRepository(), &fakeStorage{})

	_, err := service.Upload(context.Background(), UploadInput{OriginalName: "../../secret.txt", MIMEType: "text/plain", Size: 4, Content: bytes.NewReader([]byte("safe"))})
	require.ErrorIs(t, err, ErrInvalidMedia)

	_, err = service.Upload(context.Background(), UploadInput{OriginalName: "large.txt", MIMEType: "text/plain", Size: MaxUploadSize + 1, Content: bytes.NewReader([]byte("safe"))})
	require.ErrorIs(t, err, ErrMediaTooLarge)
}

func TestServiceDeletesBlobBeforeMetadata(t *testing.T) {
	storage := &fakeStorage{}
	service := NewService(NewMemoryRepository(), storage)
	item, err := service.Upload(context.Background(), UploadInput{OriginalName: "file.txt", MIMEType: "text/plain", Size: 4, Content: bytes.NewReader([]byte("data"))})
	require.NoError(t, err)

	require.NoError(t, service.Delete(context.Background(), item.ID))
	require.Equal(t, []string{item.Path}, storage.deleted)
	_, err = service.Get(context.Background(), item.ID)
	require.ErrorIs(t, err, ErrNotFound)
}
