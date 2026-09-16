package hebon

import (
	"strings"
	"testing"
	"unicode"
)

// Every worked example printed in the 添え書き of 令和7年内閣告示第4号,
// with the 本表 spelling under both long-vowel styles. The Doubled column
// is the (2) column of 添え書き 3 where the notice prints one; elsewhere
// the two styles coincide.
var soegaki = []struct {
	item                  int
	kana, macron, doubled string
}{
	{1, "あんまん", "anman", ""},
	{1, "かんぱい", "kanpai", ""},
	{1, "ぎんざ", "Ginza", ""},
	{1, "しんぶん", "shinbun", ""},
	{2, "ざっし", "zasshi", ""},
	{2, "てっぱん", "teppan", ""},
	{2, "にっちょく", "nicchoku", ""},
	{2, "やっきょく", "yakkyoku", ""},
	{3, "かあさん", "kāsan", "kaasan"},
	{3, "まあ", "mā", "maa"},
	{3, "かわいい", "kawaii", "kawaii"},
	{3, "しいたけ", "shiitake", "shiitake"},
	{3, "にいさん", "niisan", "niisan"},
	{3, "じゅうごや", "jūgoya", "juugoya"},
	{3, "ふうりゅう", "fūryū", "fuuryuu"},
	{3, "ええ", "ē", "ee"},
	{3, "ねえさん", "nēsan", "neesan"},
	{3, "やじろべえ", "yajirobē", "yajirobee"},
	{3, "ていえん", "teien", "teien"},
	{3, "とけいだい", "tokeidai", "tokeidai"},
	{3, "へいせい", "Heisei", "Heisei"},
	{3, "おおかみ", "ōkami", "ookami"},
	{3, "ほおずき", "hōzuki", "hoozuki"},
	{3, "とうほく", "Tōhoku", "Touhoku"},
	{3, "ぼうそう", "Bōsō", "Bousou"},
	{3, "おおどうぐ", "ōdōgu", "oodougu"},
	{3, "こおりどうふ", "kōridōfu", "kooridoufu"},
	{4, "たんい", "tan'i", ""},
	{4, "せんいん", "sen'in", ""},
	{4, "えんゆうかい", "en'yūkai", "en'yuukai"},
	{4, "とんや", "ton'ya", ""},
	{4, "おお'おじ", "ōoji", "oo'oji"},
	{4, "こ'うた", "ko'uta", ""},
	{5, "くたに-やき", "Kutani-yaki", ""},
	{5, "たなか-さん", "Tanaka-san", ""},
	{5, "しち-ご-さん", "shichi-go-san", ""},
	{8, "じゅうどう", "jūdō", "juudou"},
	{8, "とうきょう", "Tōkyō", "Toukyou"},
	{8, "おおたわら", "Ōtawara", "Ootawara"},
	{8, "しんばし", "Shinbashi", ""},
	{8, "らんま", "ranma", ""},
	{8, "てんぷら", "tenpura", ""},
	{8, "まっちゃ", "maccha", ""},
}

func TestSoegakiExamples(t *testing.T) {
	n := 0
	for _, tc := range soegaki {
		capitalized := unicode.IsUpper([]rune(tc.macron)[0])
		got := Romaji(tc.kana, Options{Style: Macron, Capitalize: capitalized})
		if got != tc.macron {
			t.Errorf("添え書き %d: Romaji(%q, Macron) = %q, want %q", tc.item, tc.kana, got, tc.macron)
		}
		n++
		want := tc.doubled
		if want == "" {
			want = tc.macron
		}
		got = Romaji(tc.kana, Options{Style: Doubled, Capitalize: capitalized})
		if got != want {
			t.Errorf("添え書き %d: Romaji(%q, Doubled) = %q, want %q", tc.item, tc.kana, got, want)
		}
		n++
	}
	if n != 86 {
		t.Errorf("checked %d spellings, the README says 86", n)
	}
}

