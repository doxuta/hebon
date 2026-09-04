# hebon ヘボン

Kana → romaji the way **令和7年内閣告示第4号「ローマ字のつづり方」** (22 December
2025) spells it, and a normaliser that rewrites 1954-style romaji (訓令式
`si ti tu hu zi sya…`, `Shimbashi`, `matcha`) into that spelling. Pure Go,
standard library only, one package and one command.

```
$ hebon -cap とうきょう しんばし ぎんざ にっちょく
Tōkyō Shinbashi Ginza Nicchoku

$ hebon -cap -kunrei -style circumflex とうきょう しんばし ぎんざ にっちょく
Tôkyô Sinbasi Ginza Nittyoku

$ hebon -n "Shimbashi tyatto Ohtawara Hukusaki Line"
Shinbashi chatto Ohtawara Fukusaki Line
```

## Why

On 22 December 2025 the Cabinet replaced 昭和29年内閣告示第1号 — the 1954
notice that put 訓令式 in its first table and ヘボン式 in its second — with a
single Hepburn-based 本表. It is the first change to how Japan spells itself in
Latin letters in 71 years; textbooks switch in 2026年度. The notice is short: a
kana table, nine 添え書き (ん, っ, long vowels, the apostrophe, hyphens,
capitals, は/へ/を, existing spellings, personal names) and a 22-row 対照表
against both 1954 tables.

