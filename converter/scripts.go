package converter

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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
	var errors []string

	// Strategy 1: Temp directory (preferred - cleaned by OS)
	dir, cleanup, err := extractToTemp()
	if err == nil {
		return &ScriptManager{scriptDir: dir, cleanup: cleanup}, nil
	}
	errors = append(errors, fmt.Sprintf("temp: %v", err))

	// Strategy 2: Persistent cache directory
	dir, err = extractToCache()
	if err == nil {
		return &ScriptManager{scriptDir: dir, cleanup: func() error { return nil }}, nil
	}
	errors = append(errors, fmt.Sprintf("cache: %v", err))

	// Strategy 3: Home directory fallback
	dir, err = extractToHome()
	if err == nil {
		return &ScriptManager{scriptDir: dir, cleanup: func() error { return nil }}, nil
	}
	errors = append(errors, fmt.Sprintf("home: %v", err))

	// Strategy 4: Current working directory (Windows-specific fallback)
	if runtime.GOOS == "windows" {
		dir, err = extractToCWD()
		if err == nil {
			return &ScriptManager{scriptDir: dir, cleanup: func() error { return nil }}, nil
		}
		errors = append(errors, fmt.Sprintf("cwd: %v", err))
	}

	return nil, fmt.Errorf("failed to extract scripts to any location:\n  %s",
		strings.Join(errors, "\n  "))
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

	// Verify all scripts exist
	scripts := []string{"detect_webp_type.py", "webp_to_gif.py", "webp_to_jpeg.py"}
	for _, script := range scripts {
		path := sm.GetScriptPath(script)
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("script %s not found: %w", script, err)
		}
	}

	return nil
}

// extractToTemp creates temporary directory and extracts scripts
func extractToTemp() (string, func() error, error) {
	tmpDir, err := os.MkdirTemp("", "webp2gif-*")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	if err := extractScripts(tmpDir); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to extract scripts to temp: %w", err)
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
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Use different directory naming for Windows vs Unix
	var dirName string
	if runtime.GOOS == "windows" {
		dirName = "webp2gifjpeg-temp" // No leading dot on Windows
	} else {
		dirName = ".webp2gifjpeg-tmp"
	}

	dir := filepath.Join(homeDir, dirName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create home directory: %w", err)
	}

	if err := extractScripts(dir); err != nil {
		return "", err
	}

	return dir, nil
}

// extractToCWD extracts to current working directory (Windows fallback)
func extractToCWD() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	dir := filepath.Join(cwd, ".webp2gifjpeg")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	if err := extractScripts(dir); err != nil {
		return "", err
	}

	return dir, nil
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
				return "", fmt.Errorf("failed to get home directory: %w", err)
			}
			cacheBase = filepath.Join(home, ".cache")
		}
	case "windows":
		// Try multiple Windows environment variables in order of preference
		cacheBase = os.Getenv("LOCALAPPDATA")
		if cacheBase == "" {
			cacheBase = os.Getenv("APPDATA")
		}
		if cacheBase == "" {
			// Fallback to USERPROFILE\AppData\Local
			userProfile := os.Getenv("USERPROFILE")
			if userProfile == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return "", fmt.Errorf("failed to determine user directory: %w", err)
				}
				userProfile = home
			}
			cacheBase = filepath.Join(userProfile, "AppData", "Local")
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

		// Use 0644 for better Windows compatibility
		if err := os.WriteFile(dstPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write %s to %s: %w", script, dstPath, err)
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
