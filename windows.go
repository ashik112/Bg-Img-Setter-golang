// +build windows

package main

import (
	/*"os"*/
)

// Desktop contains the current desktop environment on Linux.
// Empty string on all other operating systems.

func getDesktop() string {
	desktop := "Windows"
	return desktop
}
