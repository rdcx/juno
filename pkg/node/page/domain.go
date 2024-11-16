package page

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"juno/pkg/shard"
	"time"
)

var ErrPageNotFound = errors.New("page not found")
var ErrPageAlreadyExists = errors.New("page already exists")

type PageID [16]byte

func (id PageID) String() string {
	return hex.EncodeToString(id[:])
}

func NewPageID(u string) PageID {
	hash := sha256.New()
	hash.Write([]byte(u))
	fullHash := hash.Sum(nil)
	var id PageID
	copy(id[:], fullHash[:16])
	return id
}

type VersionHash [16]byte

func (h VersionHash) String() string {
	return hex.EncodeToString(h[:])
}

func NewVersionHash(data []byte) VersionHash {
	hash := sha256.New()
	hash.Write(data)
	fullHash := hash.Sum(nil)
	var h VersionHash
	copy(h[:], fullHash[:16])
	return h
}

type Version struct {
	Hash      VersionHash `json:"hash"`
	CreatedAt time.Time   `json:"created_at"`
}

func NewVersion(hash VersionHash) Version {
	return Version{
		Hash:      hash,
		CreatedAt: time.Now(),
	}
}

type Page struct {
	ID       PageID    `json:"id"`
	Shard    int       `json:"shard"`
	URL      string    `json:"url"`
	Versions []Version `json:"versions"`
}

func NewPage(url string) *Page {
	return &Page{
		ID:    NewPageID(url),
		URL:   url,
		Shard: shard.GetShard(url),
	}
}

type Repository interface {
	CreatePage(page *Page) error
	GetPage(id PageID) (*Page, error)
	AddVersion(pageID PageID, version Version) error
	GetVersions(pageID PageID) ([]Version, error)
	Iterator(fn func(*Page)) error
	Count() (int, error)
}

type Service interface {
	Create(page *Page) error
	Get(pageID PageID) (*Page, error)
	GetByURL(url string) (*Page, error)
	AddVersion(pageID PageID, version Version) error
	GetVersions(pageID PageID) ([]Version, error)
	Iterator(fn func(*Page)) error
	Count() (int, error)
}
