package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/getlantern/systray"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/sys/windows/registry"
	"gopkg.in/ini.v1"
)

var appVersion = "dev-build"

//go:embed active.ico
var activeIcon []byte

//go:embed inactive.ico
var inactiveIcon []byte

var gameTemplates []GameTemplate

type Config struct {
	RetroarchPath           string            `ini:"retroarch_path"`
	SavePath                string            `ini:"save_path"`
	SaveToOneFile           bool              `ini:"save_to_one_file"`
	Autorun                 bool              `ini:"autorun"`
	OutputToFiles           bool              `ini:"output_to_files"`
	WebPort                 int               `ini:"web_port"`
	SystemIcon              int               `ini:"system_icon"`
	Theme                   string            `ini:"theme"`
	Language                string            `ini:"language"`
	ThumbnailsPath          string            `ini:"thumbnails_path"`
	EnableThumbnails        bool              `ini:"enable_thumbnails"`
	ThumbnailSize           string            `ini:"thumbnail_size"`
	AlternateThumbnails     bool              `ini:"alternate_thumbnails"`
	ThumbnailSwitchInterval int               `ini:"thumbnail_switch_interval"`
	FadeDuration            float64           `ini:"fade_duration"`
	FadeType                string            `ini:"fade_type"`
	Systems                 map[string]string `ini:"systems"`
}
type GameTemplate struct {
	ProcessName  string `json:"process_name"`
	WindowTitle  string `json:"window_title"`
	System       string `json:"system"`
	Game         string `json:"game"`
	NamedTitles  string `json:"named_titles"`
	NamedBoxarts string `json:"named_boxarts"`
	pid          int32
	isRunning    bool
}

type Language struct {
	Code string
	Name string
}
type Translations map[string]string

var (
	currentGame    string
	currentConsole string
	configMutex    sync.RWMutex
	systemsPath    string
	themePath      string
	langPath       string
	config         Config
	templates      map[string]*template.Template
	translations   Translations
)

type ClientInfo struct {
	Conn   *websocket.Conn
	Screen string
}

var clients = make(map[*websocket.Conn]ClientInfo)
var clientsMutex sync.Mutex

