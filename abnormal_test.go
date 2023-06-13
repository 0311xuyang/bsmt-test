package main

import (
	"fmt"
	"testing"

	bsmt "github.com/bnb-chain/zkbnb-smt"
	"github.com/stretchr/testify/assert"
)

func TestNoCommit(t *testing.T) {
	for _, env := range prepareEnv() {
		smt, _ := initSMT(env)
		// set a text
		err := smt.Set(7, []byte("some data"))
		assert.NoError(t, err)
		// get will fail without commit
		_, err = smt.Get(7, nil)
		assert.Error(t, err)
	}
}

func TestVersionNotMatch(t *testing.T) {
	for _, env := range prepareEnv() {
		smt, _ := initSMT(env)

		ver0 := bsmt.Version(0)
		// set a text
		err := smt.Set(7, []byte("some data"))
		assert.NoError(t, err)
		fmt.Printf("latest version: %v recent version: %v\n", smt.LatestVersion(), smt.RecentVersion())

		// commit it, latest version incr 0 -> 1
		ver, err := smt.Commit(nil)
		fmt.Printf("latest version: %v recent version: %v\n", smt.LatestVersion(), smt.RecentVersion())
		assert.NoError(t, err)

		// trying get a higher version, will fail
		ver5 := bsmt.Version(5)
		_, err = smt.Get(7, &ver5)
		assert.EqualError(t, err, "the version is higher than the latest version")

		// set a text with a lower version, will fail
		err = smt.SetWithVersion(7, []byte("some other data"), ver0)
		assert.EqualError(t, err, "the version is lower than the latest version")

		// set a text with a specific higher version, latest version incr 1 -> 5
		err = smt.SetWithVersion(7, []byte("some other data"), ver5)
		assert.NoError(t, err)

		// when set with a specific version, commit should be consostent
		ver, err = smt.CommitWithNewVersion(&ver, &ver5)
		fmt.Printf("latest version: %v recent version: %v\n", smt.LatestVersion(), smt.RecentVersion())

		// the new text is set
		res, err := smt.Get(7, nil)
		assert.Equal(t, string(res), "some other data")

		// try to set with a higher version 10
		ver10 := bsmt.Version(10)
		err = smt.SetWithVersion(7, []byte("some latest data"), ver10)
		assert.NoError(t, err)

		// if use normal commit, latest version incr 5 -> 6, won't contain the version 10
		ver, err = smt.Commit(nil)
		fmt.Printf("latest version: %v recent version: %v\n", smt.LatestVersion(), smt.RecentVersion())

		// the text still version 5
		res, err = smt.Get(7, nil)
		assert.Equal(t, string(res), "some other data")
	}
}

func TestProofNotMatch(t *testing.T) {
	for _, env := range prepareEnv() {
		smt, _ := initSMT(env)

		/*
			1
			|
			|-- 2
			|   |
			|   |-- 4
			|   |-- 5
			|
			|-- 3
			    |
			    |-- 6
			    |-- 7
		*/

		err := smt.Set(1, []byte("data 1"))
		assert.NoError(t, err)
		err = smt.Set(2, []byte("data 2"))
		assert.NoError(t, err)
		err = smt.Set(3, []byte("data 3"))
		assert.NoError(t, err)
		err = smt.Set(4, []byte("data same"))
		assert.NoError(t, err)
		err = smt.Set(5, []byte("data same"))
		assert.NoError(t, err)
		err = smt.Set(6, []byte("data 6"))
		assert.NoError(t, err)
		err = smt.Set(7, []byte("data 7"))
		assert.NoError(t, err)

		_, err = smt.Commit(nil)
		assert.NoError(t, err)

		pf5, err := smt.GetProof(5)
		assert.NoError(t, err)

		// proof is correct
		ok := smt.VerifyProof(5, pf5)
		assert.True(t, ok)

		// proof won't work with other node
		ok = smt.VerifyProof(6, pf5)
		assert.False(t, ok)

		// proof will work if the routes are same
		ok = smt.VerifyProof(4, pf5)
		assert.True(t, ok)

		// proof won't work if changed
		pf5[1][2] = 0
		ok = smt.VerifyProof(5, pf5)
		assert.False(t, ok)
	}
}
