package trie_test

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/31z4/harvest2/internal/trie"
)

func TestTrie(t *testing.T) {
	type prefixData struct {
		Prefix string
		Count  uint
	}

	tests := []struct {
		insertValues     []string
		expectedPrefixes []prefixData
	}{
		{
			[]string{},
			nil,
		},
		{
			[]string{"test"},
			[]prefixData{
				{"test", 1},
			},
		},
		{
			[]string{"test", "slow"},
			[]prefixData{
				{"test", 1},
				{"slow", 1},
			},
		},
		{
			[]string{"test", "slow", "water", "slower"},
			[]prefixData{
				{"test", 1},
				{"slow", 2},
				{"slower", 1},
				{"water", 1},
			},
		},
		{
			[]string{"tester", "test"},
			[]prefixData{
				{"test", 2},
				{"tester", 1},
			},
		},
		{
			[]string{"test", "team"},
			[]prefixData{
				{"te", 2},
				{"test", 1},
				{"team", 1},
			},
		},
		{
			[]string{"test", "toaster", "toasting", "slow", "slowly"},
			[]prefixData{
				{"t", 3},
				{"test", 1},
				{"toast", 2},
				{"toaster", 1},
				{"toasting", 1},
				{"slow", 2},
				{"slowly", 1},
			},
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tree := trie.New()
			for _, v := range tc.insertValues {
				tree.Insert(v)
			}

			var walkedPrefixes []prefixData
			tree.Walk(func(prefix string, count uint) {
				walkedPrefixes = append(walkedPrefixes, prefixData{prefix, count})
			})

			if !cmp.Equal(walkedPrefixes, tc.expectedPrefixes) {
				t.Errorf("unexpected walkedPrefixes:\n%s", cmp.Diff(walkedPrefixes, tc.expectedPrefixes))
			}
		})
	}
}
