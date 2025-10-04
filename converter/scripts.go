package converter

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

//go:embed scripts
var scriptsFS embed.FS

const scriptVersion = "1.0.0"

// ScriptManager handles extraction and management of embedded Python scripts
type ScriptManager struct {
	scriptDir string
	cleanup   func() error
}

// NewScriptManager creates a new script manager with fallback extraction strategy
func NewScriptManager() (*ScriptManager, error) {
	// Strategy 1: Temp directory (preferred - cleaned by OS)
	dir, cleanup, err := extractToTemp()
	if err == nil {
		return &ScriptManager{scriptDir: dir, cleanup: cleanup}, nil
	}

	// Strategy 2: Persistent cache directory
	dir, err = extractToCache()
	if err == nil {
		return &ScriptManager{scriptDir: dir, cleanup: func() error { return nil }}, nil
	}

	// Strategy 3: Home directory fallback
	dir, err = extractToHome()
	if err == nil {
		return &ScriptManager{scriptDir: dir, cleanup: func() error { return nil }}, nil
	}

	return nil, fmt.Errorf("failed to extract scripts to any location")
}

// GetScriptPath returns the full path to a named script
func (sm *ScriptManager) GetScriptPath(name string) string {
	return filepath.Join(sm.scriptDir, name)
}

// Cleanup removes temporary files if applicable
func (sm *ScriptManager) Cleanup() error {
	if sm.cleanup != nil {
		return sm.cleanup()
	}
	return nil
}

// Validate checks if Python and Pillow are available
func (sm *ScriptManager) Validate() error {
	// Check Python availability
	if _, err := exec.LookPath("python3"); err != nil {
		return fmt.Errorf(`python3 not found in PATH

This tool requires Python 3.x to be installed.

Installation instructions:
  Ubuntu/Debian: sudo apt install python3
  macOS:         brew install python3
  Windows:       Download from python.org

After installing Python, also install Pillow:
  pip3 install Pillow`)
	}

	// Check Pillow availability
	cmd := exec.Command("python3", "-c", "import PIL")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(`Pillow library not installed

This tool requires the Pillow library for image processing.

Installation:
  pip3 install Pillow

On some systems you may need:
  pip3 install --break-system-packages Pillow`)
	}

	// Verify all scripts exist and are executable
	scripts := []string{"detect_webp_type.py", "webp_to_gif.py", "webp_to_jpeg.py"}
	for _, script := range scripts {
		path := sm.GetScriptPath(script)
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("script %s not found: %w", script, err)
		}
		// Check if executable (on Unix-like systems)
		if runtime.GOOS != "windows" && info.Mode().Perm()&0100 == 0 {
			return fmt.Errorf("script %s is not executable", script)
		}
	}

	return nil
}

// extractToTemp creates temporary directory and extracts scripts
func extractToTemp() (string, func() error, error) {
	tmpDir, err := os.MkdirTemp("", "webp2gif-*")
	if err != nil {
		return "", nil, err
	}

	if err := extractScripts(tmpDir); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}

	cleanup := func() error {
		return os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup, nil
}

// extractToCache uses persistent cache directory with versioning
func extractToCache() (string, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return "", err
	}

	// Check if scripts already exist and are current version
	versionFile := filepath.Join(cacheDir, ".version")
	if currentVersion, _ := os.ReadFile(versionFile); string(currentVersion) == scriptVersion {
		// Verify scripts exist
		if scriptsExist(cacheDir) {
			return cacheDir, nil
		}
	}

	// Extract scripts
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	if err := extractScripts(cacheDir); err != nil {
		return "", err
	}

	// Write version file
	if err := os.WriteFile(versionFile, []byte(scriptVersion), 0644); err != nil {
		return "", err
	}

	return cacheDir, nil
}

// extractToHome fallback to home directory
func extractToHome() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(homeDir, ".webp2gifjpeg-tmp")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return dir, extractScripts(dir)
}

// getCacheDir returns platform-specific cache directory
func getCacheDir() (string, error) {
	var cacheBase string

	switch runtime.GOOS {
	case "linux", "darwin":
		if xdgCache := os.Getenv("XDG_CACHE_HOME"); xdgCache != "" {
			cacheBase = xdgCache
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			cacheBase = filepath.Join(home, ".cache")
		}
	case "windows":
		cacheBase = os.Getenv("LOCALAPPDATA")
		if cacheBase == "" {
			return "", fmt.Errorf("LOCALAPPDATA not set")
		}
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return filepath.Join(cacheBase, "webp2gifjpeg"), nil
}

// extractScripts writes all embedded scripts to target directory
func extractScripts(targetDir string) error {
	scripts := []string{"detect_webp_type.py", "webp_to_gif.py", "webp_to_jpeg.py"}

	for _, script := range scripts {
		srcPath := filepath.Join("scripts", script)
		dstPath := filepath.Join(targetDir, script)

		data, err := scriptsFS.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("failed to read embedded %s: %w", script, err)
		}

		if err := os.WriteFile(dstPath, data, 0755); err != nil {
			return fmt.Errorf("failed to write %s: %w", script, err)
		}
	}

	return nil
}

// scriptsExist checks if all required scripts exist in directory
func scriptsExist(dir string) bool {
	scripts := []string{"detect_webp_type.py", "webp_to_gif.py", "webp_to_jpeg.py"}
	for _, script := range scripts {
		if _, err := os.Stat(filepath.Join(dir, script)); err != nil {
			return false
		}
	}
	return true
}
