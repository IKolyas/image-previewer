package memory

import (
	"context"
	"sync"

	"slices"

	"github.com/IKolyas/image-previewer/internal/core/image"
	"github.com/IKolyas/image-previewer/internal/storage/source"
)

type LRUStorage struct {
	capacity int
	cache    map[string][]byte
	order    []string
	mu       sync.Mutex
}

func NewLRUStorage(capacity int) *LRUStorage {
	return &LRUStorage{
		capacity: capacity,
		cache:    make(map[string][]byte),
		order:    make([]string, 0, capacity),
	}
}

func (s *LRUStorage) Get(ctx context.Context, imgData *image.ImgData) ([]byte, error) {
	key := imgData.String()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check cache first
	if imgData, ok := s.cache[key]; ok {
		// Move to front of order
		s.moveToFront(key)
		return imgData, nil
	}

	// Download and process image
	data, err := source.Get(ctx, imgData)
	if err != nil {
		return nil, err
	}

	// Add to cache
	s.addToCache(key, data)

	return data, nil
}

func (s *LRUStorage) addToCache(key string, imgData []byte) {
	if len(s.order) >= s.capacity {
		// Remove least recently used
		oldest := s.order[len(s.order)-1]
		delete(s.cache, oldest)
		s.order = s.order[:len(s.order)-1]
	}

	s.cache[key] = imgData
	s.order = append([]string{key}, s.order...)
}

func (s *LRUStorage) moveToFront(key string) {
	for i, k := range s.order {
		if k == key {
			s.order = slices.Delete(s.order, i, i+1)
			s.order = append([]string{key}, s.order...)
			break
		}
	}
}

func (s *LRUStorage) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache = make(map[string][]byte)
	s.order = make([]string, 0, s.capacity)
}
