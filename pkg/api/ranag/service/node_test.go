package service

import (
	"juno/pkg/api/ranag"
	"juno/pkg/api/ranag/repo/mem"
	"juno/pkg/api/user"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func testRanagMatches(t *testing.T, id, ownerID uuid.UUID, address string, shardAssignments [][2]int, n *ranag.Ranag) bool {
	if n.ID != id {
		t.Errorf("Expected ID %s, got %s", id, n.ID)
		return false
	}

	if n.OwnerID != ownerID {
		t.Errorf("Expected OwnerID %s, got %s", ownerID, n.OwnerID)
		return false
	}

	if n.Address != address {
		t.Errorf("Expected Address %s, got %s", address, n.Address)
		return false
	}

	if len(n.ShardAssignments) != len(shardAssignments) {
		t.Errorf("Expected %d shard assignments, got %d", len(shardAssignments), len(n.ShardAssignments))
		return false
	}

	for i, s := range shardAssignments {
		if n.ShardAssignments[i][0] != s[0] || n.ShardAssignments[i][1] != s[1] {
			t.Errorf("Expected shard assignment %v, got %v", s, n.ShardAssignments[i])
			return false
		}
	}

	return true
}

func TestAllShardsRanags(t *testing.T) {
	repo := mem.New()
	svc := New(repo)

	n1 := &ranag.Ranag{
		ID:               uuid.New(),
		OwnerID:          uuid.New(),
		Address:          "example.com",
		ShardAssignments: [][2]int{{0, 1000}, {1000, 1000}},
	}

	err := repo.Create(n1)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	n2 := &ranag.Ranag{
		ID:               uuid.New(),
		OwnerID:          uuid.New(),
		Address:          "example.org",
		ShardAssignments: [][2]int{{2000, 1000}, {3000, 2000}},
	}

	err = repo.Create(n2)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	shards, err := svc.AllShardsRanags()

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if len(shards) != 5000 {
		t.Errorf("Expected 5000 shards, got %d", len(shards))
	}

	if len(shards[0]) != 1 {
		t.Errorf("Expected 1 ranag for shard 0, got %d", len(shards[0]))
	}

	if len(shards[1000]) != 1 {
		t.Errorf("Expected 1 ranag for shard 1000, got %d", len(shards[1000]))
	}

	for i := 1; i < 1000; i++ {
		if shards[i][0].ID != n1.ID {
			t.Errorf("Expected ranag %s for shard %d, got %s", n1.ID, i, shards[i][0].ID)
		}

		if shards[i+1000][0].ID != n1.ID {
			t.Errorf("Expected ranag %s for shard %d, got %s", n1.ID, i+1000, shards[i+1000][0].ID)
		}

		if shards[i+2000][0].ID != n2.ID {
			t.Errorf("Expected ranag %s for shard %d, got %s", n2.ID, i+2000, shards[i+2000][0].ID)
		}

		if shards[i+3000][0].ID != n2.ID {
			t.Errorf("Expected ranag %s for shard %d, got %s", n2.ID, i+3000, shards[i+3000][0].ID)
		}

		if shards[i+4000][0].ID != n2.ID {
			t.Errorf("Expected ranag %s for shard %d, got %s", n2.ID, i+4000, shards[i+4000][0].ID)
		}
	}
}

func TestValidateShardAssignments(t *testing.T) {
	tests := []struct {
		name    string
		shards  [][2]int
		wantErr bool
	}{
		{
			name:    "valid",
			shards:  [][2]int{{0, 1}, {1, 2}},
			wantErr: false,
		},
		{
			name:   "valid 1k",
			shards: [][2]int{{0, 1000}, {1000, 1000}},
		},
		{
			name:    "invalid",
			shards:  [][2]int{{1, 100_001}, {1000, 1000}},
			wantErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			err := validateShardAssignments(tt.shards)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %s", err)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		u := &user.User{
			ID: uuid.New(),
		}

		addr := "example.com:7000"

		n, err := svc.Create(u.ID, addr, [][2]int{{0, 1}, {1, 2}})

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		ranag, err := repo.Get(n.ID)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testRanagMatches(t, n.ID, n.OwnerID, addr, [][2]int{{0, 1}, {1, 2}}, ranag) {
			t.Errorf("Ranag does not match")
		}
	})

	t.Run("validation", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		u := &user.User{
			ID: uuid.New(),
		}

		addr := "bad address"

		n, err := svc.Create(u.ID, addr, [][2]int{{1, 100_001}, {1000, 1000}})

		if err == nil {
			t.Fatal("Expected an error")
		}

		if !strings.Contains(err.Error(), "invalid address") {
			t.Errorf("Expected ErrInvalidAddress, got %v", err)
		}

		if !strings.Contains(err.Error(), "invalid shards") {
			t.Errorf("Expected ErrInvalidShards, got %v", err)
		}

		if n != nil {
			t.Errorf("Expected nil, got %v", n)
		}
	})

	t.Run("error address already exists", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &ranag.Ranag{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
		}

		u := &user.User{
			ID: n.OwnerID,
		}

		err := repo.Create(n)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		addr := "http://example.com"

		n, err = svc.Create(u.ID, addr, [][2]int{{0, 1}, {1, 2}})

		if err != ranag.ErrAddressExists {
			t.Errorf("Expected error, got nil")
		}

		if n != nil {
			t.Errorf("Expected nil, got %v", n)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &ranag.Ranag{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "example.com:8000",
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		ranag, err := svc.Get(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testRanagMatches(t, n.ID, n.OwnerID, n.Address, n.ShardAssignments, ranag) {
			t.Errorf("Ranag does not match")
		}
	})

	t.Run("error not found", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &ranag.Ranag{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
		}

		_, err := svc.Get(n.ID)
		if err != ranag.ErrNotFound {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestListByOwnerID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		ownerID := uuid.New()

		n1 := &ranag.Ranag{
			ID:               uuid.New(),
			OwnerID:          ownerID,
			Address:          "example.com:8000",
			ShardAssignments: [][2]int{{0, 1}, {1, 2}},
		}

		err := repo.Create(n1)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		n2 := &ranag.Ranag{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "example2.com:8080",
		}

		err = repo.Create(n2)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		ranags, err := svc.ListByOwnerID(ownerID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if len(ranags) != 1 {
			t.Errorf("Expected 1 ranags, got %d", len(ranags))
		}

		if !testRanagMatches(t, n1.ID, n1.OwnerID, n1.Address, n1.ShardAssignments, ranags[0]) {
			t.Errorf("Ranag does not match")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &ranag.Ranag{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "ranag.com:9392",
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		n.Address = "example2.com:8080"

		n, err = svc.Update(n.ID, n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		ranag, err := repo.Get(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testRanagMatches(t, n.ID, n.OwnerID, n.Address, n.ShardAssignments, ranag) {
			t.Errorf("Ranag does not match")
		}
	})

	t.Run("validation", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &ranag.Ranag{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "valid.com:8000",
		}

		err := repo.Create(n)

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		nn := &ranag.Ranag{
			ID:      n.ID,
			OwnerID: n.OwnerID,
		}

		nn.Address = "bad address"

		updated, err := svc.Update(n.ID, nn)

		if err == nil {
			t.Fatal("Expected an error")
		}

		if updated != nil {
			t.Fatal("Expected ranag, got nil")
		}

		if !strings.Contains(err.Error(), "invalid address") {
			t.Errorf("Expected ErrInvalidAddress, got %v", err)
		}

		check, err := repo.Get(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if !testRanagMatches(t, n.ID, n.OwnerID, "valid.com:8000", n.ShardAssignments, check) {
			t.Errorf("Ranag does not match")
		}
	})

	t.Run("error not found", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &ranag.Ranag{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
		}

		n, err := svc.Update(uuid.New(), n)
		if err != ranag.ErrNotFound {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n := &ranag.Ranag{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			Address: "http://example.com",
		}

		err := repo.Create(n)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		err = svc.Delete(n.ID)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		_, err = repo.Get(n.ID)
		if err != ranag.ErrNotFound {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("error not found", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		err := svc.Delete(uuid.New())
		if err != ranag.ErrNotFound {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestGroupByRange(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mem.New()
		svc := New(repo)

		n1 := &ranag.Ranag{
			ID:               uuid.New(),
			OwnerID:          uuid.New(),
			Address:          "example.com:8000",
			ShardAssignments: [][2]int{{0, 10}, {10, 10}},
		}

		err := repo.Create(n1)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		n2 := &ranag.Ranag{
			ID:               uuid.New(),
			OwnerID:          uuid.New(),
			Address:          "example2.com:8080",
			ShardAssignments: [][2]int{{20, 10}, {30, 10}},
		}

		err = repo.Create(n2)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		groups, err := svc.GroupByRange()

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if len(groups) != 4 {
			t.Errorf("Expected 4 groups, got %d", len(groups))
		}

		if len(groups[[2]int{0, 10}]) != 1 {
			t.Errorf("Expected 1 ranag for shard range 0-10, got %d", len(groups[[2]int{0, 10}]))
		}

		if len(groups[[2]int{10, 10}]) != 1 {
			t.Errorf("Expected 1 ranag for shard range 10-10, got %d", len(groups[[2]int{10, 10}]))
		}

		if len(groups[[2]int{20, 10}]) != 1 {
			t.Errorf("Expected 1 ranag for shard range 20-10, got %d", len(groups[[2]int{20, 10}]))
		}

		if len(groups[[2]int{30, 10}]) != 1 {
			t.Errorf("Expected 1 ranag for shard range 30-10, got %d", len(groups[[2]int{30, 10}]))
		}
	})
}