func isRetroarchRunning() (bool, int32) {
	processes, err := process.Processes()
	if err != nil {
		log.Printf("Error getting process list: %v", err)
		return false, 0
	}
	for _, p := range processes {
		name, err := p.Name()
		if err == nil && strings.ToLower(name) == "retroarch.exe" {
			return true, p.Pid
		}
	}
	return false, 0
}
func setAutorun(enable bool, appName string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer func() {
		if err := key.Close(); err != nil {
			log.Printf("failed close: %v", err)
		}
	}()

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	if enable {
		err = key.SetStringValue(appName, exePath)
		if err != nil {
			return err
		}
		log.Printf("Program added to autorun: %s", exePath)
	} else {
		err = key.DeleteValue(appName)
		if err != nil && !errors.Is(err, registry.ErrNotExist) {
			return err
		}
		log.Println("Program removed from autorun")
	}
	return nil
}
func updateConfig(newConfig Config) error {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		return err
	}
	cfg.Section("").Key("retroarch_path").SetValue(newConfig.RetroarchPath)
	cfg.Section("").Key("save_path").SetValue(newConfig.SavePath)
	cfg.Section("").Key("save_to_one_file").SetValue(strconv.FormatBool(newConfig.SaveToOneFile))
	cfg.Section("").Key("autorun").SetValue(strconv.FormatBool(newConfig.Autorun))
	cfg.Section("").Key("output_to_files").SetValue(strconv.FormatBool(newConfig.OutputToFiles))
	cfg.Section("").Key("web_port").SetValue(strconv.Itoa(newConfig.WebPort))
	cfg.Section("").Key("system_icon").SetValue(strconv.Itoa(newConfig.SystemIcon))
	cfg.Section("").Key("theme").SetValue(newConfig.Theme)
	cfg.Section("").Key("language").SetValue(newConfig.Language)
	cfg.Section("").Key("thumbnails_path").SetValue(newConfig.ThumbnailsPath)
	cfg.Section("").Key("enable_thumbnails").SetValue(strconv.FormatBool(newConfig.EnableThumbnails))
	cfg.Section("").Key("thumbnail_size").SetValue(newConfig.ThumbnailSize)
	cfg.Section("").Key("alternate_thumbnails").SetValue(strconv.FormatBool(newConfig.AlternateThumbnails))
	cfg.Section("").Key("thumbnail_switch_interval").SetValue(strconv.Itoa(newConfig.ThumbnailSwitchInterval))
	cfg.Section("").Key("fade_duration").SetValue(strconv.FormatFloat(newConfig.FadeDuration, 'f', 2, 64))
	cfg.Section("").Key("fade_type").SetValue(newConfig.FadeType)
	return cfg.SaveTo("config.ini")
}
func loadTranslations(language string) (Translations, string, error) {
	langFile := filepath.Join(langPath, language+".json")
	data, err := os.ReadFile(langFile)
	if err != nil {
		return nil, "", fmt.Errorf("error loading translation %s: %v", langFile, err)
	}

	var rawData map[string]interface{}
	err = json.Unmarshal(data, &rawData)
	if err != nil {
		return nil, "", fmt.Errorf("error parsing translation %s: %v", langFile, err)
	}

	langName, ok := rawData["language_name"].(string)
	if !ok {
		langName = language
	}

	translations := make(Translations)
	for key, value := range rawData {
		if key != "language_name" {
			if strValue, ok := value.(string); ok {
				translations[key] = strValue
			}
		}
	}

	return translations, langName, nil
}
func loadTemplates(theme string) error {
	templates = make(map[string]*template.Template)
	themeDir := filepath.Join(themePath, theme)
	defaultDir := filepath.Join(themePath, "default")
	files := []string{
		"index.html",
		"game.html",
		"system.html",
		"all.html",
		"settings.html",
		"thumbnails.html",
		"settings-games.html",
	}

	for _, file := range files {
		tmplPath := filepath.Join(themeDir, file)
		var tmpl *template.Template

		if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
			tmplPath = filepath.Join(defaultDir, file)
			if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
				return fmt.Errorf("template %s not found in theme %s or default theme", file, theme)
			}
			tmpl, _ = template.ParseFiles(tmplPath)
			log.Printf("Loaded template %s from default theme (not found in %s)", file, theme)
		} else {
			tmpl, _ = template.ParseFiles(tmplPath)
			log.Printf("Loaded template %s from theme %s", file, theme)
		}

		templates[file] = tmpl
	}
	return nil
}
func getAvailableThemes() []string {
	var themes []string
	dir, err := os.Open(themePath)
	if err != nil {
		log.Printf("Error reading Theme folder: %v", err)
		return []string{"default"}
	}
	defer func() {
		if err := dir.Close(); err != nil {
			log.Printf("failed close dir: %v", err)
		}
	}()
	dirs, err := dir.Readdir(-1)
	if err != nil {
		log.Printf("Error reading Theme contents: %v", err)
		return []string{"default"}
	}

	for _, d := range dirs {
		if d.IsDir() {
			themes = append(themes, d.Name())
		}
	}
	if len(themes) == 0 {
		return []string{"default"}
	}
	return themes
}
func getAvailableLanguages() []Language {
	var languages []Language
	dir, err := os.Open(langPath)
	if err != nil {
		log.Printf("Error reading lang folder: %v", err)
		return []Language{{Code: "en", Name: "English"}}
	}

	defer func() {
		if err := dir.Close(); err != nil {
			log.Printf("failed close: %v", err)
		}
	}()

	files, err := dir.Readdir(-1)
	if err != nil {
		log.Printf("Error reading lang contents: %v", err)
		return []Language{{Code: "en", Name: "English"}}
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			code := strings.TrimSuffix(f.Name(), ".json")
			_, langName, err := loadTranslations(code)
			if err != nil {
				log.Printf("Error loading language %s: %v", code, err)
				langName = code
			}
			languages = append(languages, Language{Code: code, Name: langName})
		}
	}

	if len(languages) == 0 {
		return []Language{{Code: "en", Name: "English"}}
	}
	return languages
}
func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	if _, ok := templates[tmplName]; !ok {
		log.Printf("Template %s not found", tmplName)
		http.Error(w, fmt.Sprintf("Server error: template %s not found", tmplName), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates[tmplName].Execute(w, data)
	if err != nil {
		log.Printf("Error rendering %s: %v", tmplName, err)
		http.Error(w, "Server error: failed to render template", http.StatusInternalServerError)
		return
	}
}
func writeOutputFiles(savePath, gameName, consoleName string, saveToOneFile bool) {
	if saveToOneFile {
		output := consoleName + ": " + gameName
		if err := os.WriteFile(filepath.Join(savePath, "output.txt"), []byte(output), 0644); err == nil {
			log.Printf(translations["data_updated_output"], output)
		} else {
			log.Printf("Error writing to output.txt: %v", err)
		}
	} else {
		if err := os.WriteFile(filepath.Join(savePath, "game.txt"), []byte(gameName), 0644); err == nil {
			log.Printf(translations["game_updated"], gameName)
		} else {
			log.Printf("Error writing to game.txt: %v", err)
		}
		if err := os.WriteFile(filepath.Join(savePath, "console.txt"), []byte(consoleName), 0644); err == nil {
			log.Printf(translations["system_updated"], consoleName)
		} else {
			log.Printf("Error writing to console.txt: %v", err)
		}
	}
}
func clearOutputFiles(savePath string, saveToOneFile bool) {
	if saveToOneFile {
		if err := os.WriteFile(filepath.Join(savePath, "output.txt"), []byte(""), 0644); err != nil {
			log.Printf("Error clearing output.txt: %v", err)
		} else {
			log.Println(translations["data_cleared_output"])
		}
	} else {
		if err := os.WriteFile(filepath.Join(savePath, "game.txt"), []byte(""), 0644); err != nil {
			log.Printf("Error clearing game.txt: %v", err)
		} else {
			log.Println(translations["data_cleared_game"])
		}
		if err := os.WriteFile(filepath.Join(savePath, "console.txt"), []byte(""), 0644); err != nil {
			log.Printf("Error clearing console.txt: %v", err)
		} else {
			log.Println(translations["data_cleared_console"])
		}
	}
}
func getThumbnailPaths(config Config, currentConsole, currentGame, theme string) ([]string, string, string) {

	var thumbnailPaths []string
	var thumbnailWidth, thumbnailHeight string
	if config.EnableThumbnails && config.ThumbnailsPath != "" && currentConsole != "" && currentGame != "" {
		currentGame = strings.TrimSpace(currentGame)
		currentConsole = strings.TrimSpace(currentConsole)

		titlesPath := filepath.Join(config.ThumbnailsPath, currentConsole, "Named_Titles", currentGame+".png")
		for _, filename := range []string{titlesPath, strings.ReplaceAll(titlesPath, "&", "_")} {
			if _, err := os.Stat(filename); !os.IsNotExist(err) {
				tmp := strings.ReplaceAll(currentGame, "&", "_")
				thumbnailPaths = append(thumbnailPaths, fmt.Sprintf("/thumbnails/%s/Named_Titles/%s.png", currentConsole, tmp))
				break
			}
		}

		boxartsPath := filepath.Join(config.ThumbnailsPath, currentConsole, "Named_Boxarts", currentGame+".png")
		for _, filename := range []string{boxartsPath, strings.ReplaceAll(boxartsPath, "&", "_")} {
			if _, err := os.Stat(filename); !os.IsNotExist(err) {
				tmp := strings.ReplaceAll(currentGame, "&", "_")
				thumbnailPaths = append(thumbnailPaths, fmt.Sprintf("/thumbnails/%s/Named_Boxarts/%s.png", currentConsole, tmp))
				break
			}
		}

		if len(thumbnailPaths) == 0 {
			noImagePath := filepath.Join(themePath, theme, "noimage.png")
			if _, err := os.Stat(noImagePath); !os.IsNotExist(err) {
				thumbnailPaths = append(thumbnailPaths, fmt.Sprintf("/theme/%s/noimage.png", theme))
			} else {
				defaultNoImagePath := filepath.Join(themePath, "default", "noimage.png")
				if _, err := os.Stat(defaultNoImagePath); !os.IsNotExist(err) {
					thumbnailPaths = append(thumbnailPaths, "/theme/default/noimage.png")
				} else {
					log.Printf("Warning: noimage.png not found in theme %s or default", theme)
				}
			}
		}

		if config.ThumbnailSize != "" && config.ThumbnailSize != "0" {
			parts := strings.Split(config.ThumbnailSize, "x")
			if len(parts) == 2 {
				if parts[0] != "" && parts[1] != "" {
					thumbnailWidth = parts[0] + "px"
					thumbnailHeight = parts[1] + "px"
				} else if parts[0] != "" {
					thumbnailWidth = parts[0] + "px"
				} else if parts[1] != "" {
					thumbnailHeight = parts[1] + "px"
				}
			}
		}
	}
	log.Printf("Returning thumbnail paths: %v", thumbnailPaths)
	return thumbnailPaths, thumbnailWidth, thumbnailHeight
}
func startWebServer(port int) {
	addr := fmt.Sprintf(":%d", port)

	http.Handle("/systems/", http.StripPrefix("/systems/", http.FileServer(http.Dir(systemsPath))))
	http.Handle("/theme/", http.StripPrefix("/theme/", http.FileServer(http.Dir(themePath))))
	http.Handle("/thumbnails/", http.StripPrefix("/thumbnails/", http.FileServer(http.Dir(config.ThumbnailsPath))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		configMutex.RLock()
		defer configMutex.RUnlock()
		translations, _, err := loadTranslations(config.Language)
		if err != nil {
			log.Printf("Error loading translations: %v", err)
			http.Error(w, "Server error: failed to load translations", http.StatusInternalServerError)
			return
		}
		var isRunning, _ = isRetroarchRunning()
		data := struct {
			Running                 bool
			CurrentGame             string
			CurrentConsole          string
			Theme                   string
			T                       Translations
			EnableThumbnails        bool
			ThumbnailPaths          []string
			ThumbnailWidth          string
			ThumbnailHeight         string
			AlternateThumbnails     bool
			ThumbnailSwitchInterval int
			Version                 string
			Port                    int
		}{
			Running:                 isRunning,
			CurrentGame:             currentGame,
			CurrentConsole:          currentConsole,
			Theme:                   config.Theme,
			T:                       translations,
			EnableThumbnails:        config.EnableThumbnails,
			ThumbnailPaths:          []string{},
			ThumbnailWidth:          "",
			ThumbnailHeight:         "",
			AlternateThumbnails:     config.AlternateThumbnails,
			ThumbnailSwitchInterval: config.ThumbnailSwitchInterval,
			Version:                 appVersion,
			Port:                    config.WebPort,
		}
		thumbnailPaths, thumbnailWidth, thumbnailHeight := getThumbnailPaths(config, currentConsole, currentGame, config.Theme)
		data.ThumbnailPaths = thumbnailPaths
		data.ThumbnailWidth = thumbnailWidth
		data.ThumbnailHeight = thumbnailHeight
		renderTemplate(w, "index.html", data)
		footer := `
        <script>
            var container = document.querySelector(".container.index-container");
            if (container) {
                container.innerHTML += '<div id="donate-section" style="text-align: center; padding: 10px; font-size: 14px;">' +
                    '<p class="donate-text">Please consider supporting the project with a donation: <a class="donate-link" href="https://www.donationalerts.com/r/ork8bit">DonationAlerts</a>' +
                    '<img src="/theme/heart.png" alt="Heart" class="donate-heart"></p>' +
                    '<p class="donate-text">©ork8bit aka Dmitriy Anatolyevich | <a class="donate-link" href="mailto:dmitrius.true@proton.me">dmitrius.true@proton.me</a></p>' +
                    '</div>';
            }
            function protectFooter() {
                var container = document.querySelector(".container.index-container");
                var footer = document.getElementById("donate-section");
                if (!footer || footer.style.display === "none" || !container.querySelector("#donate-section")) {
                    container.innerHTML += '<div id="donate-section" style="text-align: center; padding: 10px; font-size: 14px;">' +
                        '<p class="donate-text">Please consider supporting the project with a donation: <a class="donate-link" href="https://www.donationalerts.com/r/ork8bit">DonationAlerts</a>' +
                        '<img src="/theme/heart.png" alt="Heart" class="donate-heart"></p>' +
                        '<p class="donate-text">©ork8bit aka Dmitriy Anatolyevich | <a class="donate-link" href="mailto:dmitrius.true@proton.me">dmitrius.true@proton.me</a></p>' +
                        '</div>';
                }
            }
            setInterval(protectFooter, 1000);
        </script>
    `
		if _, err := w.Write([]byte(footer)); err != nil {
			log.Printf("failed to write footer: %v", err)
		}
	})
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		configMutex.RLock()
		defer configMutex.RUnlock()
		data := struct {
			CurrentGame string
			Theme       string
			Port        int
		}{
			CurrentGame: currentGame,
			Theme:       config.Theme,
			Port:        config.WebPort,
		}
		renderTemplate(w, "game.html", data)
	})
	http.HandleFunc("/system", func(w http.ResponseWriter, r *http.Request) {
		configMutex.RLock()
		defer configMutex.RUnlock()
		data := struct {
			SystemIcon     int
			CurrentConsole string
			IconFile       string
			Theme          string
			Port           int
		}{
			SystemIcon:     config.SystemIcon,
			CurrentConsole: currentConsole,
			Theme:          config.Theme,
			Port:           config.WebPort,
		}
		if config.SystemIcon > 0 && currentConsole != "" {
			if iconFile, exists := config.Systems[currentConsole]; exists {
				data.IconFile = iconFile
			}
		}
		renderTemplate(w, "system.html", data)
	})
	http.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		configMutex.RLock()
		defer configMutex.RUnlock()
		data := struct {
			SystemIcon     int
			CurrentConsole string
			CurrentGame    string
			IconFile       string
			Theme          string
			Port           int
		}{
			SystemIcon:     config.SystemIcon,
			CurrentConsole: currentConsole,
			CurrentGame:    currentGame,
			Theme:          config.Theme,
			Port:           config.WebPort,
		}
		if config.SystemIcon > 0 && currentConsole != "" {
			if iconFile, exists := config.Systems[currentConsole]; exists {
				data.IconFile = iconFile
			}
		}
		renderTemplate(w, "all.html", data)
	})
	http.HandleFunc("/thumbnails", func(w http.ResponseWriter, r *http.Request) {
		configMutex.RLock()
		defer configMutex.RUnlock()
		Width, Height := parseSizeToInt(config.ThumbnailSize)
		data := struct {
			CurrentGame             string
			Theme                   string
			ThumbnailPaths          []string
			ThumbnailWidth          string
			ThumbnailHeight         string
			AlternateThumbnails     bool
			ThumbnailSwitchInterval int
			FadeDuration            float64
			FadeType                string
			Port                    int
			Width                   int
			Height                  int
		}{
			CurrentGame:             currentGame,
			Theme:                   config.Theme,
			ThumbnailPaths:          []string{},
			ThumbnailWidth:          "",
			ThumbnailHeight:         "",
			AlternateThumbnails:     config.AlternateThumbnails,
			ThumbnailSwitchInterval: config.ThumbnailSwitchInterval,
			FadeDuration:            config.FadeDuration,
			FadeType:                config.FadeType,
			Port:                    config.WebPort,
			Width:                   Width,
			Height:                  Height,
		}

		data.ThumbnailPaths, data.ThumbnailWidth, data.ThumbnailHeight = getThumbnailPaths(config, currentConsole, currentGame, config.Theme)
		log.Printf("Serving /thumbnails with paths: %v", data.ThumbnailPaths)
		renderTemplate(w, "thumbnails.html", data)
	})
	http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		configMutex.RLock()
		currentConfig := config
		configMutex.RUnlock()
		if r.Method == "GET" {
			translations, _, err := loadTranslations(currentConfig.Language)
			if err != nil {
				log.Printf("Error loading translations: %v", err)
				http.Error(w, "Server error: failed to load translations", http.StatusInternalServerError)
				return
			}
			data := struct {
				Config    Config
				Themes    []string
				Languages []Language
				T         Translations
			}{
				Config:    currentConfig,
				Themes:    getAvailableThemes(),
				Languages: getAvailableLanguages(),
				T:         translations,
			}
			renderTemplate(w, "settings.html", data)
		} else if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}
			configMutex.Lock()
			defer configMutex.Unlock()
			config.RetroarchPath = r.FormValue("retroarch_path")
			config.SavePath = r.FormValue("save_path")
			config.SaveToOneFile = r.FormValue("save_to_one_file") == "on"
			config.Autorun = r.FormValue("autorun") == "on"
			config.OutputToFiles = r.FormValue("output_to_files") == "on"
			if port, err := strconv.Atoi(r.FormValue("web_port")); err == nil && port > 0 && port <= 65535 {
				config.WebPort = port
			}
			if icon, err := strconv.Atoi(r.FormValue("system_icon")); err == nil && icon >= 0 && icon <= 2 {
				config.SystemIcon = icon
			}
			newTheme := r.FormValue("theme")
			if _, err := os.Stat(filepath.Join(themePath, newTheme)); !os.IsNotExist(err) {
				config.Theme = newTheme
				if err := loadTemplates(config.Theme); err != nil {
					log.Printf("Error loading theme %s: %v", config.Theme, err)
				}
			}
			newLanguage := r.FormValue("language")
			if _, err := os.Stat(filepath.Join(langPath, newLanguage+".json")); !os.IsNotExist(err) {
				config.Language = newLanguage
				translations, _, err = loadTranslations(config.Language)
				if err != nil {
					log.Printf("Error reloading translations: %v", err)
				}
			}
			config.ThumbnailsPath = r.FormValue("thumbnails_path")
			config.EnableThumbnails = r.FormValue("enable_thumbnails") == "on"
			config.ThumbnailSize = r.FormValue("thumbnail_size")
			config.AlternateThumbnails = r.FormValue("alternate_thumbnails") == "on"
			if switchInterval, err := strconv.Atoi(r.FormValue("thumbnail_switch_interval")); err == nil && switchInterval >= 1 {
				config.ThumbnailSwitchInterval = switchInterval
			}
			fadeDurationStr := strings.Replace(r.FormValue("fade_duration"), ",", ".", -1)
			if duration, err := strconv.ParseFloat(fadeDurationStr, 64); err == nil && duration >= 0.1 {
				config.FadeDuration = duration
			}
			fadeType := r.FormValue("fade_type")
			validFadeTypes := []string{"ease", "ease-in", "ease-out", "ease-in-out", "linear"}
			for _, ft := range validFadeTypes {
				if fadeType == ft {
					config.FadeType = fadeType
					break
				}
			}
			if err := updateConfig(config); err != nil {
				http.Error(w, "Error saving settings", http.StatusInternalServerError)
				log.Printf("Error saving config.ini: %v", err)
				return
			}
			if config.Autorun != (r.FormValue("autorun") != "on") {
				if err := setAutorun(config.Autorun, "TrackGameName"); err != nil {
					log.Printf("Error updating autorun: %v", err)
				}
			}
			log.Println(translations["settings_updated"])
			http.Redirect(w, r, "/settings", http.StatusSeeOther)
		}
	})
	http.HandleFunc("/settings-games", func(w http.ResponseWriter, r *http.Request) {
		configMutex.RLock()
		currentConfig := config
		configMutex.RUnlock()

		log.Printf("Handling /settings-games, Method: %s", r.Method)

		if r.Method == "GET" {
			translations, _, err := loadTranslations(config.Language)
			if err != nil {
				log.Printf("Error loading translations: %v", err)
				http.Error(w, "Server error: failed to load translations", http.StatusInternalServerError)
				return
			}
			data := struct {
				Config        Config
				GameTemplates []GameTemplate
				T             Translations
				Port          int
			}{
				Config:        currentConfig,
				GameTemplates: gameTemplates,
				T:             translations,
				Port:          config.WebPort,
			}
			log.Println("Rendering settings-games.html")
			renderTemplate(w, "settings-games.html", data)
		} else if r.Method == "POST" {
			if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
				log.Printf("Error parsing form: %v", err)
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}

			processName := r.FormValue("process_name_display")
			windowTitle := r.FormValue("window_title")

			system := "Windows"
			game := strings.TrimSuffix(processName, ".exe")
			if processName == "retroarch.exe" {
				system = currentConsole
				game = currentGame
			}

			namedTitlesPath := ""
			namedBoxartsPath := ""
			if file, _, err := r.FormFile("named_titles"); err == nil {
				defer func() {
					if err := file.Close(); err != nil {
						log.Printf("failed close: %v", err)
					}
				}()
				titlesDir := filepath.Join(currentConfig.ThumbnailsPath, system, "Named_Titles")
				if err := os.MkdirAll(titlesDir, 0755); err != nil {
					log.Printf("Error creating Named_Titles dir: %v", err)
				} else {
					destPath := filepath.Join(titlesDir, game+".png")
					dest, err := os.Create(destPath)
					if err == nil {
						defer func() {
							if err := dest.Close(); err != nil {
								log.Printf("failed to close dest: %v", err)
							}
						}()
						if _, err := io.Copy(dest, file); err != nil {
							log.Printf("failed to copy data: %v", err)
						}
						namedTitlesPath = filepath.Join(system, "Named_Titles", game+".png")
						log.Printf("Saved named_titles to %s", destPath)
					}
				}
			}
			if file, _, err := r.FormFile("named_boxarts"); err == nil {
				defer func() {
					if err := file.Close(); err != nil {
						log.Printf("failed to close file: %v", err)
					}
				}()
				boxartsDir := filepath.Join(currentConfig.ThumbnailsPath, system, "Named_Boxarts")
				if err := os.MkdirAll(boxartsDir, 0755); err != nil {
					log.Printf("Error creating Named_Boxarts dir: %v", err)
				} else {
					destPath := filepath.Join(boxartsDir, game+".png")
					dest, err := os.Create(destPath)
					if err == nil {
						defer func() {
							if err := dest.Close(); err != nil {
								log.Printf("failed to close dest: %v", err)
							}
						}()
						if _, err := io.Copy(dest, file); err != nil {
							log.Printf("failed to copy data: %v", err)
						}
						namedBoxartsPath = filepath.Join(system, "Named_Boxarts", game+".png")
						log.Printf("Saved named_boxarts to %s", destPath)
					}
				}
			}

			gameTemplates = append(gameTemplates, GameTemplate{
				ProcessName:  processName,
				WindowTitle:  windowTitle,
				System:       system,
				Game:         game,
				NamedTitles:  namedTitlesPath,
				NamedBoxarts: namedBoxartsPath,
			})
			if err := saveGameTemplates(currentConfig.SavePath); err != nil {
				log.Printf("Error saving game templates: %v", err)
			}

			log.Println("Redirecting to /settings-games after POST")
			http.Redirect(w, r, "/settings-games", http.StatusSeeOther)
		}
	})
	http.HandleFunc("/settings-games/templates", func(w http.ResponseWriter, r *http.Request) {
		configMutex.RLock()
		defer configMutex.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(gameTemplates)
		if err != nil {
			log.Printf("Error encoding game templates: %v", err)
			http.Error(w, "Server error: failed to encode templates", http.StatusInternalServerError)
		}
	})
	http.HandleFunc("/startport", handleWebSocket)

	log.Printf("Web server started at http://localhost:%d", port)
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Printf("Web server error: %v", err)
		}
	}()
}

