// +build linux

package main

import (
	"os/exec"
	"strings"

	"os/user"
	"path/filepath"

	"github.com/go-ini/ini"
	"strconv"
)

func getDesktop() string {
	//desktop := os.Getenv("XDG_CURRENT_DESKTOP")
	desktop := "Linux"
	return desktop
}

// Get returns the current wallpaper.
func Get() (string, error) {
	switch Desktop {
	case "GNOME", "Unity", "Pantheon", "Budgie:GNOME":
		return parseDconf("dconf", "read", "/org/gnome/desktop/background/picture-uri")
	/*case "KDE":
	return parseKDEConfig()*/
	case "X-Cinnamon":
		return parseDconf("dconf", "read", "/org/cinnamon/desktop/background/picture-uri")
	case "MATE":
		return parseDconf("dconf", "read", "/org/mate/desktop/background/picture-filename")
	case "XFCE":
		output, err := exec.Command("xfconf-query", "-c", "xfce4-desktop", "-p", "/backdrop/screen0/monitor0/workspace0/last-image").Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(output)), nil
	case "LXDE":
		return parseLXDEConfig()
	case "Deepin":
		return parseDconf("dconf", "read", "/com/deepin/wrap/gnome/desktop/background/picture-uri")
	default:
		return "", ErrUnsupportedDE
	}
}

// SetFromFile sets wallpaper from a file path.
func SetFromFile(file string) error {
	switch Desktop {
	case "GNOME", "Unity", "Pantheon", "Budgie:GNOME":
		return exec.Command("dconf", "write", "/org/gnome/desktop/background/picture-uri", strconv.Quote("file://"+file)).Run()
	/*case "KDE":
	return setKDEBackground("file://" + file)*/
	case "X-Cinnamon":
		return exec.Command("dconf", "write", "/org/cinnamon/desktop/background/picture-uri", strconv.Quote("file://"+file)).Run()
	case "MATE":
		return exec.Command("dconf", "write", "/org/mate/desktop/background/picture-filename", strconv.Quote(file)).Run()
	case "XFCE":
		return exec.Command("xfconf-query", "-c", "xfce4-desktop", "-p", "/backdrop/screen0/monitor0/workspace0/last-image", "-s", file).Run()
	case "LXDE":
		return exec.Command("pcmanfm", "-w", file).Run()
	case "Deepin":
		return exec.Command("dconf", "write", "/com/deepin/wrap/gnome/desktop/background/picture-uri", strconv.Quote("file://"+file)).Run()
	default:
		return ErrUnsupportedDE
	}
}

/*func getCacheDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".cache"), nil
}*/

func unquote(quoted string) string {
	if len(quoted) >= 2 && quoted[0] == '\'' && quoted[len(quoted)-2] == '\'' {
		return quoted[1 : len(quoted)-2]
	}

	return quoted
}

func removeProtocol(output string) string {
	if len(output) >= 7 && output[:7] == "file://" {
		return output[7:]
	}

	return output
}

func parseDconf(command string, args ...string) (string, error) {
	output, err := exec.Command(command, args...).Output()
	if err != nil {
		return "", err
	}

	return removeProtocol(unquote(string(output))), nil
}

func parseLXDEConfig() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	cfg, err := ini.Load(filepath.Join(usr.HomeDir, ".config/pcmanfm/LXDE/desktop-items-0.conf"))
	if err != nil {
		return "", err
	}

	key, err := cfg.Section("*").GetKey("wallpaper")
	if err != nil {
		return "", err
	}
	return key.String(), err
}
