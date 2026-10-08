// landingdiff says which of a set of landing pages is the odd one out.
//
// The fleet's landings share one styles.html, so "aligned" cannot be checked
// by diffing the design system — it is already identical. What drifts is what
// each page PUTS IN it: a headline too long for the column it is set in, a
// band with one more item than fits on a line, a navigation with twice as many
// links as its neighbours.
//
// ⛔ It reports the SPREAD and names the outlier rather than asserting a rule.
// The first cut of this checked "three cards to a grid" because two pages had
// three — and the fleet turned out to use four and nine as happily, so the
// rule was invented rather than observed. A column that disagrees everywhere
// is a column with no house style, and saying so is the useful answer.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// A measure is one countable property of a page.
type measure struct {
	name string
	of   func(string) int
}

var measures = []measure{
	{"headline line 1, characters", func(s string) int {
		m := regexp.MustCompile(`<h1 class="head">([^<]*)`).FindStringSubmatch(s)
		if m == nil {
			return -1
		}
		return len([]rune(m[1]))
	}},
	{"headline, total characters", func(s string) int {
		m := regexp.MustCompile(`(?s)<h1 class="head">(.*?)</h1>`).FindStringSubmatch(s)
		if m == nil {
			return -1
		}
		return len([]rune(strings.TrimSpace(tags.ReplaceAllString(m[1], ""))))
	}},
	{"items in the trust band", func(s string) int {
		return countIn(s, `<div class="trust">`, `</div>`, "<span>")
	}},
	{"buttons in the hero", func(s string) int {
		return countIn(s, `<div class="cta">`, `</div>`, `class="btn`)
	}},
	{"links in the navigation", func(s string) int {
		return countIn(s, "<nav", "</nav>", "<a ")
	}},
	{"blocks in the strip", func(s string) int { return strings.Count(s, `class="mods"`) }},
	{"sections", func(s string) int { return strings.Count(s, "<section") }},
}

var tags = regexp.MustCompile(`<[^>]*>`)

// countIn counts needles between the first open and the next close after it.
func countIn(s, open, close, needle string) int {
	i := strings.Index(s, open)
	if i < 0 {
		return -1
	}
	rest := s[i+len(open):]
	if j := strings.Index(rest, close); j >= 0 {
		rest = rest[:j]
	}
	return strings.Count(rest, needle)
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "landingdiff: give two or more layouts/index.html to compare")
		os.Exit(2)
	}

	names := make([]string, 0, flag.NArg())
	pages := make([]string, 0, flag.NArg())
	for _, at := range flag.Args() {
		b, err := os.ReadFile(at)
		if err != nil {
			fmt.Fprintln(os.Stderr, "landingdiff:", err)
			os.Exit(1)
		}
		names = append(names, orgOf(at))
		pages = append(pages, string(b))
	}

	w := 0
	for _, n := range names {
		if len(n) > w {
			w = len(n)
		}
	}
	fmt.Printf("%-30s", "")
	for _, n := range names {
		fmt.Printf(" %*s", w, n)
	}
	fmt.Println("   spread")

	odd := map[string]int{}
	for _, m := range measures {
		got := make([]int, len(pages))
		for i, p := range pages {
			got[i] = m.of(p)
		}
		fmt.Printf("%-30s", m.name)
		for _, v := range got {
			fmt.Printf(" %*d", w, v)
		}
		lo, hi, who := spread(got)
		if hi-lo > 0 {
			fmt.Printf("   %d..%d", lo, hi)
			if who >= 0 {
				fmt.Printf("  <- %s", names[who])
				odd[names[who]]++
			}
		}
		fmt.Println()
	}

	if len(odd) == 0 {
		fmt.Println("\nNothing is an outlier: every measure agrees across the set.")
		return
	}
	type pair struct {
		name string
		n    int
	}
	var ps []pair
	for k, v := range odd {
		ps = append(ps, pair{k, v})
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].n > ps[j].n })
	fmt.Println()
	for _, p := range ps {
		fmt.Printf("%s is alone on %d measure(s)\n", p.name, p.n)
	}
}

// spread is the range, and which page sits OUTSIDE the range of all the
// others.
//
// A page is named only when the others AGREE and it does not: one value held
// by everybody else, one page holding something different. That is a house
// style with an exception in it, which is the only thing a count can honestly
// report.
//
// ⛔ It deliberately does NOT flag "furthest from the rest". The rule was
// briefly leave-one-out — outside the range of the others — and with four
// pages whose values are all different that names the smallest and the largest
// every time, by arithmetic rather than by observation: every column grew an
// outlier and the output stopped meaning anything.
//
// ⛔ And it is blind on purpose to the defect that started this. A headline of
// 23 characters against 10, 14 and 19 is not a categorical exception, it is a
// line too long for the column it is set in — a fact about RENDERING that no
// count of anything can see. That one was found by screenshotting the page and
// looking at it, and there is no honest way to automate it here.
func spread(v []int) (lo, hi, odd int) {
	lo, hi = bounds(v)
	count := map[int]int{}
	for _, x := range v {
		count[x]++
	}
	if len(count) != 2 {
		return lo, hi, -1
	}
	for i, x := range v {
		if count[x] == 1 {
			return lo, hi, i
		}
	}
	return lo, hi, -1
}

func bounds(v []int) (lo, hi int) {
	lo, hi = v[0], v[0]
	for _, x := range v {
		if x < lo {
			lo = x
		}
		if x > hi {
			hi = x
		}
	}
	return lo, hi
}

// orgOf is the organisation a path belongs to: .../landing-go-authn/layouts/index.html
func orgOf(at string) string {
	d := filepath.Dir(filepath.Dir(at))
	return strings.TrimPrefix(filepath.Base(d), "landing-")
}
