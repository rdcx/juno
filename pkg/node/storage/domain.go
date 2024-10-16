package storage

import "juno/pkg/node/page"

type Service interface {
	Write(hash page.VersionHash, data []byte) error
}
