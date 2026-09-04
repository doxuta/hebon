// Package hebon spells Japanese kana in Latin letters the way
// 令和7年内閣告示第4号「ローマ字のつづり方」(22 December 2025) says to, and
// rewrites romaji written under the 1954 notice (訓令式) into that spelling.
//
// The notice consists of a 本表 (the kana table), nine 添え書き (rules for
// ん, っ, long vowels, the apostrophe, hyphens, capitals and particles) and a
// 対照表 that lines the new spellings up against both 1954 tables. All three
// are reproduced here as data; the engine is the 添え書き.
package hebon

import (
	"strings"
	"unicode"
)

// Style is how a long vowel is written (添え書き 3).
type Style int

const (
	// Macron writes kāsan — 添え書き 3(1), the notice's first choice.
	Macron Style = iota
	// Circumflex writes kâsan — 3(1) allows "^" where a macron is unavailable.
	Circumflex
	// Doubled writes kaasan — 3(2), letters follow the kana as in 現代仮名遣い.
	Doubled
)

// Options control Romaji and Normalize.
type Options struct {
	Style Style
	// Kunrei spells with the 1954 第1表 (訓令式: si, ti, tu, hu, zi, sya…)
	// instead of the 本表. Useful for showing what changed.
	Kunrei bool
	// Capitalize upper-cases the first letter of every whitespace-separated
	// word (添え書き 6 and 7: proper nouns and sentence starts).
	Capitalize bool
}

// honpyo is the 本表 of the notice, keyed by hiragana. Katakana input is
// folded to hiragana before lookup. ぢ/づ/ヲ follow the ※ note: no
// distinction from じ/ず/オ.
var honpyo = map[string]string{
	"あ": "a", "い": "i", "う": "u", "え": "e", "お": "o",
	"か": "ka", "き": "ki", "く": "ku", "け": "ke", "こ": "ko", "きゃ": "kya", "きゅ": "kyu", "きょ": "kyo",
	"さ": "sa", "し": "shi", "す": "su", "せ": "se", "そ": "so", "しゃ": "sha", "しゅ": "shu", "しょ": "sho",
	"た": "ta", "ち": "chi", "つ": "tsu", "て": "te", "と": "to", "ちゃ": "cha", "ちゅ": "chu", "ちょ": "cho",
	"な": "na", "に": "ni", "ぬ": "nu", "ね": "ne", "の": "no", "にゃ": "nya", "にゅ": "nyu", "にょ": "nyo",
	"は": "ha", "ひ": "hi", "ふ": "fu", "へ": "he", "ほ": "ho", "ひゃ": "hya", "ひゅ": "hyu", "ひょ": "hyo",
	"ま": "ma", "み": "mi", "む": "mu", "め": "me", "も": "mo", "みゃ": "mya", "みゅ": "myu", "みょ": "myo",
	"や": "ya", "ゆ": "yu", "よ": "yo",
	"ら": "ra", "り": "ri", "る": "ru", "れ": "re", "ろ": "ro", "りゃ": "rya", "りゅ": "ryu", "りょ": "ryo",
	"わ": "wa", "を": "o",
	"が": "ga", "ぎ": "gi", "ぐ": "gu", "げ": "ge", "ご": "go", "ぎゃ": "gya", "ぎゅ": "gyu", "ぎょ": "gyo",
	"ざ": "za", "じ": "ji", "ず": "zu", "ぜ": "ze", "ぞ": "zo", "じゃ": "ja", "じゅ": "ju", "じょ": "jo",
	"だ": "da", "ぢ": "ji", "づ": "zu", "で": "de", "ど": "do", "ぢゃ": "ja", "ぢゅ": "ju", "ぢょ": "jo",
	"ば": "ba", "び": "bi", "ぶ": "bu", "べ": "be", "ぼ": "bo", "びゃ": "bya", "びゅ": "byu", "びょ": "byo",
	"ぱ": "pa", "ぴ": "pi", "ぷ": "pu", "ぺ": "pe", "ぽ": "po", "ぴゃ": "pya", "ぴゅ": "pyu", "ぴょ": "pyo",
}

