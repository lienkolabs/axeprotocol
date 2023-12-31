package state

import (
	"github.com/lienkolabs/axeprotocol/social/actions"
	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/protocol/chain"
)

type MutatingState struct {
	Epoch     uint64
	state     *State
	mutations *Mutations
}

func (m *MutatingState) Mutations() chain.Mutations {
	return m.mutations
}

func (s *MutatingState) SetNewGrantPower(token, attorney crypto.Token) bool {
	join := append(token[:], attorney[:]...)
	hash := crypto.Hasher(join)
	s.mutations.GrantPower[hash] = struct{}{}
	return true
}

func (s *MutatingState) SetNewRevokePower(token, attorney crypto.Token) bool {
	join := append(token[:], attorney[:]...)
	hash := crypto.Hasher(join)
	s.mutations.GrantPower[hash] = struct{}{}
	return true
}

func (s *MutatingState) SetNewMember(token crypto.Token, handle string) bool {
	if (!s.HasHandle(handle)) && (!s.state.HasMember(token)) {
		captionHash := crypto.Hasher([]byte(handle))
		tokenHash := crypto.HashToken(token)
		s.mutations.NewMembers[tokenHash] = struct{}{}
		s.mutations.NewCaption[captionHash] = struct{}{}
		return true
	}
	return false
}

func (s *MutatingState) PowerOfAttorney(token, attorney crypto.Token) bool {
	if token.Equal(attorney) {
		return true
	}
	join := append(token[:], attorney[:]...)
	hash := crypto.Hasher(join)
	_, ok := s.mutations.GrantPower[hash]
	return ok || s.state.Attorneys.ExistsHash(hash)
}

func (s *MutatingState) HasMember(token crypto.Token) bool {
	hash := crypto.HashToken(token)
	_, ok := s.mutations.NewMembers[hash]
	return ok || s.state.Members.ExistsHash(hash)
}

func (s *MutatingState) HasHandle(handle string) bool {
	hash := crypto.Hasher([]byte(handle))
	_, ok := s.mutations.NewCaption[hash]
	return ok || s.state.Captions.ExistsHash(hash)
}

func (v *MutatingState) Validate(data []byte) bool {
	kind := actions.Kind(data)
	if kind == actions.Invalid {
		return false
	}
	switch kind {
	case actions.JoinNetworkType:
		join := actions.ParseJoinNetwork(data)
		return v.SetNewMember(join.Author, join.Handle)
	case actions.UpdateInfoType:
		update := actions.ParseUpdateInfo(data)
		if !v.PowerOfAttorney(update.Author, update.Signer) {
			return false
		}
		return v.HasMember(update.Author)
	case actions.GrantPowerOfAttorneyType:
		grant := actions.ParseGrantPowerOfAttorney(data)
		if !v.HasMember(grant.Author) {
			return false
		}
		return v.SetNewGrantPower(grant.Author, grant.Attorney)
	case actions.RevokePowerOfAttorneyType:
		revoke := actions.ParseRevokePowerOfAttorney(data)
		if !v.HasMember(revoke.Author) {
			return false
		}
		return v.SetNewRevokePower(revoke.Author, revoke.Attorney)
	case actions.VoidType:
		void := actions.ParseVoid(data)
		if !v.HasMember(void.Author) {
			return false
		}
		return v.PowerOfAttorney(void.Author, void.Signer)
	}
	return false
}