type processInfo struct {
	Name string `json:"name"`
	Pid  int32  `json:"pid"`
}

func getProcesses() ([]processInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	excludedUsers := map[string]struct{}{
		"СИСТЕМА":         {},
		"SYSTEM":          {},
		"LOCAL SERVICE":   {},
		"NETWORK SERVICE": {},
		"DWM-1":           {},
		"UMFD-1":          {},
		"UMFD-0":          {},
	}
	var userProcesses []processInfo
	for _, p := range processes {
		username, err := p.Username()
		if err != nil {
			continue
		}
		// Исключаем процессы, принадлежащие пользователям SYSTEM, LOCAL SERVICE и NETWORK SERVICE
		tmpName, tmpPid := "", int32(0)
		if _, excluded := excludedUsers[username]; !excluded {
			tmpName, _ = p.Name()
			tmpPid = p.Pid
			userProcesses = append(userProcesses, struct {
				Name string `json:"name"`
				Pid  int32  `json:"pid"`
			}{
				Name: tmpName,
				Pid:  tmpPid,
			})
		}
	}

	return userProcesses, nil
}
func getProcessInfo(pid interface{}) (interface{}, error) {
	var pPid int32
	switch p := pid.(type) {
	case int32:
		pPid = p
	case string:
		pidInt, _ := strconv.Atoi(p)
		pPid = int32(pidInt)
	default:
		return nil, fmt.Errorf("unsupported type: %T", pid)
	}
	p, _ := process.NewProcess(pPid)
	name, err := p.Name()
	if err != nil {
		name = ""
		log.Printf("Failed to get process name for PID %d: %v", pPid, err)
	}

	title, _ := getWindowTitle(pPid)
	if err != nil {
		log.Printf("Failed to get window title for PID %d: %v", pPid, err)
		title = ""
	}
	data := struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	}{Name: name, Title: title}
	return data, nil
}