// kunrei is the 1954 第1表 wherever it differs from the 本表 — the middle
// column of the notice's 対照表.
var kunrei = map[string]string{
	"し": "si", "ち": "ti", "つ": "tu", "ふ": "hu", "じ": "zi", "ぢ": "zi",
	"しゃ": "sya", "しゅ": "syu", "しょ": "syo", "ちゃ": "tya", "ちゅ": "tyu", "ちょ": "tyo",
	"じゃ": "zya", "じゅ": "zyu", "じょ": "zyo", "ぢゃ": "zya", "ぢゅ": "zyu", "ぢょ": "zyo",
}

// gairaigo covers the kana that 外来語の表記 (平成3年内閣告示第2号) adds
// for loanwords: the 第1表 additions and the whole 第2表. The romanization
// notice excludes these sounds (前書き 4), so the spellings below are the
// usual Hepburn ones by analogy with the 本表, not text of the notice.
var gairaigo = map[string]string{
	// 第1表 additions
	"しぇ": "she", "ちぇ": "che", "つぁ": "tsa", "つぇ": "tse", "つぉ": "tso",
	"てぃ": "ti", "でぃ": "di", "ふぁ": "fa", "ふぃ": "fi", "ふぇ": "fe", "ふぉ": "fo",
	"じぇ": "je", "でゅ": "dyu",
	// 第2表
	"いぇ": "ye", "うぃ": "wi", "うぇ": "we", "うぉ": "wo",
	"くぁ": "kwa", "くぃ": "kwi", "くぇ": "kwe", "くぉ": "kwo", "ぐぁ": "gwa",
	"つぃ": "tsi", "とぅ": "tu", "どぅ": "du",
	"ゔぁ": "va", "ゔぃ": "vi", "ゔ": "vu", "ゔぇ": "ve", "ゔぉ": "vo",
	"てゅ": "tyu", "ふゅ": "fyu", "ゔゅ": "vyu",
}

// extra: small kana standing alone and the two historical kana. Not in the
// notice; here so that real text does not fall through unconverted.
var extra = map[string]string{
	"ぁ": "a", "ぃ": "i", "ぅ": "u", "ぇ": "e", "ぉ": "o",
	"ゃ": "ya", "ゅ": "yu", "ょ": "yo", "ゎ": "wa",
	"ゐ": "i", "ゑ": "e",
}

func lookup(k string, kunreiOn bool) (string, bool) {
	if kunreiOn {
		if v, ok := kunrei[k]; ok {
			return v, true
		}
	}
	if v, ok := honpyo[k]; ok {
		return v, true
	}
	if v, ok := gairaigo[k]; ok {
		return v, true
	}
	v, ok := extra[k]
	return v, ok
}

func toHiragana(s string) []rune {
	rs := []rune(s)
	for i, r := range rs {
		if r >= 'ァ' && r <= 'ヶ' {
			rs[i] = r - 0x60
		}
	}
	return rs
}

func isVowel(r rune) bool { return strings.ContainsRune("aiueo", r) }

// longPair reports whether vowel b after vowel a is read as one long
// vowel (添え書き 3): the same vowel twice, ou, or ei.
func longPair(a, b rune) bool {
	return a == b || (a == 'o' && b == 'u') || (a == 'e' && b == 'i')
}

const (
	plain      = "aiueoAIUEO"
	macrons    = "āīūēōĀĪŪĒŌ"
	circumflex = "âîûêôÂÎÛÊÔ"
)

func mark(v rune, s Style) rune {
	i := strings.IndexRune(plain, v)
	if i < 0 {
		return v
	}
	if s == Circumflex {
		return []rune(circumflex)[i]
	}
	return []rune(macrons)[i]
}

