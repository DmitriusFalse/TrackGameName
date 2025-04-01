package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
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

//go:embed active.ico
var activeIcon []byte

//go:embed inactive.ico
var inactiveIcon []byte

type Config struct {
	RetroarchPath    string            `ini:"retroarch_path"`
	SavePath         string            `ini:"save_path"`
	SaveToOneFile    bool              `ini:"save_to_one_file"`
	Autorun          bool              `ini:"autorun"`
	OutputToFiles    bool              `ini:"output_to_files"`
	WebPort          int               `ini:"web_port"`
	SystemIcon       int               `ini:"system_icon"`
	RefreshInterval  int               `ini:"refresh_interval"`
	Theme            string            `ini:"theme"`
	Language         string            `ini:"language"`
	ThumbnailsPath   string            `ini:"thumbnails_path"`
	EnableThumbnails bool              `ini:"enable_thumbnails"`
	ThumbnailSize    string            `ini:"thumbnail_size"`
	Systems          map[string]string `ini:"systems"`
}

type Translations map[string]string

var (
	currentGame         string
	currentGameFullName string
	currentConsole      string
	configMutex         sync.RWMutex
	systemsPath         string
	themePath           string
	langPath            string
	config              Config
	templates           map[string]*template.Template
	translations        Translations // Глобальная переменная для переводов
)

func findFirstLine(filePath, search string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

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

	// Полное название (всё между кавычек)
	fullName := line[start:end]

	// Короткое название (обрезаем до первого скобочного уточнения)
	shortName := strings.Split(fullName, "(")[0]
	shortName = strings.TrimSpace(shortName)

	return shortName, fullName
}

func isRetroarchRunning() bool {
	processes, err := process.Processes()
	if err != nil {
		log.Printf("Error getting process list: %v", err)
		return false
	}
	for _, p := range processes {
		name, err := p.Name()
		if err == nil && strings.ToLower(name) == "retroarch.exe" {
			return true
		}
	}
	return false
}

