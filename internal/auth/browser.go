package auth

import (
	"os"
	"os/exec"
	"runtime"
)

// Headless reports whether we appear to be in a non-interactive environment
// where auto-opening a browser is not possible or useful.
func Headless() bool {
	if os.Getenv("MIDAZ_NO_BROWSER") != "" {
		return true
	}
	if os.Getenv("SSH_CONNECTION") != "" || os.Getenv("SSH_CLIENT") != "" {
		return true
	}
	if runtime.GOOS == "linux" && os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		return true
	}
	if os.Getenv("CI") != "" {
		return true
	}
	return false
}

// OpenBrowser best-effort opens the given URL in the default browser. Returns
// an error if the platform-specific opener fails. Never blocks.
func OpenBrowser(url string) error {
	if Headless() {
		return errNoBrowser
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

// errNoBrowser is returned when Headless() suppresses opening a browser.
var errNoBrowser = browserErr("headless environment — browser not opened")

type browserErr string

func (e browserErr) Error() string { return string(e) }