// The (付)対照表: 本表 spelling, 1954 第1表 spelling, 1954 第2表 spelling.
var taishohyo = []struct{ kana, honpyo, dai1, dai2 string }{
	{"シ", "shi", "si", "shi"},
	{"チ", "chi", "ti", "chi"},
	{"ツ", "tsu", "tu", "tsu"},
	{"フ", "fu", "hu", "fu"},
	{"ヲ", "o", "o", "wo"},
	{"ジ", "ji", "zi", "ji"},
	{"ヂ", "ji", "zi", "di"},
	{"ヅ", "zu", "zu", "du"},
	{"シャ", "sha", "sya", "sha"},
	{"シュ", "shu", "syu", "shu"},
	{"ショ", "sho", "syo", "sho"},
	{"チャ", "cha", "tya", "cha"},
	{"チュ", "chu", "tyu", "chu"},
	{"チョ", "cho", "tyo", "cho"},
	{"ジャ", "ja", "zya", "ja"},
	{"ジュ", "ju", "zyu", "ju"},
	{"ジョ", "jo", "zyo", "jo"},
	{"ヂャ", "ja", "zya", "dya"},
	{"ヂュ", "ju", "zyu", "dyu"},
	{"ヂョ", "jo", "zyo", "dyo"},
	{"カ", "ka", "ka", "kwa"},
	{"ガ", "ga", "ga", "gwa"},
}

func TestTaishohyo(t *testing.T) {
	if len(taishohyo) != 22 {
		t.Fatalf("対照表 has 22 rows, table has %d", len(taishohyo))
	}
	for _, row := range taishohyo {
		if got := Romaji(row.kana, Options{}); got != row.honpyo {
			t.Errorf("Romaji(%q) = %q, want 本表 %q", row.kana, got, row.honpyo)
		}
		if got := Romaji(row.kana, Options{Kunrei: true}); got != row.dai1 {
			t.Errorf("Romaji(%q, Kunrei) = %q, want 第1表 %q", row.kana, got, row.dai1)
		}
		// Every 第1表 spelling normalises to the 本表 …
		if got := Normalize(row.dai1, Options{}); got != row.honpyo {
			t.Errorf("Normalize(%q) = %q, want %q", row.dai1, got, row.honpyo)
		}
		// … and back again with Kunrei.
		if got := Normalize(row.honpyo, Options{Kunrei: true}); got != row.dai1 {
			t.Errorf("Normalize(%q, Kunrei) = %q, want %q", row.honpyo, got, row.dai1)
		}
	}
}