// Romaji converts hiragana or katakana to romaji under the notice. Anything
// that is not kana passes through unchanged, except 、。　 which become
// , . and a space (添え書き 7). An ASCII or typographic apostrophe in the
// input marks a sound break (添え書き 4): it stops two vowels from merging
// into a long vowel and is written out only where the letters would
// otherwise be misread (おお'おじ → oo'oji, but ōoji with a macron).
func Romaji(kana string, opt Options) string {
	rs := toHiragana(kana)
	e := emitter{opt: opt, wordStart: true}
	for i := 0; i < len(rs); {
		r := rs[i]
		switch {
		case r == '\'' || r == '’':
			e.pendingBreak = true
			i++
		case r == 'っ':
			e.sokuon = true
			i++
		case r == 'ん':
			e.hatsuon()
			i++
		case r == 'ー':
			e.choon()
			i++
		default:
			if i+1 < len(rs) {
				if v, ok := lookup(string(rs[i:i+2]), opt.Kunrei); ok {
					e.mora(v)
					i += 2
					continue
				}
			}
			if v, ok := lookup(string(r), opt.Kunrei); ok {
				e.mora(v)
			} else {
				e.raw(r)
			}
			i++
		}
	}
	return string(e.out)
}

type emitter struct {
	opt          Options
	out          []rune
	sokuon       bool // a っ waiting for its consonant (添え書き 2)
	pendingN     bool // an ん just written; needs ' before a vowel or y (添え書き 4)
	pendingBreak bool // an apostrophe seen in the input
	wordStart    bool
	prevVowel    rune // last vowel letter written and still open for lengthening
}

func (e *emitter) write(r rune) {
	if e.opt.Capitalize && e.wordStart {
		r = unicode.ToUpper(r)
	}
	e.wordStart = false
	e.out = append(e.out, r)
}

func (e *emitter) mora(r string) {
	first := rune(r[0])
	if e.pendingN {
		if isVowel(first) || first == 'y' {
			e.out = append(e.out, '\'')
		}
		e.pendingN, e.pendingBreak = false, false
	}
	if e.pendingBreak {
		if e.prevVowel != 0 && isVowel(first) && longPair(e.prevVowel, first) {
			e.out = append(e.out, '\'')
		}
		e.pendingBreak = false
		e.prevVowel = 0
	}
	if len(r) == 1 && e.prevVowel != 0 && longPair(e.prevVowel, first) {
		e.sokuon = false // っ before a vowel has no letter to double
		e.lengthen(first)
		return
	}
	if e.sokuon {
		if !isVowel(first) {
			e.write(first) // s of sh, c of ch, t of ts: the first letter (添え書き 2)
		}
		e.sokuon = false
	}
	for _, c := range r {
		e.write(c)
	}
	last := rune(r[len(r)-1])
	e.prevVowel = 0
	if isVowel(last) {
		e.prevVowel = last
	}
}

// lengthen merges the vowel v into the vowel just written.
func (e *emitter) lengthen(v rune) {
	switch {
	case e.opt.Style == Doubled:
		e.out = append(e.out, v)
		e.prevVowel = v
	case v == 'i':
		// イ列 and エ列+い keep two letters (添え書き 3, the 〈 〉 forms are
		// listed but "母音字を並べたつづりを用いるのが一般的").
		e.out = append(e.out, 'i')
		e.prevVowel = 0
	default:
		e.out[len(e.out)-1] = mark(e.out[len(e.out)-1], e.opt.Style)
		e.prevVowel = 0
	}
}

func (e *emitter) hatsuon() {
	e.sokuon = false
	e.write('n')
	e.prevVowel = 0
	e.pendingN, e.pendingBreak = true, false
}

// choon handles ー, which has no kana of its own to double.
func (e *emitter) choon() {
	e.pendingN, e.pendingBreak, e.sokuon = false, false, false
	if e.prevVowel == 0 {
		e.out = append(e.out, '-')
		return
	}
	if e.opt.Style == Doubled {
		e.out = append(e.out, e.prevVowel)
		return
	}
	e.out[len(e.out)-1] = mark(e.out[len(e.out)-1], e.opt.Style)
	e.prevVowel = 0
}

func (e *emitter) raw(r rune) {
	e.pendingN, e.pendingBreak, e.sokuon = false, false, false
	switch r {
	case '、':
		r = ','
	case '。':
		r = '.'
	case '　':
		r = ' '
	}
	e.out = append(e.out, r)
	e.prevVowel = 0
	// A new word starts after anything but a letter, a digit, a hyphen
	// (Kutani-yaki, 添え書き 5) or an apostrophe.
	e.wordStart = !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '\'' && r != '’'
}