Nobody had put it in code. A GitHub search for 「ローマ字のつづり方」 finds no
repository; the one Go kana library with stars, [`gojp/kana`](https://github.com/gojp/kana),
stopped in 2020 and handles none of ん-before-vowel, long vowels or the
apostrophe the way the notice does. Meanwhile, in the last three months:

- [`TrainLCD/StationAPI#1653`](https://github.com/TrainLCD/StationAPI/issues/1653) —
  a production transit app fixing `Nirehana`→`Nirehara`, `Fukusai`→`Fukusaki`
  by hand, then running a "ローマ字変換の一括検証" over the dataset.
- [`shimatoshi/transit-pwa#12`](https://github.com/shimatoshi/transit-pwa/issues/12) —
  2,464 of 9,901 stations with no romaji at all; proposed fix: "駅名からのヘボン式変換".
- [`finelagusaz/Snotra#523`](https://github.com/finelagusaz/Snotra/issues/523) —
  a search that missed `chatto` against an entry named `tyatto`.
- [`Iktahana/illusions#1729`](https://github.com/Iktahana/illusions/issues/1729) —
  a proofreading rule for 訓令式/ヘボン式 mixing, "実装可能", unimplemented.

`hebon.Romaji` is the first two; `hebon.Normalize` is the last two.

## Install

```
go install github.com/doxuta/hebon/cmd/hebon@latest
```

or, without installing anything:

```
go run github.com/doxuta/hebon/cmd/hebon@latest -cap とうきょう しんばし ぎんざ
```

The command takes text as arguments or one line at a time on stdin.

| Flag | Meaning |
|---|---|
| `-style macron` (default) | long vowels with a macron: `kāsan` — 添え書き 3(1) |
| `-style circumflex` | `kâsan` — 3(1) allows `^` where a macron is unavailable |
| `-style doubled` | `kaasan`, letters follow the kana — 3(2) |
| `-kunrei` | spell with the 1954 第1表 instead, to see what changed |
| `-cap` | capitalise the first letter of each word — 添え書き 6, 7 |
| `-n` | input is romaji: rewrite it to the notice's spelling |

## Library

```go
import "github.com/doxuta/hebon"

hebon.Romaji("とうきょう しんばし ぎんざ", hebon.Options{Capitalize: true})
// Tōkyō Shinbashi Ginza
hebon.Romaji("とうきょう", hebon.Options{Style: hebon.Doubled})
// toukyou
hebon.Romaji("しんじゅく", hebon.Options{Kunrei: true})
// sinzyuku

hebon.Normalize("Shimbashi tyatto Ohtawara Line", hebon.Options{})
// Shinbashi chatto Ohtawara Line
hebon.Normalize("Shinjuku", hebon.Options{Kunrei: true})
// Sinzyuku
```

`Romaji` folds katakana to hiragana and looks each mora up in the 本表; anything
that is not kana passes through. `Normalize` splits a line into words, reads
each word as romaji under the union of the 本表, the 1954 第1表 and the loanword
kana, and re-spells it; a word that does not parse completely is returned
exactly as it was, which is how `Line` and `Ohtawara` survive above.

## What the notice says, and what this does

| 添え書き | Rule | Here |
|---|---|---|
| 1 | 撥音 ん is `n` (`shinbun`, not `shimbun`) | always `n`; `Normalize` turns `m` before b/m/p back into `n` |
| 2 | 促音 っ doubles the next consonant; for `sh`/`ch` the first letter (`zasshi`, `nicchoku`) | first letter of the next mora; `Normalize` reads `tch` as っ |
| 3 | long vowels: macron, or `^` if needed, or doubled letters as in 現代仮名遣い; イ列 and エ列+い are usually doubled (`kawaii`, `Heisei`) | `Macron` / `Circumflex` / `Doubled`; `ii` and `ei` stay two letters in every style |
| 4 | `'` separates ん from a following vowel or y (`tan'i`, `ton'ya`) and marks a non-long vowel pair (`oo'oji`, `ko'uta`) | automatic after ん; an apostrophe in the kana input marks the break and is written only where the letters would otherwise merge (`oo'oji` doubled, `ōoji` with a macron) |
| 5 | `-` may join the parts of a compound (`Kutani-yaki`) | hyphens pass through and do not start a new word for `Capitalize` |
| 6, 7 | proper nouns and sentences start with a capital; `，` and `．`; は/へ/を are `wa`/`e`/`o` | `Capitalize`; 、。 become `,` `.`; を is `o` from the 本表 — は and へ are **not** detected, see Limitations |
| 8 | spellings from other fields (`Tokyo`, `judo`, `Ohtawara`, `Shimbashi`, `tempura`, `matcha`) are not forced to change | `Normalize` rewrites the last three; the first three cannot be recovered from the short spelling and are left alone |
| 9 | personal and organisation names follow the bearer's wish | the 1954 第2表 spellings the notice names as used in proper nouns (`di du dya dyu dyo wo kwa gwa`) are never rewritten |
| 対照表 | 22 rows against both 1954 tables | `Kunrei: true` spells with the first; `Normalize` maps its spellings to the 本表 and, with `Kunrei`, back |

Tests pin all 86 spellings printed in the notice's own examples (43 words in
both long-vowel styles), all 22 rows of the 対照表 in both directions, the
four threads above, and — by fuzzing — that `Normalize` is idempotent and that
`Romaji` never leaks a っ.

## Limitations

- **Kanji are not read.** Input is kana. The transit datasets above carry a
  読み仮名 column; for free text use a morphological analyser such as
  [`ikawaha/kagome`](https://github.com/ikawaha/kagome) first.
- **Particles は and へ are not detected** — that needs grammar, not a table.
  Write わ and え in the kana, or post-process. を is always `o`, which is what
  the 本表 says.
- **Long vowels need the kana.** おう always merges to `ō`; a genuine two-vowel
  boundary (大伯父 おお'おじ, 小唄 こ'うた) has to be marked with an apostrophe
  in the input, exactly as the notice does in print. `Normalize` on a macron
  spelling cannot recover the kana, so `Tōkyō` in `Doubled` style becomes
  `Tookyoo`, not the notice's `Toukyou`.
- **ー in イ列 and エ列** becomes `ī`/`ē` under `Macron` (コーヒー → `kōhī`) even
  though the notice prefers doubled letters for native words in those rows —
  there is no kana to double. `Doubled` gives `koohii`.
- **Loanword kana** (ティ, ファ, ヴ, ウィ, トゥ …) are outside the notice
  (前書き 4). The spellings here follow common Hepburn usage by analogy, taken
  from the kana lists in 外来語の表記 (平成3年内閣告示第2号) 第1表 and 第2表.
  In `Normalize`, `tu`, `ti` and `tyu` are read as the 対照表 rows つ, ち, ちゅ,
  so `tisshu` becomes `chisshu`.
- `Normalize` does not know `oh` for a long o (`Ohtawara`), `nn` for ん before a
  vowel in typing style, or `l`/`x`/`q` — such words are left untouched.
- Half-width katakana (ｶﾀｶﾅ) are not folded; normalise with NFKC first.
- っ with nothing after it is dropped; ー with nothing before it is `-`;
  ヵ/ヶ pass through; the apostrophe written is ASCII `'`, where the notice
  prints `’`.

## Prior art and sources

- [`gojp/kana`](https://github.com/gojp/kana) (Go, 123★, 2013–2020) by Herman
  Schaaf and Shawn Smith — trie-based kana↔romaji; read before writing this.
  [`chanyeinthaw/kuroshiro.go`](https://github.com/chanyeinthaw/kuroshiro.go) and
  [`deelawn/wanakana`](https://github.com/deelawn/wanakana) are Go ports of the
  JS libraries [kuroshiro](https://github.com/hexenq/kuroshiro) and
  [WanaKana](https://github.com/WaniKani/WanaKana), whose
  [#65](https://github.com/WaniKani/WanaKana/issues/65) has asked for a 訓令式
  option since 2017.
- The notice: [「ローマ字のつづり方」(文化庁)](https://www.bunka.go.jp/kokugo_nihongo/sisaku/joho/joho/kijun/naikaku/roma/index2.html)
  and the [text as circulated to the 言語資源小委員会 on 2025-12-23 (PDF)](https://www.bunka.go.jp/seisaku/bunkashingikai/kokugo/gengo/gengo_10/pdf/94304301_05.pdf),
  which is the version read for this package. Loanword kana:
  [外来語の表記 (文化庁)](https://www.bunka.go.jp/kokugo_nihongo/sisaku/joho/joho/kijun/naikaku/gairai/index.html).

Built with an AI coding agent (Claude) under human review, as with the other
repositories on this account. The notice was read in full and every example it
prints is a test; the loanword table is the one part that is not the notice's
text, and it says so.

## 日本語要約

令和7年内閣告示第4号「ローマ字のつづり方」に従って仮名をローマ字に直す Go
ライブラリと CLI です。本表・添え書き（撥音 n、促音の子音重ね、長音の
「¯」「^」「母音字を並べる」の三方式、切れ目の「’」、「-」、大文字）・対照表を
そのまま実装し、告示に印刷された用例 43 語（両方式で 86 通り）と対照表 22 行を
テストで固定しています。`Normalize` は訓令式や `Shimbashi`・`matcha` のような
慣用表記を告示のつづりに書き換え、`Kunrei` で逆方向にも動きます。漢字の読み、
助詞「は」「へ」の判定、長音かどうかの語ごとの判断はしません（仮名側で
「’」を入れてください）。外来語の仮名（ティ・ファ・ヴ等）は告示の範囲外で、
外来語の表記の表に基づく慣用のつづりです。

## License

MIT © Xuan Tai Doan