func TestNormalize(t *testing.T) {
	cases := []struct{ in, want string }{
		// 添え書き 8: spellings from other fields, and what the notice writes.
		{"Shimbashi", "Shinbashi"},
		{"ramma", "ranma"},
		{"tempura", "tenpura"},
		{"matcha", "maccha"},
		// judo, Tokyo, Ohtawara: the notice's long vowels cannot be recovered
		// from the short spelling, and "oh" is not a spelling this package
		// knows, so they pass through untouched.
		{"judo Tokyo Ohtawara", "judo Tokyo Ohtawara"},
		// finelagusaz/Snotra#523: 訓令式 and ヘボン式 must meet.
		{"tyatto", "chatto"},
		{"Nirehara Hukusaki Kōzimati", "Nirehara Fukusaki Kōjimachi"},
		// Words that are not romaji stay as they are, including inside a line.
		{"Abukuma Kyūkō Line", "Abukuma Kyūkō Line"},
		{"Keihin-Tôhoku sen", "Keihin-Tōhoku sen"},
		// Long-vowel notation follows Style.
		{"Toukyou", "Tōkyō"},
		// n before a vowel keeps its apostrophe; nn before a vowel is ん+な行.
		{"tan'i sen’in konnichiwa", "tan'i sen'in konnichiwa"},
		// 第2表 proper-noun spellings and loanword kana are left alone.
		{"Kwansei disuku wonda", "Kwansei disuku wonda"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := Normalize(tc.in, Options{}); got != tc.want {
			t.Errorf("Normalize(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
	if got := Normalize("Toukyou", Options{Style: Doubled}); got != "Toukyou" {
		t.Errorf("Normalize(Toukyou, Doubled) = %q", got)
	}
}

func TestGairaigo(t *testing.T) {
	cases := []struct{ kana, want string }{
		{"ティッシュ", "tisshu"},
		{"ファックス", "fakkusu"},
		{"ヴァイオリン", "vaiorin"},
		{"ウィスキー", "wisukī"},
		{"ジェット", "jetto"},
		{"デュエット", "dyuetto"},
		{"トゥールーズ", "tūrūzu"},
		{"シェフ", "shefu"},
	}
	for _, tc := range cases {
		if got := Romaji(tc.kana, Options{}); got != tc.want {
			t.Errorf("Romaji(%q) = %q, want %q", tc.kana, got, tc.want)
		}
	}
	if got := Romaji("ウィスキー", Options{Style: Doubled}); got != "wisukii" {
		t.Errorf("Doubled ー = %q", got)
	}
}

func TestStylesAndEdges(t *testing.T) {
	cases := []struct {
		kana string
		opt  Options
		want string
	}{
		{"とうきょう", Options{Style: Circumflex, Capitalize: true}, "Tôkyô"},
		{"とうきょう", Options{Kunrei: true, Style: Circumflex}, "tôkyô"},
		{"しんじゅく", Options{Kunrei: true}, "sinzyuku"},
		{"まっちゃ", Options{Kunrei: true}, "mattya"},
		{"ざっし", Options{Kunrei: true}, "zassi"},
		// っ with nothing to double is dropped; ー with nothing to lengthen is a dash.
		{"あっ", Options{}, "a"},
		{"ーい", Options{}, "-i"},
		// Katakana folds to the same table; punctuation follows 添え書き 7.
		{"トウキョウ、オオサカ。", Options{Capitalize: true}, "Tōkyō,Ōsaka."},
		{"わたし　は　がくせい　です", Options{}, "watashi ha gakusei desu"},
		// ん before y and before a vowel, at a word start, and at the end.
		{"きんようび", Options{}, "kin'yōbi"},
		{"んあ", Options{Capitalize: true}, "N'a"},
		{"ほん", Options{}, "hon"},
		// Non-kana passes through and ends any pending state.
		{"東京 とうきょう", Options{}, "東京 tōkyō"},
		{"", Options{}, ""},
		// Three vowels: only the first pair merges.
		{"おおお", Options{}, "ōo"},
		{"おおお", Options{Style: Doubled}, "ooo"},
		// Small kana alone and historical kana.
		{"ぁゐゑ", Options{}, "aie"},
		// ヵ/ヶ are not rows of the 本表: in 三ヶ月, 茅ヶ崎 they abbreviate 箇
		// and are read か/が, not the small kana their codepoints are named
		// after. They pass through, unfolded (README, Limitations).
		{"三ヶ月", Options{}, "三ヶ月"},
		{"ヵヶ", Options{}, "ヵヶ"},
		// ヴ is the last katakana that does fold.
		{"ヴ", Options{}, "vu"},
	}
	for _, tc := range cases {
		if got := Romaji(tc.kana, tc.opt); got != tc.want {
			t.Errorf("Romaji(%q, %+v) = %q, want %q", tc.kana, tc.opt, got, tc.want)
		}
	}
}

// Everything the 本表 and the 1954 table emit is lower-case ASCII letters,
// and every spelling round-trips through Normalize unchanged.
func TestTablesRoundTrip(t *testing.T) {
	// The three loanword spellings that coincide with a 1954 第1表 row are
	// read as that row (tu → つ, ti → ち, tyu → ちゅ); the README says so.
	collisions := map[string]bool{"tu": true, "ti": true, "tyu": true}
	for _, tbl := range []map[string]string{honpyo, kunrei, gairaigo} {
		for kana, romaji := range tbl {
			for _, c := range romaji {
				if c < 'a' || c > 'z' {
					t.Errorf("%q → %q is not ASCII lower-case", kana, romaji)
				}
			}
			if collisions[romaji] && gairaigo[kana] == romaji {
				continue
			}
			if got := Normalize(romaji, Options{}); got != Romaji(kana, Options{}) {
				t.Errorf("Normalize(%q) = %q, Romaji(%q) = %q", romaji, got, kana, Romaji(kana, Options{}))
			}
		}
	}
}

func FuzzNormalizeIdempotent(f *testing.F) {
	for _, s := range []string{"Shimbashi", "tyatto", "Tōkyō", "oo'oji", "konnichiwa", "n'", "tch", "Ōsaka", "xyz", "ティッシュ"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		once := Normalize(s, Options{})
		if twice := Normalize(once, Options{}); twice != once {
			t.Fatalf("Normalize is not idempotent on %q: %q then %q", s, once, twice)
		}
		if r := Romaji(s, Options{Capitalize: true}); strings.ContainsRune(r, 'っ') {
			t.Fatalf("Romaji(%q) leaked っ: %q", s, r)
		}
	})
}
