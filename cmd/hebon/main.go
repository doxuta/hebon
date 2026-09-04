// Command hebon spells kana in romaji under 令和7年内閣告示第4号, or rewrites
// older romaji into that spelling.
//
//	hebon [-style macron|circumflex|doubled] [-kunrei] [-cap] [-n] [text ...]
//
// With no arguments it converts standard input line by line.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/doxuta/hebon"
)

func main() {
	style := flag.String("style", "macron", "long vowels: macron (kāsan), circumflex (kâsan) or doubled (kaasan)")
	kunrei := flag.Bool("kunrei", false, "spell with the 1954 訓令式 table instead, for comparison")
	capitalize := flag.Bool("cap", false, "capitalise the first letter of every word")
	normalize := flag.Bool("n", false, "input is romaji: rewrite it to the notice's spelling")
	flag.Parse()

	opt := hebon.Options{Kunrei: *kunrei, Capitalize: *capitalize}
	switch *style {
	case "macron":
	case "circumflex":
		opt.Style = hebon.Circumflex
	case "doubled":
		opt.Style = hebon.Doubled
	default:
		fmt.Fprintf(os.Stderr, "hebon: unknown -style %q\n", *style)
		os.Exit(2)
	}
	convert := hebon.Romaji
	if *normalize {
		convert = hebon.Normalize
	}

	if flag.NArg() > 0 {
		fmt.Println(convert(strings.Join(flag.Args(), " "), opt))
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		fmt.Println(convert(sc.Text(), opt))
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "hebon:", err)
		os.Exit(1)
	}
}
