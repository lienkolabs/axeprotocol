package state

import (
	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/protocol/chain"
)

type Mutations struct {
	GrantPower  map[crypto.Hash]struct{}
	RevokePower map[crypto.Hash]struct{}
	NewMembers  map[crypto.Hash]struct{}
	NewCaption  map[crypto.Hash]struct{}
}

func (m *Mutations) HasGrantPower(hash crypto.Hash) bool {
	_, ok := m.GrantPower[hash]
	return ok
}

func (m *Mutations) HasRevokePower(hash crypto.Hash) bool {
	_, ok := m.RevokePower[hash]
	return ok
}

func (m *Mutations) HasMember(hash crypto.Hash) bool {
	_, ok := m.NewMembers[hash]
	return ok
}

func (m *Mutations) HasCaption(hash crypto.Hash) bool {
	_, ok := m.NewCaption[hash]
	return ok
}

func (m *Mutations) Append(array []chain.Mutations) chain.Mutations {
	grouped := &Mutations{
		GrantPower:  make(map[crypto.Hash]struct{}),
		RevokePower: make(map[crypto.Hash]struct{}),
		NewMembers:  make(map[crypto.Hash]struct{}),
		NewCaption:  make(map[crypto.Hash]struct{}),
	}
	all := []*Mutations{m}
	if len(array) > 0 {
		for _, a := range array {
			if mutation, ok := a.(*Mutations); ok {
				all = append(all, mutation)
			}
		}
	}
	for _, mutations := range all {
		for hash := range mutations.GrantPower {
			grouped.GrantPower[hash] = struct{}{}
		}
		for hash := range mutations.RevokePower {
			grouped.RevokePower[hash] = struct{}{}
			delete(grouped.GrantPower, hash)
		}

		for hash := range mutations.NewMembers {
			grouped.NewMembers[hash] = struct{}{}
		}

		for hash := range mutations.NewCaption {
			grouped.NewCaption[hash] = struct{}{}
		}
	}
	return grouped
}
