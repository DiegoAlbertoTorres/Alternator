package alternator

import (
	"container/list"
	"fmt"
)

// Members stores the members of a node
type Members struct {
	List *list.List
	// Slice []*Peer
	Map map[Key]*list.Element
}

// Init initializes a members struct
func (members *Members) Init() {
	members.List = list.New()
	members.Map = make(map[Key]*list.Element)
	// altNode.Members.Add(&selfExt)
}

// GetPeer gets a peer from an element of Members.List
func GetPeer(e *list.Element) *Peer {
	if e != nil {
		return e.Value.(*Peer)
	}
	return nil
}

// Inserts a member
// func (members *Members) sliceInsert(new *Peer, i int) {
// 	// Grow the slice by one element.
// 	members.Slice = append(members.Slice, &Peer{})
// 	copy(members.Slice[i+1:], members.Slice[i:])
// 	members.Slice[i] = new
// }

// FindSuccessor finds the successor of a key in the ring
func (members *Members) FindSuccessor(key Key) (*list.Element, error) {
	// Find ID of successor
	for e := members.List.Front(); e != nil; e = e.Next() {
		if Compare(GetPeer(e).ID, key) > 0 {
			return e, nil
		}
	}
	// Nothing bigger in ring, successor is first node
	return (members.List.Front()), nil
}

// Insert adds a node to the members if not already there, returns true if added
func (members *Members) Insert(peer *Peer) *list.Element {
	i := 0
	// Iterate through members
	for e := members.List.Front(); e != nil; e = e.Next() {
		if e == nil {
			break
		}
		current := GetPeer(e)
		var prevID Key
		if prev := e.Prev(); prev != nil {
			prevID = GetPeer(prev).ID
		} else {
			prevID = MinKey
		}
		// Already in Members
		if Compare(peer.ID, current.ID) == 0 {
			return nil
		} else if (Compare(peer.ID, prevID) > 0) && (Compare(peer.ID, current.ID) < 0) {
			// Correct spot to add, add to list, slice and map
			// members.sliceInsert(peer, i)
			e := members.List.InsertBefore(peer, e)
			members.Map[peer.ID] = e
			return e
		}
		i++
	}
	// Append at end if not in members and not added by loop
	// members.sliceInsert(peer, 0)
	e := members.List.PushBack(peer)
	members.Map[peer.ID] = e
	return e
}

// Remove deletes a node from members
func (members *Members) Remove(del *Peer) {
	for e := members.List.Front(); e != nil; e = e.Next() {
		if GetPeer(e).ID == del.ID {
			members.List.Remove(e)
			delete(members.Map, del.ID)
		}
	}
}

func (members Members) String() (str string) {
	i := 0
	for e := members.List.Front(); e != nil; e = e.Next() {
		str += fmt.Sprintf("member %d: %s\n", i, GetPeer(e).String())
		i++
	}
	return
}

// Compares member table to neighbors, adds new
// func (altNode *Alternator) updateMembers() {
// 	var successorMembers []*Peer
// 	// Get successors members
// 	if altNode.Successor != nil {
// 		MakeRemoteCall(altNode.Successor, "GetMembers", struct{}{}, &successorMembers)
// 	}
//
// 	for _, member := range successorMembers {
// 		altNode.Members.AddIfMissing(member)
// 		// elem := altNode.Members.AddIfMissing(member)
// 		// if elem != nil {
// 		// 	// Send this node the keys that belong to it
// 		// 	altNode.expelForeignKeys(elem)
// 		// }
// 	}
// }

// GetRandomMember returns a random member from the ring
func (members Members) GetRandomMember() *Peer {
	// i := 0
	// random := rand.Intn(len(members.Map))
	// for _, member := range members.Map {
	// 	if i == random {
	// 		return GetPeer(member)
	// 	}
	// 	i++
	// }
	for _, member := range members.Map {
		return GetPeer(member)
	}
	return nil
}