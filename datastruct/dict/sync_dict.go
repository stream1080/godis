package dict

import (
	"sync"
)

type SyncDict struct {
	m sync.Map
}

func MakeSyncDict() *SyncDict {
	return &SyncDict{}
}

func (dict *SyncDict) Get(key string) (val interface{}, exists bool) {
	return dict.m.Load(key)
}

func (dict *SyncDict) Len() int {
	length := 0
	dict.m.Range(func(key, value any) bool {
		length++
		return true
	})

	return length
}

func (dict *SyncDict) Put(key string, val interface{}) (result int) {
	_, existed := dict.m.Load(key)
	dict.m.Store(key, val)
	if existed {
		return 0
	}

	return 1
}

func (dict *SyncDict) PutIfAbsent(key string, val interface{}) (result int) {
	_, existed := dict.m.Load(key)
	if existed {
		return 0
	}
	dict.m.Store(key, val)
	return 1
}

func (dict *SyncDict) PutIfExists(key string, val interface{}) (result int) {
	_, existed := dict.m.Load(key)
	if existed {
		dict.m.Store(key, val)
		return 1
	}

	return 0
}

func (dict *SyncDict) Remove(key string) (result int) {
	_, existed := dict.m.Load(key)
	dict.m.Delete(key)
	if existed {
		return 1
	}

	return 0
}

func (dict *SyncDict) ForEach(consumer Consumer) {
	dict.m.Range(func(key, value any) bool {
		consumer(key.(string), value)
		return true
	})
}

func (dict *SyncDict) Keys() []string {
	result := make([]string, dict.Len())
	i := 0
	dict.m.Range(func(key, value any) bool {
		result[i] = key.(string)
		i++
		return true
	})

	return result
}

func (dict *SyncDict) RandomKeys(limit int) []string {
	result := make([]string, dict.Len())
	for i := 0; i < limit; i++ {
		dict.m.Range(func(key, value any) bool {
			result[i] = key.(string)
			return false
		})
	}

	return result
}

func (dict *SyncDict) RandomDistinctKeys(limit int) []string {
	result := make([]string, dict.Len())
	i := 0

	dict.m.Range(func(key, value any) bool {
		result[i] = key.(string)
		i++
		return i != limit
	})

	return result
}

func (dict *SyncDict) Clear() {
	*dict = *MakeSyncDict()
}
