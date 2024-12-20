package service

import (
	"juno/pkg/node/page"
	"os"
)

type Service struct {
	dir string
}

func New(dir string) *Service {

	// create if not exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	return &Service{
		dir: dir,
	}
}

func (s *Service) Write(hash page.VersionHash, data []byte) error {
	return os.WriteFile(s.dir+"/"+hash.String(), data, 0644)
}

func (s *Service) Read(hash page.VersionHash) ([]byte, error) {
	return os.ReadFile(s.dir + "/" + hash.String())
}
