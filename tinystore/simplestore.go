package tinystore

import (
	"sync"
	"errors"
	"io/ioutil"
	"encoding/json"
)

// SimpleStore implements  Store
type SimpleStore struct {
	mutex sync.Mutex
	items []*Credential
}

// All implements Store.All
func (s *SimpleStore) All() []*Credential {
	return s.items
}

// Length implements Store.Length
func(store *SimpleStore) Length() (int, error ) {
	if store ==nil {
		return 0 , errors.New("Nil Items")

	}
	return len(store.items), nil
}

// Find implements Store.Find
func (s *SimpleStore) Find(f Filter) (*Credential, error) {
	for _, x := range s.items {
		if f(x) {
			return x, nil
		}
	}
	return &Credential{}, ErrNotFound
}

// Clear Implements Store.Clear
func (s *SimpleStore) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.items = make([]*Credential, 0)
}

// Remove implements Store.Remove
func (s *SimpleStore) Remove(c *Credential) error {

	if ex:= c.Validate() ; ex!= nil { return ex }

	s.mutex.Lock()
	defer s.mutex.Unlock()

	result := make([]*Credential, 0)

	var e error = ErrNotFound

	for _, x := range s.items {
		if x.Equals(c) {
			e = nil
			continue
		} else {
			result = append(result, x)
		}
	}

	if e == nil {
		s.items = result
	}

	return e
}

// Add always...
func (store *SimpleStore) Add(item *Credential) error {

	if ex:= item.Validate() ; ex != nil {
		return ex
	}

	found := Exists(store, func(x *Credential) bool {
		return store.GetKey(x) == store.GetKey(item)
	})

	if !found {
		store.mutex.Lock()
		store.items = append(store.items, item)
		store.mutex.Unlock()
		return nil
	}

	return ErrAlreadyExists
}

// RemoveWhere should go , .. should be <tinystore>.RemoveWhere(store, filter ) error
func (s *SimpleStore) RemoveWhere(find Filter) error  {

	result := make([]*Credential, 0)
	var e error = ErrNotFound

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, x := range s.items[:] {
		if find(x) { // Skip
			e = nil
			continue
		}
		result = append(result, x)
	}
	if e ==nil {
		s.items = result
	}

	return e
}

// ForEach implements Store.ForEach
func (s *SimpleStore) ForEach(f Mutator) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var e error = nil
	for i, x := range s.items {
		r, e := f(x)
		if e == nil {
			s.items[i] = r
		}
	}
	return e
}

// ForEachWhere implements Store.ForEachWhere
func (s *SimpleStore) ForEachWhere(find Filter, transform Mutator) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var err error = ErrNotFound
	for i, x := range s.items {
		if find(x) {
			if r, e := transform(x) ;e == nil {
				s.items[i] = r
				return e
			}
		}
	}
	return err
}

// Specific t0 SimpleCredentialStore
func (store *SimpleStore) LoadJson(path string) error {

	bytes, e := ioutil.ReadFile(path)
	if e != nil {
		return e
	}
	e = json.Unmarshal(bytes, &store.items)
	if e != nil {
		return e
	}
	return nil
}

// GetKey implements Store.GetKey
func (store *SimpleStore) GetKey(item *Credential) interface{} {
	return item.Username
}

// Load implements Sore.Load, does not do error  checking
func (store *SimpleStore) Load(c ...*Credential) error {
	store.items = c
	// satisfy Interface
	return nil
}