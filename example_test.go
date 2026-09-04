package hebon_test

import (
	"fmt"

	"github.com/doxuta/hebon"
)

func ExampleRomaji() {
	fmt.Println(hebon.Romaji("とうきょう しんばし ぎんざ", hebon.Options{Capitalize: true}))
	fmt.Println(hebon.Romaji("とうきょう しんばし ぎんざ", hebon.Options{Capitalize: true, Style: hebon.Doubled}))
	fmt.Println(hebon.Romaji("とうきょう しんばし ぎんざ", hebon.Options{Capitalize: true, Kunrei: true, Style: hebon.Circumflex}))
	// Output:
	// Tōkyō Shinbashi Ginza
	// Toukyou Shinbashi Ginza
	// Tôkyô Sinbasi Ginza
}

func ExampleNormalize() {
	fmt.Println(hebon.Normalize("Shimbashi tyatto Ohtawara Hukusaki Line", hebon.Options{}))
	// Output:
	// Shinbashi chatto Ohtawara Fukusaki Line
}
