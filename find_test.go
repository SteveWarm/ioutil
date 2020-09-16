package ioutil

import "testing"

func TestFindWithOutMatch(t *testing.T) {
	frr, drr, err := Find(Config{
		Dir:        "/tmp",
		AppendFile: true,
		AppendDir:  true,
		MinDepth:   2,
		MaxDepth:   3})

	if err != nil {
		t.Fatal(err)
	}

	t.Log(frr)
	t.Log(drr)
}

func TestFindWithRegexMatch(t *testing.T) {
	frr, drr, err := Find(Config{
		Dir:        "/tmp",
		RegexMatch: []string{"[0-9]+$", "\\.conf$"},
		AppendFile: true,
		AppendDir:  true,
		MinDepth:   2,
		MaxDepth:   5})

	if err != nil {
		t.Fatal(err)
	}

	t.Log(frr)
	t.Log(drr)
}