func setAutorun(enable bool, appName string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()

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
		if err != nil && err != registry.ErrNotExist {
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
	cfg.Section("").Key("refresh_interval").SetValue(strconv.Itoa(newConfig.RefreshInterval))
	cfg.Section("").Key("theme").SetValue(newConfig.Theme)
	cfg.Section("").Key("language").SetValue(newConfig.Language)
	cfg.Section("").Key("thumbnails_path").SetValue(newConfig.ThumbnailsPath)
	cfg.Section("").Key("enable_thumbnails").SetValue(strconv.FormatBool(newConfig.EnableThumbnails))
	cfg.Section("").Key("thumbnail_size").SetValue(newConfig.ThumbnailSize)
	return cfg.SaveTo("config.ini")
}

func loadTranslations(language string) (Translations, error) {
	langFile := filepath.Join(langPath, language+".json")
	data, err := os.ReadFile(langFile)
	if err != nil {
		return nil, fmt.Errorf("error loading translation %s: %v", langFile, err)
	}
	var translations Translations
	err = json.Unmarshal(data, &translations)
	if err != nil {
		return nil, fmt.Errorf("error parsing translation %s: %v", langFile, err)
	}
	return translations, nil
}

func loadTemplates(theme string) error {
	templates = make(map[string]*template.Template)
	themeDir := filepath.Join(themePath, theme)
	files := []string{"index.html", "game.html", "system.html", "all.html", "settings.html", "thumbnails.html"}
	for _, file := range files {
		tmplPath := filepath.Join(themeDir, file)
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			return fmt.Errorf("error loading template %s: %v", file, err)
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
	defer dir.Close()

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

func getAvailableLanguages() []string {
	var languages []string
	dir, err := os.Open(langPath)
	if err != nil {
		log.Printf("Error reading lang folder: %v", err)
		return []string{"en"}
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		log.Printf("Error reading lang contents: %v", err)
		return []string{"en"}
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			languages = append(languages, strings.TrimSuffix(f.Name(), ".json"))
		}
	}
	if len(languages) == 0 {
		return []string{"en"}
	}
	return languages
}

// Новая функция для рендеринга шаблонов
func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates[tmplName].Execute(w, data); err != nil {
		log.Printf("Error rendering %s: %v", tmplName, err)
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}

// Новая функция для работы с configMutex
func withConfigLock(fn func()) {
	configMutex.RLock()
	defer configMutex.RUnlock()
	fn()
}

// Новая функция для записи данных в файлы
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

// Новая функция для очистки файлов
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

// Функция для получения пути к миниатюре
func getThumbnailPath(config Config, currentConsole, currentGame, theme string) (string, string, string) {
	var thumbnailPath, thumbnailWidth, thumbnailHeight string
	if config.EnableThumbnails && config.ThumbnailsPath != "" && currentConsole != "" && currentGame != "" {
		currentGame = strings.TrimSpace(currentGame)
		currentConsole = strings.TrimSpace(currentConsole)
		fullPath := filepath.Join(config.ThumbnailsPath, currentConsole, "Named_Titles", currentGame+".png")
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			thumbnailPath = fmt.Sprintf("/thumbnails/%s/Named_Titles/%s.png", currentConsole, currentGame)
		} else {
			thumbnailPath = fmt.Sprintf("/theme/%s/noimage.png", theme)
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
	return thumbnailPath, thumbnailWidth, thumbnailHeight
}

func startWebServer(port int) {
	addr := fmt.Sprintf(":%d", port)

	http.Handle("/systems/", http.StripPrefix("/systems/", http.FileServer(http.Dir(systemsPath))))
	http.Handle("/theme/", http.StripPrefix("/theme/", http.FileServer(http.Dir(themePath))))

	http.HandleFunc("/thumbnails/", func(w http.ResponseWriter, r *http.Request) {
		withConfigLock(func() {
			if !config.EnableThumbnails || config.ThumbnailsPath == "" {
				http.Error(w, "Thumbnails are disabled or path not set", http.StatusNotFound)
				return
			}
			http.StripPrefix("/thumbnails/", http.FileServer(http.Dir(config.ThumbnailsPath))).ServeHTTP(w, r)
		})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		withConfigLock(func() {
			translations, err := loadTranslations(config.Language)
			if err != nil {
				log.Printf("Error loading translations: %v", err)
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}
			data := struct {
				RefreshInterval  int
				Running          bool
				CurrentGame      string
				CurrentConsole   string
				Theme            string
				T                Translations
				EnableThumbnails bool
				ThumbnailPath    string
				ThumbnailWidth   string
				ThumbnailHeight  string
			}{
				RefreshInterval:  config.RefreshInterval,
				Running:          isRetroarchRunning(),
				CurrentGame:      currentGame,
				CurrentConsole:   currentConsole,
				Theme:            config.Theme,
				T:                translations,
				EnableThumbnails: config.EnableThumbnails,
				ThumbnailPath:    "",
				ThumbnailWidth:   "",
				ThumbnailHeight:  "",
			}
			data.ThumbnailPath, data.ThumbnailWidth, data.ThumbnailHeight = getThumbnailPath(config, currentConsole, currentGameFullName, config.Theme)
			fmt.Println(data.ThumbnailPath)
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
			w.Write([]byte(footer))
		})
	})

	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		withConfigLock(func() {
			data := struct {
				RefreshInterval int
				CurrentGame     string
				Theme           string
			}{
				RefreshInterval: config.RefreshInterval,
				CurrentGame:     currentGame,
				Theme:           config.Theme,
			}
			renderTemplate(w, "game.html", data)
		})
	})

	http.HandleFunc("/system", func(w http.ResponseWriter, r *http.Request) {
		withConfigLock(func() {
			data := struct {
				RefreshInterval int
				SystemIcon      int
				CurrentConsole  string
				IconFile        string
				Theme           string
			}{
				RefreshInterval: config.RefreshInterval,
				SystemIcon:      config.SystemIcon,
				CurrentConsole:  currentConsole,
				Theme:           config.Theme,
			}
			if config.SystemIcon > 0 && currentConsole != "" {
				if iconFile, exists := config.Systems[currentConsole]; exists {
					data.IconFile = iconFile
				}
			}
			renderTemplate(w, "system.html", data)
		})
	})

	http.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		withConfigLock(func() {
			data := struct {
				RefreshInterval int
				SystemIcon      int
				CurrentConsole  string
				CurrentGame     string
				IconFile        string
				Theme           string
			}{
				RefreshInterval: config.RefreshInterval,
				SystemIcon:      config.SystemIcon,
				CurrentConsole:  currentConsole,
				CurrentGame:     currentGame,
				Theme:           config.Theme,
			}
			if config.SystemIcon > 0 && currentConsole != "" {
				if iconFile, exists := config.Systems[currentConsole]; exists {
					data.IconFile = iconFile
				}
			}
			renderTemplate(w, "all.html", data)
		})
	})

	http.HandleFunc("/thumbnails", func(w http.ResponseWriter, r *http.Request) {
		withConfigLock(func() {
			data := struct {
				RefreshInterval  int
				CurrentGame      string
				CurrentConsole   string
				Theme            string
				EnableThumbnails bool
				ThumbnailPath    string
				ThumbnailWidth   string
				ThumbnailHeight  string
			}{
				RefreshInterval:  config.RefreshInterval,
				CurrentGame:      currentGame,
				CurrentConsole:   currentConsole,
				Theme:            config.Theme,
				EnableThumbnails: config.EnableThumbnails,
				ThumbnailPath:    "",
				ThumbnailWidth:   "",
				ThumbnailHeight:  "",
			}
			fmt.Println(currentGameFullName)
			data.ThumbnailPath, data.ThumbnailWidth, data.ThumbnailHeight = getThumbnailPath(config, currentConsole, currentGameFullName, config.Theme)
			renderTemplate(w, "thumbnails.html", data)
		})
	})

	http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		configMutex.RLock()
		currentConfig := config
		configMutex.RUnlock()

		if r.Method == "GET" {
			translations, err := loadTranslations(currentConfig.Language)
			if err != nil {
				log.Printf("Error loading translations: %v", err)
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}
			data := struct {
				Config    Config
				Themes    []string
				Languages []string
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
			if interval, err := strconv.Atoi(r.FormValue("refresh_interval")); err == nil && interval >= 1 {
				config.RefreshInterval = interval
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
				translations, err = loadTranslations(config.Language)
				if err != nil {
					log.Printf("Error reloading translations: %v", err)
				}
			}
			config.ThumbnailsPath = r.FormValue("thumbnails_path")
			config.EnableThumbnails = r.FormValue("enable_thumbnails") == "on"
			config.ThumbnailSize = r.FormValue("thumbnail_size")

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

	log.Printf("Web server started at http://localhost:%d", port)
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Printf("Web server error: %v", err)
		}
	}()
}

