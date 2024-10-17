//go:build !solution

package keylock

import (
	"sort"
	"sync"
)

type KeyLock struct {
	locks map[string]chan struct{}
	cond  *sync.Cond
}

func New() *KeyLock {
	return &KeyLock{
		locks: make(map[string]chan struct{}),
		cond:  sync.NewCond(&sync.Mutex{}),
	}
}

func (kl *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	sortedKeys := kl.SortKeys(keys)
	for i, key := range sortedKeys {
		kl.cond.L.Lock()
		lockChan := kl.GetOrCreateLock(key)
		kl.cond.L.Unlock()
		select {
		case <-cancel:
			canceled = true
			unlock = nil
			kl.UnlockKeys(sortedKeys[:i])
			return
		case lockChan <- struct{}{}:
		}
	}
	canceled = false
	unlock = func() {
		kl.UnlockKeys(sortedKeys)
	}
	return
}

func (kl *KeyLock) SortKeys(keys []string) []string {
	sortedKeys := make([]string, len(keys))
	copy(sortedKeys, keys)
	sort.Strings(sortedKeys)
	return sortedKeys
}

func (kl *KeyLock) GetOrCreateLock(key string) chan struct{} {
	if kl.locks[key] == nil {
		kl.locks[key] = make(chan struct{}, 1)
	}
	return kl.locks[key]
}

func (kl *KeyLock) UnlockKeys(keys []string) {
	for _, key := range keys {
		kl.cond.L.Lock()
		<-kl.locks[key]
		kl.cond.L.Unlock()
	}
}
