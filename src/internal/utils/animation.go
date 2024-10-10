package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// ShowLoadingAnimation displays a loading animation in the console
func ShowLoadingAnimation(done chan bool) {
	words := []string{"Thinking"} // Add more words to the list to pick randomly
	loadingText := words[rand.Intn(len(words))]
	shades := []string{
		"\033[1;90m", // Light Dark Gray
		"\033[1;37m", // Light Gray
		"\033[0;37m", // White
	}
	resetColor := "\033[0m"

	for {
		select {
		case <-done:
			// Clear the animation before stopping
			fmt.Print("\r\033[K")
			return
		default:
			// Create a wave effect by iterating over each character and applying shades
			for waveStart := 0; waveStart < len(loadingText)+len(shades); waveStart++ {
				fmt.Printf("\r")
				for i := 0; i < len(loadingText); i++ {
					shadeOffset := waveStart - i
					if shadeOffset >= 0 && shadeOffset < len(shades) {
						fmt.Printf("%s%c%s", shades[len(shades)-1-shadeOffset], loadingText[i], resetColor)
					} else {
						fmt.Printf("%s%c%s", shades[0], loadingText[i], resetColor)
					}
				}
				time.Sleep(150 * time.Millisecond)
			}
		}
	}
}
