package banner

import (
	"fmt"
)

// prints the version message
const version = "v0.0.3"

func PrintVersion() {
	fmt.Printf("Current socialfinder version %s\n", version)
}

// Prints the Colorful banner
func PrintBanner() {
	banner := `
                       _         __ ____ _             __           
   _____ ____   _____ (_)____ _ / // __/(_)____   ____/ /___   _____
  / ___// __ \ / ___// // __  // // /_ / // __ \ / __  // _ \ / ___/
 (__  )/ /_/ // /__ / // /_/ // // __// // / / // /_/ //  __// /    
/____/ \____/ \___//_/ \__,_//_//_/  /_//_/ /_/ \__,_/ \___//_/
`
	fmt.Printf("%s\n%65s\n\n", banner, "Current socialfinder version "+version)
}