func saveProcessInfo(processName string, windowTitle string) (bool, error) {

	system := "Windows"
	game := strings.TrimSuffix(processName, ".exe")
	if processName == "retroarch.exe" {
		return false, fmt.Errorf("retroarch.exe is not a valid process name")
	}

	namedTitlesPath := tmpNamedTitles
	namedBoxartsPath := tmpNamedBoxarts

	boxartsDir := filepath.Join(config.ThumbnailsPath, system, "Named_Boxarts")
	if err := os.MkdirAll(boxartsDir, 0755); err != nil {
		log.Printf("Error creating Named_Boxarts dir: %v", err)
	}
	titlesDir := filepath.Join(config.ThumbnailsPath, system, "Named_Titles")
	if err := os.MkdirAll(titlesDir, 0755); err != nil {
		log.Printf("Error creating Named_Boxarts dir: %v", err)
	}
	isFileExist, _ := isFile(namedTitlesPath)
	newNamedTitles := ""
	if isFileExist {
		newNamedTitles, _ = copyFileToDir(namedTitlesPath, titlesDir)
		newNamedTitles = system + "\\Named_Titles\\" + filepath.Base(newNamedTitles)
	}

	isFileExist, _ = isFile(namedBoxartsPath)
	newBoxArts := ""
	if isFileExist {
		newBoxArts, _ = copyFileToDir(namedBoxartsPath, boxartsDir)
		newBoxArts = system + "\\Named_Boxarts\\" + filepath.Base(newBoxArts)
	}
	gameTemplates = append(gameTemplates, GameTemplate{
		ProcessName:  processName,
		WindowTitle:  windowTitle,
		System:       system,
		Game:         game,
		NamedTitles:  newNamedTitles,
		NamedBoxarts: newBoxArts,
	})

	if err := saveGameTemplates(config.SavePath); err != nil {
		log.Printf("Error saving game templates: %v", err)
		return false, err
	}
	defer os.Remove(namedTitlesPath)
	defer os.Remove(namedBoxartsPath)
	tmpNamedTitles = ""
	tmpNamedBoxarts = ""
	return true, nil
}
func isFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		// Если ошибка, проверяем, существует ли путь
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	// Проверяем, является ли путь файлом
	return !info.IsDir(), nil
}
func removeGameTemplate(processName string) {
	filteredTemplates := []GameTemplate{}
	for _, templates := range gameTemplates {
		if templates.ProcessName != processName {
			filteredTemplates = append(filteredTemplates, templates)
		}
	}
	gameTemplates = filteredTemplates
}
func copyFileToDir(src string, dstDir string) (string, error) {
	// Проверяем, существует ли исходный файл
	sourceFile, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer sourceFile.Close()

	// Получаем имя файла из исходного пути
	fileName := filepath.Base(src)

	// Создаем полный путь для целевого файла
	dst := filepath.Join(dstDir, fileName)

	// Создаем целевой файл
	destinationFile, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer destinationFile.Close()

	// Копируем содержимое
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return "", err
	}

	// Копируем права доступа
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return "", err
	}
	err = os.Chmod(dst, sourceInfo.Mode())
	if err != nil {
		return "", err
	}

	// Возвращаем полный путь к новому файлу
	return dst, nil
}
func parseSizeToInt(sizeStr string) (int, int) {
	sizeStr = strings.TrimSpace(sizeStr)
	if sizeStr == "0" || sizeStr == "" {
		return 0, 0
	}

	parts := strings.Split(sizeStr, "x")
	if len(parts) != 2 {
		return 0, 0
	}

	var width, height int

	if parts[0] != "" {
		if w, err := strconv.Atoi(parts[0]); err == nil {
			width = w
		}
	}
	if parts[1] != "" {
		if h, err := strconv.Atoi(parts[1]); err == nil {
			height = h
		}
	}

	return width, height
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10 * 1024 * 1024,
	WriteBufferSize: 10 * 1024 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var tmpNamedTitles string
var tmpNamedBoxarts string

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка апгрейда WebSocket:", err)
		return
	}
	defer func() {
		clientsMutex.Lock()
		delete(clients, conn)
		clientsMutex.Unlock()
		if err := conn.Close(); err != nil {
			log.Printf("failed close conn: %v", err)
		}

	}()

	// Инициализируем клиента сразу при подключении
	clientsMutex.Lock()
	clients[conn] = ClientInfo{
		Conn:   conn,
		Screen: "", // Screen будет обновлен позже при регистрации
	}
	clientsMutex.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket отключён:", err)
			break
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(msg, &payload); err != nil {
			log.Println("Некорректный JSON:", err)
			continue
		}
		switch payload["type"] {
		case "register": // Регистрация клиента
			screen, _ := payload["screen"].(string)
			clientsMutex.Lock()
			info := clients[conn]
			info.Screen = screen
			clients[conn] = info
			clientsMutex.Unlock()
		case "get_data": // Запрос данных
			var data SendData
			switch payload["dataType"] {
			case "gameTemplates":
				data = SendData{
					Type:    "gameTemplates",
					Screen:  "settings-games",
					Payload: gameTemplates, // Отправляем список шаблонов игр
				}
			case "processes":
				processes, _ := getProcesses()
				data = SendData{
					Type:    "processes",
					Screen:  "settings-games",
					Payload: processes, // Отправляем список шаблонов игр
				}
			case "infoProcess":
				info, _ := getProcessInfo(payload["pid"])
				data = SendData{
					Type:    "infoProcess",
					Screen:  "settings-games",
					Payload: info, // Отправляем данные о процессе
				}
			}

			response, err := json.Marshal(data)
			if err != nil {
				log.Println("Ошибка сериализации JSON:", err)
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, response); err != nil {
				log.Println("Ошибка при отправке данных:", err)
			}
		case "saveData":

			// Обработка сохранения данных
			switch payload["dataType"] {
			case "saveFile":
				fileData, _ := payload["fileData"].(string)

				filename, _ := payload["name"].(string)

				var err error
				var file string
				// Обработка сохранения данных в зависимости от типа
				switch payload["imgType"] {
				case "named_titles":
					// Обработка сохранения named_titles
					file, err = saveToFile(fileData, filename, "/named_titles")
					tmpNamedTitles = file
				case "named_boxarts":
					// Обработка сохранения named_boxarts
					file, err = saveToFile(fileData, filename, "/named_boxarts/")
					tmpNamedBoxarts = file
				}
				fmt.Println("file:", file)
				if err != nil {
					log.Printf("Ошибка: %s", err)
				}

			case "saveProcess":
				dataForm, _ := payload["dataForm"].(map[string]interface{})
				processName, _ := dataForm["process_name_display"].(string)
				windowTitle, _ := dataForm["window_title"].(string)
				_, err := saveProcessInfo(processName, windowTitle)
				if err != nil {
					log.Printf("Error saving process info: %v", err)
				}
				data := SendData{
					Type:    "refresh",
					Screen:  "settings-games",
					Payload: true, // Отправляем данные о процессе
				}
				response, err := json.Marshal(data)
				if err := conn.WriteMessage(websocket.TextMessage, response); err != nil {
					log.Println("Ошибка при отправке данных:", err)
				}
			}
		case "delete":
			// Обработка удаления данных
			switch payload["dataType"] {
			case "deleteGameTemplate":
				processName, _ := payload["processName"].(string)
				removeGameTemplate(processName)
				data := SendData{
					Type:    "refresh",
					Screen:  "settings-games",
					Payload: true, // Отправляем данные о процессе
				}
				response, _ := json.Marshal(data)
				if err := conn.WriteMessage(websocket.TextMessage, response); err != nil {
					log.Println("Ошибка при отправке данных:", err)
				}
			}
		}
	}
}
func saveToFile(base64Data string, nameFile string, subdir string) (string, error) {
	// Удаляем префикс "data:image/png;base64,"
	parts := strings.SplitN(base64Data, ",", 2)
	if len(parts) != 2 {
		log.Printf("invalid base64 data")
		return "", fmt.Errorf("invalid base64 data")
	}
	// Декодируем base64 строку
	decodedData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		log.Printf("invalid base64 data")
		return "", fmt.Errorf("invalid base64 data")
	}

	// Создаём директорию
	tempDir := os.TempDir()
	subDir := filepath.Join(tempDir, subdir)
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		log.Printf("Ошибка создания папки: %v\n", err)
		return "", fmt.Errorf("Ошибка создания папки: %v\n", err)
	}
	// Создаём файл с фиксированным именем
	filePath := filepath.Join(subDir, nameFile+".png")
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("Ошибка создания файла: %v", err)
		return "", fmt.Errorf("Ошибка создания файла: %v", err)
	}
	// Записываем данные в файл
	_, err2 := file.Write(decodedData)
	if err2 != nil {
		log.Printf("Ошибка записи в файл: %v", err2)
		return "", fmt.Errorf("Ошибка записи в файл: %v", err2)
	}
	defer file.Close()
	return filePath, nil
}