func main() {
	logFile, err := os.OpenFile("trackgamename.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)

	// Загружаем конфигурацию
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
		cfg.Section("").Key("refresh_interval").SetValue("20")
		cfg.Section("").Key("theme").SetValue("default")
		cfg.Section("").Key("language").SetValue("en")
		cfg.Section("").Key("thumbnails_path").SetValue("")
		cfg.Section("").Key("enable_thumbnails").SetValue("false")
		cfg.Section("").Key("thumbnail_size").SetValue("0")
		cfg.Section("systems").Key("Nintendo").SetValue("nes.png")
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

	langPath = filepath.Join("lang")
	if err := os.MkdirAll(langPath, 0755); err != nil {
		log.Printf("Error creating lang folder: %v", err)
		os.Exit(1)
	}

	if err := loadTemplates(config.Theme); err != nil {
		log.Printf("Error loading theme %s: %v", config.Theme, err)
		config.Theme = "default"
		if err := loadTemplates(config.Theme); err != nil {
			log.Printf("Error loading default theme: %v", err)
			os.Exit(1)
		}
	}

	// Загружаем переводы
	translations, err = loadTranslations(config.Language)
	if err != nil {
		log.Printf("Error loading translations: %v", err)
		config.Language = "en" // По умолчанию английский
		translations, err = loadTranslations(config.Language)
		if err != nil {
			log.Printf("Error loading default translations: %v", err)
			os.Exit(1)
		}
	}
	log.Println(translations["app_started"])

	if err := setAutorun(config.Autorun, "TrackGameName"); err != nil {
		log.Printf("Error setting autorun at startup: %v", err)
	}

	startWebServer(config.WebPort)

	lplPath := filepath.Join(config.RetroarchPath, "content_history.lpl")
	log.Printf("Path to content_history.lpl: %s", lplPath)

	systray.Run(onReady(savePath), onExit)
}

func onReady(savePath string) func() {
	return func() {
		if isRetroarchRunning() && len(activeIcon) > 0 {
			systray.SetIcon(activeIcon)
			log.Println(translations["retroarch_running_icon"])
		} else if len(inactiveIcon) > 0 {
			systray.SetIcon(inactiveIcon)
			log.Println(translations["retroarch_not_running_icon"])
		} else {
			log.Println(translations["icons_not_loaded"])
		}

		systray.SetTitle(translations["title"])
		systray.SetTooltip(translations["title"])
		log.Println(translations["systray_initialized"])

		gameItem := systray.AddMenuItem(translations["game_not_detected"], translations["game_not_detected"])
		consoleItem := systray.AddMenuItem(translations["system_not_detected"], translations["system_not_detected"])
		systray.AddSeparator()
		openWebItem := systray.AddMenuItem(translations["open_web_page"], translations["open_web_page_tip"])
		quitItem := systray.AddMenuItem(translations["exit"], translations["exit_tip"])
		log.Println(translations["menu_items_added"])

		gamename := ""
		var lastState bool

		go func() {
			for {
				currentState := isRetroarchRunning()
				if currentState != lastState {
					if currentState && len(activeIcon) > 0 {
						systray.SetIcon(activeIcon)
						log.Println(translations["retroarch_running_icon"])
					} else if len(inactiveIcon) > 0 {
						systray.SetIcon(inactiveIcon)
						log.Println(translations["retroarch_closed_icon"])
						configMutex.RLock()
						if config.OutputToFiles {
							clearOutputFiles(savePath, config.SaveToOneFile)
						}
						configMutex.RUnlock()
						gameItem.SetTitle(translations["process_not_running"])
						consoleItem.SetTitle(translations["process_not_running"])
						gamename = ""
						currentGame = ""
						currentConsole = ""
					}
					lastState = currentState
				}

				if currentState {
					configMutex.RLock()
					currentLplPath := filepath.Join(config.RetroarchPath, "content_history.lpl")
					currentSavePath := config.SavePath
					if currentSavePath == "" {
						currentSavePath = savePath
					}
					configMutex.RUnlock()

					newGameLine, err1 := findFirstLine(currentLplPath, `"label":`)
					newCoreLine, err2 := findFirstLine(currentLplPath, `"db_name":`)
					if err1 != nil || err2 != nil {
						log.Printf("Error reading file: %v, %v", err1, err2)
						time.Sleep(1 * time.Second)
						continue
					}
					newGamename, newGameFullName := extractValue(newGameLine, `"label": "`)
					shortCoreName, _ := extractValue(newCoreLine, `"db_name": "`)
					consoleName := cut(shortCoreName, ".", 0)

					if gamename != newGamename {
						configMutex.RLock()
						if config.OutputToFiles {
							writeOutputFiles(currentSavePath, newGamename, consoleName, config.SaveToOneFile)
						}
						configMutex.RUnlock()
						gameItem.SetTitle("Game: " + newGamename)
						consoleItem.SetTitle("System: " + consoleName)
						gamename = newGamename
						currentGame = newGamename
						currentGameFullName = newGameFullName
						currentConsole = consoleName
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
				case <-quitItem.ClickedCh:
					log.Println("Exit clicked received")
					systray.Quit()
					return
				}
			}
		}()
	}
}

func openBrowser(url string) error {
	return exec.Command("cmd", "/c", "start", url).Start()
}

func onExit() {
	log.Println(translations["app_exited"])
	os.Exit(0)
}
