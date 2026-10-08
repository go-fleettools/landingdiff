package main

import (
	"path/filepath"
	"testing"
)

func TestSpreadNamesTheOneTheOthersDisagreeWith(t *testing.T) {
	for _, c := range []struct {
		name     string
		v        []int
		lo, hi   int
		odd      int
		whyNoOdd string
	}{
		{name: "one exception to a consensus", v: []int{3, 3, 3, 4}, lo: 3, hi: 4, odd: 3},
		{name: "the exception first", v: []int{4, 3, 3, 3}, lo: 3, hi: 4, odd: 0},
		{name: "everybody agrees", v: []int{2, 2, 2}, lo: 2, hi: 2, odd: -1,
			whyNoOdd: "a column everybody agrees on has no exception in it"},

		// ⛔ The case the rule was briefly got wrong for. With every value
		// different there is no consensus to be outside of, and a rule that
		// names the furthest from the rest names the smallest and the largest
		// by arithmetic — every column grows an outlier and the report stops
		// meaning anything.
		{name: "all different", v: []int{14, 10, 19, 23}, lo: 10, hi: 23, odd: -1,
			whyNoOdd: "four different numbers are a range, not a rule with an exception"},

		// Two against two is a disagreement, not an exception.
		{name: "two and two", v: []int{1, 1, 2, 2}, lo: 1, hi: 2, odd: -1,
			whyNoOdd: "two pages each is not one page being alone"},
	} {
		t.Run(c.name, func(t *testing.T) {
			lo, hi, odd := spread(c.v)
			if lo != c.lo || hi != c.hi {
				t.Errorf("%v spans %d..%d, want %d..%d", c.v, lo, hi, c.lo, c.hi)
			}
			if odd != c.odd {
				t.Errorf("%v named %d as the odd one, want %d — %s", c.v, odd, c.odd, c.whyNoOdd)
			}
		})
	}
}

// A page with one of everything the measures look for.
const page = `<div class="nav"><nav><a href="#a">A</a><a href="#b">B</a></nav></div>
<section class="hero"><div class="wrap">
  <h1 class="head">Two questions,<br><span class="accent">and they are not the same.</span></h1>
  <div class="cta"><a class="btn btn-primary" href="#">One</a><a class="btn btn-ghost" href="#">Two</a></div>
  <div class="trust"><span>a</span><span>b</span><span>c</span></div>
</div></section>
<section class="strip"><div class="mods"><a>x</a></div></section>`

func TestEachMeasureCountsWhatItSaysItDoes(t *testing.T) {
	want := map[string]int{
		"headline line 1, characters": 14,
		"headline, total characters":  40,
		"items in the trust band":     3,
		"buttons in the hero":         2,
		"links in the navigation":     2,
		"blocks in the strip":         1,
		"sections":                    2,
	}
	for _, m := range measures {
		w, ok := want[m.name]
		if !ok {
			t.Fatalf("there is a measure called %q that this test says nothing about", m.name)
		}
		if got := m.of(page); got != w {
			t.Errorf("%s: %d, want %d", m.name, got, w)
		}
	}
	if len(want) != len(measures) {
		t.Errorf("%d measures and %d expectations", len(measures), len(want))
	}
}

func TestAPageMissingThePartAMeasureLooksForSaysSo(t *testing.T) {
	// ⛔ -1, not 0. A page with no trust band at all and a page with an empty
	// one are different things, and reporting both as zero would put them in
	// the same column and call one of them an exception.
	for _, m := range measures {
		got := m.of("<html><body>nothing here</body></html>")
		if got > 0 {
			t.Errorf("%s found %d of them in a page that has none", m.name, got)
		}
	}
}

func TestOrgOfIsTheDirectoryTwoAbove(t *testing.T) {
	// The path shape the fleet keeps its landings in, built with the
	// separator this machine uses so the three CI lanes measure the same
	// thing.
	at := filepath.Join("whatever", "landing-go-authn", "layouts", "index.html")
	if got := orgOf(at); got != "go-authn" {
		t.Errorf("orgOf(%q) = %q", at, got)
	}
	bare := filepath.Join("go-pkgx.github.io", "layouts", "index.html")
	if got := orgOf(bare); got != "go-pkgx.github.io" {
		t.Errorf("orgOf(%q) = %q, and a directory with no landing- prefix keeps its name", bare, got)
	}
}
