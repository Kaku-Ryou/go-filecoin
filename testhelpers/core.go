package testhelpers

import (
	"context"
	"math/big"

	"gx/ipfs/QmR8BauakNcBa3RbE4nbQu76PDiJgoQgz8AJdhJuiU4TAw/go-cid"
	"gx/ipfs/QmRXf2uUSdGSunRJsM9wXSUNVwLUGCY3So5fAs7h2CBJVf/go-hamt-ipld"
	"gx/ipfs/QmS2aqUZLJp8kF1ihE5rvDGE5LvmKDPnx32w9Z1BW9xLV5/go-ipfs-blockstore"
	"gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"
	"gx/ipfs/Qmf4xQhNomPNhrtZc67qSnfJSjxjXs9LWvknJtSXwimPrM/go-datastore"

	"github.com/filecoin-project/go-filecoin/actor"
	"github.com/filecoin-project/go-filecoin/actor/builtin"
	"github.com/filecoin-project/go-filecoin/actor/builtin/account"
	"github.com/filecoin-project/go-filecoin/actor/builtin/miner"
	"github.com/filecoin-project/go-filecoin/address"
	"github.com/filecoin-project/go-filecoin/filnet"
	"github.com/filecoin-project/go-filecoin/state"
	"github.com/filecoin-project/go-filecoin/types"
	"github.com/filecoin-project/go-filecoin/vm"
	"github.com/stretchr/testify/require"
)

// RequireMakeStateTree takes a map of addresses to actors and stores them on
// the state tree, requiring that all its steps succeed.
func RequireMakeStateTree(require *require.Assertions, cst *hamt.CborIpldStore, acts map[address.Address]*actor.Actor) (cid.Cid, state.Tree) {
	ctx := context.Background()
	t := state.NewEmptyStateTreeWithActors(cst, builtin.Actors)

	for addr, act := range acts {
		err := t.SetActor(ctx, addr, act)
		require.NoError(err)
	}

	c, err := t.Flush(ctx)
	require.NoError(err)

	return c, t
}

// RequireNewEmptyActor creates a new empty actor with the given starting
// value and requires that its steps succeed.
func RequireNewEmptyActor(require *require.Assertions, value *types.AttoFIL) *actor.Actor {
	return &actor.Actor{Balance: value}
}

// RequireNewAccountActor creates a new account actor with the given starting
// value and requires that its steps succeed.
func RequireNewAccountActor(require *require.Assertions, value *types.AttoFIL) *actor.Actor {
	act, err := account.NewActor(value)
	require.NoError(err)
	return act
}

// RequireNewMinerActor creates a new miner actor with the given owner, pledge, and collateral,
// and requires that its steps succeed.
func RequireNewMinerActor(require *require.Assertions, vms vm.StorageMap, addr address.Address, owner address.Address, key []byte, pledge uint64, pid peer.ID, coll *types.AttoFIL) *actor.Actor {
	act := actor.NewActor(types.MinerActorCodeCid, types.NewZeroAttoFIL())
	storage := vms.NewStorage(addr, act)
	initializerData := miner.NewState(owner, key, big.NewInt(int64(pledge)), pid, coll)
	err := (&miner.Actor{}).InitializeState(storage, initializerData)
	require.NoError(err)
	require.NoError(storage.Flush())
	return act
}

// RequireNewFakeActor instantiates and returns a new fake actor and requires
// that its steps succeed.
func RequireNewFakeActor(require *require.Assertions, vms vm.StorageMap, addr address.Address, codeCid cid.Cid) *actor.Actor {
	return RequireNewFakeActorWithTokens(require, vms, addr, codeCid, types.NewAttoFILFromFIL(100))
}

// RequireNewFakeActorWithTokens instantiates and returns a new fake actor and requires
// that its steps succeed.
func RequireNewFakeActorWithTokens(require *require.Assertions, vms vm.StorageMap, addr address.Address, codeCid cid.Cid, amt *types.AttoFIL) *actor.Actor {
	act := actor.NewActor(codeCid, amt)
	store := vms.NewStorage(addr, act)
	err := (&actor.FakeActor{}).InitializeState(store, &actor.FakeActorStorage{})
	require.NoError(err)
	require.NoError(vms.Flush())
	return act
}

// RequireRandomPeerID returns a new libp2p peer ID or panics.
func RequireRandomPeerID() peer.ID {
	pid, err := filnet.RandPeerID()
	if err != nil {
		panic(err)
	}

	return pid
}

// VMStorage creates a new storage object backed by an in memory datastore
func VMStorage() vm.StorageMap {
	return vm.NewStorageMap(blockstore.NewBlockstore(datastore.NewMapDatastore()))
}

// MustSign signs a given address with the provided mocksigner or panics if it
// cannot.
func MustSign(s types.MockSigner, msgs ...*types.Message) []*types.SignedMessage {
	var smsgs []*types.SignedMessage
	for _, m := range msgs {
		gasLimit := types.NewGasUnits(999)
		sm, err := types.NewSignedMessage(*m, &s, types.NewGasPrice(0), gasLimit)
		if err != nil {
			panic(err)
		}
		smsgs = append(smsgs, sm)
	}
	return smsgs
}
