package main

import (
	"fmt"

	"github.com/tidwall/btree"
)

// https://github.com/tidwall/btree
// https://github.com/tidwall/btree#example-2
type Item struct {
	Key, Val string
}

// byKeys is a comparison function that compares item keys and returns true
// when a is less than b.
func byKeys(a, b Item) bool {
	return a.Key < b.Key
}

// byVals is a comparison function that compares item values and returns true
// when a is less than b.
func byVals(a, b Item) bool {
	if a.Val < b.Val {
		return true
	}
	if a.Val > b.Val {
		return false
	}
	// Both vals are equal so we should fall though
	// and let the key comparison take over.
	return byKeys(a, b)
}

func main() {
	// Create a tree for keys and a tree for values.
	// The "keys" tree will be sorted on the Keys field.
	// The "values" tree will be sorted on the Values field.
	keys := btree.NewBTreeG[Item](byKeys)
	vals := btree.NewBTreeG[Item](byVals)

	// Create some items.
	users := []Item{
		Item{Key: "user:1", Val: "Jane"},
		Item{Key: "user:2", Val: "Andy"},
		Item{Key: "user:3", Val: "Steve"},
		Item{Key: "user:4", Val: "Andrea"},
		Item{Key: "user:5", Val: "Janet"},
		Item{Key: "user:6", Val: "Andy"},
	}

	// Insert each user into both trees
	for _, user := range users {
		keys.Set(user)
		vals.Set(user)
	}

	// Iterate over each user in the key tree
	keys.Scan(func(item Item) bool {
		fmt.Printf("%s %s\n", item.Key, item.Val)
		return true
	})
	fmt.Printf("\n")

	// Iterate over each user in the val tree
	vals.Scan(func(item Item) bool {
		fmt.Printf("%s %s\n", item.Key, item.Val)
		return true
	})

	// Output:
	// user:1 Jane
	// user:2 Andy
	// user:3 Steve
	// user:4 Andrea
	// user:5 Janet
	// user:6 Andy
	//
	// user:4 Andrea
	// user:2 Andy
	// user:6 Andy
	// user:1 Jane
	// user:5 Janet
	// user:3 Steve
}
