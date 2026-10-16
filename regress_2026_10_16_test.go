package hebon

import "testing"

// Regression tests written 2026-10-16. Each one FAILS on the current code;
// the fix is deliberately left to the maintainer.

// Finding 1: an apostrophe that is not between two kana is silently dropped.
// README ("Library"): "anything that is not kana passes through".
func TestRegress20261016_RomajiKeepsApostropheOutsideKana(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Don't", "Don't"},
		{"I’m going to とうきょう", "I’m going to tōkyō"},
		{"'とうきょう'", "'tōkyō'"}, // a quote before kana is not a break marker
	}
	for _, tc := range cases {
		if got := Romaji(tc.in, Options{}); got != tc.want {
			t.Errorf("Romaji(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// Finding 2: Normalize glues quotes onto the word, then drops them and the
// capital letter behind them. README: a word that is not romaji "is returned
// exactly as it was"; the quotes are not romaji at all.
func TestRegress20261016_NormalizeKeepsQuotesAndCase(t *testing.T) {
	cases := []struct{ in, want string }{
		{"'Tokyo'", "'Tokyo'"},
		{"‘Tokyo’", "‘Tokyo’"},
		{"rock 'n' roll", "rock 'n' roll"},
	}
	for _, tc := range cases {
		if got := Normalize(tc.in, Options{}); got != tc.want {
			t.Errorf("Normalize(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// Finding 3: the particle を is merged with a following う/お into ō.
// をう / をお is never a long vowel in 現代仮名遣い (添え書き 3 (2) ties the
// doubled-letter style to it), and 添え書き 7 makes を the particle "o";
// 添え書き 4 wants an apostrophe where two vowels would be misread as one,
// exactly as the package already writes こ'うた → ko'uta.
func TestRegress20261016_ParticleWoNeverLengthens(t *testing.T) {
	cases := []struct {
		in   string
		opt  Options
		want string
	}{
		{"うたをうたう", Options{}, "utao'utau"},      // 歌を歌う
		{"てがみをおくる", Options{}, "tegamio'okuru"}, // 手紙を送る
		{"うたをうたう", Options{Style: Doubled}, "utao'utau"},
	}
	for _, tc := range cases {
		if got := Romaji(tc.in, tc.opt); got != tc.want {
			t.Errorf("Romaji(%q, %+v) = %q, want %q", tc.in, tc.opt, got, tc.want)
		}
	}
}
