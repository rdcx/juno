package crawl

type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Work() {}
