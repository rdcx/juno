package service

import (
	"juno/pkg/node/page"
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	data := []byte("test data")

	vHash := page.NewVersionHash(data)
	err := s.Write(vHash, data)
	if err != nil {
		t.Error(err)
	}

	dataWritten, err := os.ReadFile(dir + "/" + vHash.String())
	if err != nil {
		t.Error(err)
	}

	if string(dataWritten) != string(data) {
		t.Errorf("expected %s, got %s", string(data), string(dataWritten))
	}
}
