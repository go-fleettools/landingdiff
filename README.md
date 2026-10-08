# landingdiff

Which of a set of landing pages is **the odd one out**.

```
landingdiff ../landing-go-*/layouts/index.html
```

```
                               go-authn go-fileshare  go-pkgx go-pdfkit   spread
headline line 1, characters          14           10       19        23   10..23
items in the trust band               3            3        3         4   3..4  <- go-pdfkit
buttons in the hero                   2            2        2         2
links in the navigation               3            3        6         4   3..6
blocks in the strip                   1            1        1         2   1..2  <- go-pdfkit

go-pdfkit is alone on 2 measure(s)
```

## Why counting, and not diffing

The fleet's landings share one `styles.html` — three of the four above are byte
for byte identical — so "aligned" cannot be checked by diffing the design
system. What drifts is **what each page puts in it**: a band with one more item
than fits on a line, a navigation with twice as many links as its neighbours.

## What it will not tell you

⛔ **A page is named only when the others agree and it does not.** One value
held by everybody else, one page holding something different: that is a house
style with an exception in it, which is the only thing a count can honestly
report.

It deliberately does *not* flag "furthest from the rest". That rule looks
better and is worthless: with four pages whose values are all different it
names the smallest and the largest **every time**, by arithmetic rather than by
observation — every column grows an outlier and the output stops meaning
anything.

⛔ **And it is blind to layout.** A headline of 23 characters against 10, 14 and
19 is not a categorical exception; it is a line too long for the column it is
set in, which is a fact about *rendering*. That one is found by screenshotting
the pages and looking at them:

```sh
for o in go-authn go-fileshare go-pkgx go-pdfkit; do
  chrome --headless --window-size=1280,1600 \
    --screenshot="$o.png" "https://$o.github.io/"
done
```

Use both. The count finds the exceptions, the picture finds the overflows, and
neither finds the other's.

## Before you "fix" an outlier

Read the comment on the line first. On the run above, go-pdfkit's second strip
block was flagged — and it was **asked for**, derived from `[params.arches]`
rather than typed, with a comment above it saying so. An outlier that the
history explains is not drift.

## Adding a measure

One entry in `measures`: a name, and a function from the page's HTML to a
count. Return **-1**, not 0, when the part is not there at all — a page with no
trust band and a page with an empty one are different things, and reporting
both as zero puts them in the same column and calls one of them an exception.

## Licence

BSD-3-Clause.