// reverse maps every romaji spelling this package knows (本表, 1954 第1表,
// loanword kana) back to kana, so Normalize can re-spell it. Where two
// spellings collide the 対照表 row wins: "tu" is つ, not トゥ; "ti" is ち.
var reverse = func() map[string]string {
	m := map[string]string{}
	for _, t := range []map[string]string{gairaigo, honpyo, kunrei} {
		for k, v := range t {
			m[v] = k
		}
	}
	return m
}()

var vowelMarks = func() map[rune]rune {
	m := map[rune]rune{}
	for i, r := range []rune(macrons) {
		m[r] = []rune(plain)[i]
	}
	for i, r := range []rune(circumflex) {
		m[r] = []rune(plain)[i]
	}
	return m
}()

func isConsonant(r rune) bool { return r >= 'a' && r <= 'z' && !isVowel(r) }

// parseWord reads one lower-cased romaji word into kana, or reports false
// if any part of it is not romaji this package knows.
func parseWord(w []rune) ([]rune, bool) {
	// Split marked vowels so that "kō" matches "ko" followed by ー.
	var unmarked []rune
	for _, r := range w {
		if v, ok := vowelMarks[r]; ok {
			unmarked = append(unmarked, v, 'ー')
		} else {
			unmarked = append(unmarked, r)
		}
	}
	w = unmarked
	var out []rune
	for i := 0; i < len(w); {
		c := w[i]
		var next rune
		if i+1 < len(w) {
			next = w[i+1]
		}
		switch {
		case c == '\'' || c == '’':
			out = append(out, '\'')
			i++
		case c == 'm' && (next == 'b' || next == 'm' || next == 'p'):
			out = append(out, 'ん') // Shimbashi, ramma, tempura (添え書き 8)
			i++
		case c == 'n' && (next == 0 || next == '\'' || next == '’' || (!isVowel(next) && next != 'y')):
			out = append(out, 'ん')
			i++
		case c == 't' && next == 'c' && i+2 < len(w) && w[i+2] == 'h':
			out = append(out, 'っ') // matcha → maccha (添え書き 8)
			i++
		case isConsonant(c) && next == c:
			out = append(out, 'っ')
			i++
		case c == 'ー':
			out = append(out, 'ー')
			i++
		default:
			n := 3
			for ; n > 0; n-- {
				if i+n > len(w) {
					continue
				}
				if k, ok := reverse[string(w[i:i+n])]; ok {
					out = append(out, []rune(k)...)
					break
				}
			}
			if n == 0 {
				return nil, false
			}
			i += n
		}
	}
	return out, true
}

func isWordRune(r rune) bool { return unicode.IsLetter(r) || r == '\'' || r == '’' }

// Normalize rewrites romaji into the notice's spelling: 1954 第1表 spellings
// (si → shi, tu → tsu, hu → fu, zi → ji, sya → sha, tya → cha …), m before
// b/m/p (Shimbashi → Shinbashi), tch (matcha → maccha), and long vowels into
// opt.Style. Words that are not entirely romaji this package knows are left
// exactly as they were, so "Line" and "Express" survive. The 1954 第2表
// spellings the notice lists as used in proper nouns (di, du, dya, wo, kwa,
// gwa) are deliberately not rewritten: they coincide with loanword kana and
// the notice itself says to respect the bearer's spelling (添え書き 9).
// With opt.Kunrei the rewrite runs the other way, into 訓令式.
func Normalize(romaji string, opt Options) string {
	rs := []rune(romaji)
	var out []rune
	for i := 0; i < len(rs); {
		if !isWordRune(rs[i]) {
			out = append(out, rs[i])
			i++
			continue
		}
		j := i
		for j < len(rs) && isWordRune(rs[j]) {
			j++
		}
		out = append(out, normalizeWord(rs[i:j], opt)...)
		i = j
	}
	return string(out)
}

func normalizeWord(w []rune, opt Options) []rune {
	lower := make([]rune, len(w))
	for i, r := range w {
		lower[i] = unicode.ToLower(r)
	}
	kana, ok := parseWord(lower)
	if !ok {
		return w
	}
	opt.Capitalize = unicode.IsUpper(w[0])
	return []rune(Romaji(string(kana), opt))
}