type SendData struct {
	Type    string      `json:"type"`
	Screen  string      `json:"screen"`
	Payload interface{} `json:"payload"`
}

func broadcastTo(screen string, msg string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for conn, client := range clients {
		if client.Screen == screen || screen == "*" {
			err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("Ошибка при отправке:", err)
				delete(clients, conn)
			}
		}
	}
}

type OutgoingMessage struct {
	Type    string      `json:"type"`
	Screen  string      `json:"screen"`
	Payload interface{} `json:"payload"`
}

func sendUpdate(screen string, payload interface{}) {
	msg := OutgoingMessage{
		Type:    "update",
		Screen:  screen,
		Payload: payload,
	}

	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("Ошибка сериализации JSON:", err)
		return
	}

	broadcastTo(screen, string(jsonBytes))
}
func main() {
	logFile, err := os.OpenFile("trackgamename.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			log.Printf("failed to close logFile: %v", err)
		}
	}()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)

	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Error loading config.ini: %v", err)
		cfg = ini.Empty()
		cfg.Section("").Key("retroarch_path").SetValue("C:\\RetroArch-Win64")
		cfg.Section("").Key("save_path").SetValue("")
		cfg.Section("").Key("save_to_one_file").SetValue("false")
		cfg.Section("").Key("autorun").SetValue("false")
		cfg.Section("").Key("output_to_files").SetValue("true")
		cfg.Section("").Key("web_port").SetValue("3489")
		cfg.Section("").Key("system_icon").SetValue("0")
		cfg.Section("").Key("theme").SetValue("default")
		cfg.Section("").Key("language").SetValue("en")
		cfg.Section("").Key("thumbnails_path").SetValue("")
		cfg.Section("").Key("enable_thumbnails").SetValue("false")
		cfg.Section("").Key("thumbnail_size").SetValue("0")
		cfg.Section("").Key("alternate_thumbnails").SetValue("false")
		cfg.Section("").Key("thumbnail_switch_interval").SetValue("5")
		cfg.Section("").Key("fade_duration").SetValue("0.5")
		cfg.Section("").Key("fade_type").SetValue("ease-out")
		cfg.Section("systems").Key("Nintendo - Nintendo Entertainment System").SetValue("nes.png")
		err = cfg.SaveTo("config.ini")
		if err != nil {
			log.Printf("Error creating config.ini: %v", err)
			os.Exit(1)
		}
		log.Println("Config file config.ini created. Continuing execution.")
	}

	config = Config{
		Systems: make(map[string]string),
	}
	err = cfg.MapTo(&config)
	if err != nil {
		log.Printf("Error reading config.ini: %v", err)
		os.Exit(1)
	}

	systemsSection := cfg.Section("systems")
	for _, key := range systemsSection.Keys() {
		config.Systems[key.Name()] = key.String()
	}

	savePath := config.SavePath
	if savePath == "" {
		savePath, err = os.Getwd()
		if err != nil {
			log.Printf("Error getting current directory: %v", err)
			os.Exit(1)
		}
	}

	if err := os.MkdirAll(savePath, 0755); err != nil {
		log.Printf("Error creating save_path folder: %v", err)
		os.Exit(1)
	}

	systemsPath = filepath.Join(savePath, "systems")
	if err := os.MkdirAll(systemsPath, 0755); err != nil {
		log.Printf("Error creating systems folder: %v", err)
		os.Exit(1)
	}

	themePath = filepath.Join(savePath, "Theme")
	if err := os.MkdirAll(filepath.Join(themePath, "default"), 0755); err != nil {
		log.Printf("Error creating Theme/default folder: %v", err)
		os.Exit(1)
	}

	if _, err := os.Stat(filepath.Join(themePath, config.Theme)); os.IsNotExist(err) {
		log.Printf("Theme %s not found, falling back to default", config.Theme)
		config.Theme = "default"
	}
	if err := loadTemplates(config.Theme); err != nil {
		log.Printf("Error loading theme %s: %v, falling back to default", config.Theme, err)
		config.Theme = "default"
		if err := loadTemplates(config.Theme); err != nil {
			log.Printf("Error loading default theme: %v", err)
			os.Exit(1)
		}
	}

	langPath = filepath.Join("lang")
	if err := os.MkdirAll(langPath, 0755); err != nil {
		log.Printf("Error creating lang folder: %v", err)
		os.Exit(1)
	}

	configMutex.Lock()
	translations, _, err = loadTranslations(config.Language)
	if err != nil {
		log.Printf("Error loading translations: %v", err)
		config.Language = "en"
		translations, _, err = loadTranslations(config.Language)
		if err != nil {
			log.Printf("Error loading default translations: %v", err)
			os.Exit(1)
		}
	}
	configMutex.Unlock()
	log.Println(translations["app_started"])

	if err := loadGameTemplates(savePath); err != nil {
		log.Printf("Error loading game templates: %v", err)
	}

	if err := setAutorun(config.Autorun, "TrackGameName"); err != nil {
		log.Printf("Error setting autorun at startup: %v", err)
	}

	startWebServer(config.WebPort)

	lplPath := filepath.Join(config.RetroarchPath, "content_history.lpl")
	log.Printf("Path to content_history.lpl: %s", lplPath)

	systray.Run(onReady(savePath), onExit)
}
func getForegroundProcessPID() (int32, error) {
	pid := robotgo.GetPid()
	if pid == -1 {
		return 0, fmt.Errorf("could not get foreground process PID")
	}
	return int32(pid), nil
}
func onReady(savePath string) func() {
	return func() {
		systray.SetTitle(translations["title"])
		systray.SetTooltip(translations["title"])
		log.Println(translations["systray_initialized"])

		gameItem := systray.AddMenuItem(translations["game_not_detected"], translations["game_not_detected"])
		consoleItem := systray.AddMenuItem(translations["system_not_detected"], translations["system_not_detected"])
		systray.AddSeparator()
		openWebItem := systray.AddMenuItem(translations["open_web_page"], translations["open_web_page_tip"])
		openSettingsItem := systray.AddMenuItem(translations["open_settings"], translations["open_settings_tip"])
		quitItem := systray.AddMenuItem(translations["exit"], translations["exit_tip"])
		log.Println(translations["menu_items_added"])

		gamename := ""
		var lastState bool
		var initialized bool
		lastGame := ""
		lastConsole := ""
		updateInfo := func(console, game string) {
			configMutex.Lock()
			if config.OutputToFiles {
				currentSavePath := config.SavePath
				if currentSavePath == "" {
					currentSavePath = savePath
				}
				writeOutputFiles(currentSavePath, game, console, config.SaveToOneFile)
			}
			currentGame = game
			currentConsole = console
			configMutex.Unlock()

			gameItem.SetTitle("Game: " + game)
			consoleItem.SetTitle("System: " + console)
			gamename = game

			if lastGame != game || lastConsole != console {
				sendUpdate("game", map[string]string{
					"game": game,
				})

				icons := ""
				if config.SystemIcon > 0 && currentConsole != "" {
					if iconFile, exists := config.Systems[currentConsole]; exists {
						icons = iconFile
					}
				}

				thumbnailPaths, thumbnailWidth, thumbnailHeight := getThumbnailPaths(config, console, game, config.Theme)
				data := struct {
					Game   string   `json:"game"`
					Paths  []string `json:"paths"`
					Width  string   `json:"width"`
					Height string   `json:"height"`
				}{
					Game:   game,
					Paths:  thumbnailPaths,
					Width:  thumbnailWidth,
					Height: thumbnailHeight,
				}
				sendUpdate("system", map[string]string{
					"console": console,
					"icon":    icons,
				})

				sendUpdate("all", map[string]string{
					"console": console,
					"game":    game,
					"icon":    icons,
				})
				sendUpdate("thumbnails", data)
				log.Printf("Updated info: Game=%s, Console=%s", game, console)
				lastGame = game
				lastConsole = console
			}

		}

		go func() {
			for {
				currentState, _ := isRetroarchRunning() // функция проверки запущен ли RetroArch

				if !initialized || currentState != lastState { // если состояние изменилось
					if currentState && len(activeIcon) > 0 {
						systray.SetIcon(activeIcon)
						//log.Println(translations["retroarch_running_icon"])
					} else if len(inactiveIcon) > 0 {
						systray.SetIcon(inactiveIcon)
						//log.Println(translations["retroarch_closed_icon"])
						configMutex.RLock()
						if config.OutputToFiles {
							clearOutputFiles(savePath, config.SaveToOneFile)
						}
						configMutex.RUnlock()
						gameItem.SetTitle(translations["game_not_detected"])
						consoleItem.SetTitle(translations["system_not_detected"])
						gamename = ""
						currentGame = ""
						currentConsole = ""
					}
					lastState = currentState
					initialized = true // инициализация завершена
				}

				processes, err := process.Processes()
				if err != nil {
					log.Printf("Error getting processes: %v", err)
					time.Sleep(1 * time.Second)
					continue
				}
				countRunningProcesses := 0

				for i, gameProcc := range gameTemplates {
					isRunning := false
					var pid int32 = 0

					for _, proc := range processes {
						name, err := proc.Name()
						if err != nil {
							continue
						}
						if strings.EqualFold(name, gameProcc.ProcessName) {
							isRunning = true
							pid = proc.Pid
							countRunningProcesses++
							break
						}
					}
					gameTemplates[i].pid = pid
					gameTemplates[i].isRunning = isRunning
				}

				foregroundPID, err := getForegroundProcessPID()
				if err != nil {
					log.Printf("Error getting foreground process: %v", err)
					time.Sleep(1 * time.Second)
					continue
				}
			loopGame: // цикл по играм
				for _, gameProcc := range gameTemplates {
					if gameProcc.isRunning { //игра запущена
						systray.SetIcon(activeIcon)
						if countRunningProcesses > 1 { // несколько процессов
							if foregroundPID == gameProcc.pid { // фокус на игре
								if gameProcc.WindowTitle == "RetroArch" { // если игра RetroArch
									newGamename, consoleName, err1, err2 := getInfoGameRetroArch() // получаем название игры
									if err1 != nil || err2 != nil {
										log.Printf("Error reading content_history.lpl: %v, %v", err1, err2)
									} else {
										if gamename != newGamename && newGamename != "" {
											updateInfo(consoleName, newGamename) // обновляем информацию
										}
									}
									break loopGame
								} else {
									if gameProcc.isRunning { // если игра запущена
										windowTitle := gameProcc.WindowTitle
										updateInfo(gameProcc.System, windowTitle) // обновляем информацию
									}
									break loopGame
								}
							}
						} else if gameProcc.WindowTitle == "RetroArch" { //один процесс и если игра RetroArch

							newGamename, consoleName, err1, err2 := getInfoGameRetroArch() // получаем название игры
							if err1 != nil || err2 != nil {
								log.Printf("Error reading content_history.lpl: %v, %v", err1, err2)
							} else {
								if gamename != newGamename && newGamename != "" {
									updateInfo(consoleName, newGamename) // обновляем информацию
								}
							}
							break loopGame
						} else { //один процесс и не RetroArch
							for _, gameProcSecond := range gameTemplates {
								if gameProcSecond.isRunning {
									windowTitle := gameProcSecond.WindowTitle
									if windowTitle == "" {
										windowTitle, err = getWindowTitle(gameProcSecond.pid)
										if err != nil {
											windowTitle = gameProcSecond.Game
										}
									}
									updateInfo(gameProcSecond.System, windowTitle)
								}
							}
							break loopGame
						}
					}
				}
				time.Sleep(1 * time.Second)
			}
		}()

		go func() {
			for {
				select {
				case <-openWebItem.ClickedCh:
					configMutex.RLock()
					url := fmt.Sprintf("http://localhost:%d/", config.WebPort)
					configMutex.RUnlock()
					err := openBrowser(url)
					if err != nil {
						log.Printf("Error opening browser: %v", err)
					} else {
						log.Println(translations["web_page_opened"])
					}
				case <-openSettingsItem.ClickedCh:
					configMutex.RLock()
					url := fmt.Sprintf("http://localhost:%d/settings", config.WebPort)
					configMutex.RUnlock()
					err := openBrowser(url)
					if err != nil {
						log.Printf("Error opening settings page: %v", err)
					} else {
						log.Println(translations["settings_page_opened"])
					}
				case <-quitItem.ClickedCh:
					log.Println("Exit clicked received")
					systray.Quit()
					return
				}
			}
		}()
	}
}
func getInfoGameRetroArch() (string, string, error, error) {
	configMutex.RLock()
	currentLplPath := filepath.Join(config.RetroarchPath, "content_history.lpl")
	configMutex.RUnlock()

	newGameLine, err1 := findFirstLine(currentLplPath, `"label":`)
	newCoreLine, err2 := findFirstLine(currentLplPath, `"db_name":`)

	if err1 != nil || err2 != nil {
		log.Printf("Error reading content_history.lpl: %v, %v", err1, err2)
		return "", "", err1, err2
	}

	newGamename, _ := extractValue(newGameLine, `"label": "`)
	shortCoreName, _ := extractValue(newCoreLine, `"db_name": "`)
	consoleName := cut(shortCoreName, ".", 0)

	return newGamename, consoleName, nil, nil
}
func findFirstLine(filePath, search string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close key: %v", err)
		}
	}()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, search) {
			return line, nil
		}
	}
	return "", scanner.Err()
}
func cut(input, delimiter string, field int) string {
	parts := strings.Split(input, delimiter)
	if field < len(parts) {
		return parts[field]
	}
	return ""
}
func extractValue(line, pattern string) (string, string) {

	startIdx := strings.Index(line, pattern)

	if startIdx == -1 {
		return "", ""
	}
	start := startIdx + len(pattern)

	end := strings.Index(line[start:], `"`)
	if end == -1 {
		return "", ""
	}
	end += start
	if end <= start {
		return "", ""
	}
	fullName := line[start:end]
	shortName := strings.Split(fullName, "(")[0]
	shortName = strings.TrimSpace(shortName)

	return shortName, fullName
}
func loadGameTemplates(savePath string) error {
	gamesFile := filepath.Join(savePath, "games.json")
	if _, err := os.Stat(gamesFile); os.IsNotExist(err) {
		if err := os.WriteFile(gamesFile, []byte("[]"), 0644); err != nil {
			return fmt.Errorf("error creating games.json: %v", err)
		}

		log.Println("Created empty games.json")
		return nil
	}
	gameTemplates = []GameTemplate{}
	game := GameTemplate{}
	game.WindowTitle = "RetroArch"
	game.ProcessName = "retroarch.exe"
	game.isRunning, game.pid = isRetroarchRunning()
	data, err := os.ReadFile(gamesFile)
	if err != nil {
		return fmt.Errorf("error reading games.json: %v", err)
	}
	if err := json.Unmarshal(data, &gameTemplates); err != nil {
		return fmt.Errorf("error parsing games.json: %v", err)
	}

	gameTemplates = append(gameTemplates, game)
	gameTemplates = removeDuplicates(gameTemplates)
	log.Println("Loaded game templates from games.json")
	return nil
}
func saveGameTemplates(savePath string) error {
	gamesFile := filepath.Join(savePath, "games.json")
	gameTemplates = removeDuplicates(gameTemplates)
	data, err := json.MarshalIndent(gameTemplates, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshaling game templates: %v", err)
	}
	if err := os.WriteFile(gamesFile, data, 0644); err != nil {
		return fmt.Errorf("error writing games.json: %v", err)
	}
	log.Println("Saved game templates to games.json")
	return nil
}
func removeDuplicates(templates []GameTemplate) []GameTemplate {
	unique := make(map[string]GameTemplate)
	for _, tmplt := range templates {
		// Создаем уникальный ключ на основе полей структуры
		key := tmplt.ProcessName + "|" + tmplt.WindowTitle
		unique[key] = tmplt
	}

	// Преобразуем карту обратно в срез
	result := make([]GameTemplate, 0, len(unique))
	for _, tmplt := range unique {
		result = append(result, tmplt)
	}
	return result
}
func getWindowTitle(pid int32) (string, error) {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/FO", "CSV", "/V")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("tasklist failed: %v", err)
	}

	output := out.String()
	if strings.Contains(output, "No tasks are running") {
		return "", nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("unexpected tasklist output")
	}

	fields := strings.Split(lines[1], ",")
	if len(fields) < 9 {
		return "", fmt.Errorf("invalid tasklist output format")
	}

	title := strings.Trim(fields[8], `"`)
	return title, nil
}
func openBrowser(url string) error {
	return exec.Command("cmd", "/c", "start", url).Start()
}
func onExit() {
	log.Println(translations["app_exited"])
	os.Exit(0)
}
