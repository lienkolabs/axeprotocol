package state

import (
	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/protocol/chain"
)

type State struct {
	Epoch     uint64
	Members   *hashVault
	Captions  *hashVault
	Attorneys *hashVault
}

func NewGenesisState(dataPath string) *State {
	state := State{
		Epoch:     0,
		Members:   NewHashVault("members", 0, 8, dataPath),
		Captions:  NewHashVault("captions", 0, 8, dataPath),
		Attorneys: NewHashVault("poa", 0, 8, dataPath),
	}
	return &state
}

func (s *State) Validator(v chain.Mutations, epoch uint64) chain.MutatingState {
	m, ok := v.(*Mutations)
	if !ok {
		return nil
	}

	return &MutatingState{
		Epoch:     epoch,
		state:     s,
		mutations: m,
	}
}

func (s *State) Incorporate(v chain.MutatingState, token crypto.Token) {
	m, ok := v.(*MutatingState)
	if !ok {
		return
	}
	for hash := range m.mutations.GrantPower {
		s.Attorneys.InsertHash(hash)
	}
	for hash := range m.mutations.RevokePower {
		s.Attorneys.RemoveHash(hash)
	}
	for hash := range m.mutations.NewMembers {
		s.Members.InsertHash(hash)
	}
	for hash := range m.mutations.NewCaption {
		s.Captions.ExistsHash(hash)
	}
}

func (s *State) NewMutations() chain.Mutations {
	return &Mutations{
		GrantPower:  make(map[crypto.Hash]struct{}),
		RevokePower: make(map[crypto.Hash]struct{}),
		NewMembers:  make(map[crypto.Hash]struct{}),
		NewCaption:  make(map[crypto.Hash]struct{}),
	}
}

func (s *State) PowerOfAttorney(token, attorney crypto.Token) bool {
	if token.Equal(attorney) {
		return true
	}
	join := append(token[:], attorney[:]...)
	hash := crypto.Hasher(join)
	return s.Attorneys.ExistsHash(hash)
}

func (s *State) HasMember(token crypto.Token) bool {
	hash := crypto.HashToken(token)
	return s.Members.ExistsHash(hash)
}

func (s *State) HasHandle(handle string) bool {
	hash := crypto.Hasher([]byte(handle))
	return s.Captions.ExistsHash(hash)
}

func (s *State) Shutdown() {
	s.Members.Close()
	s.Attorneys.Close()
	s.Captions.Close()
}
