# TrackGameName

## English

### Overview
TrackGameName is a lightweight Windows tool designed to track and output RetroArch game data for use in various applications, such as OBS Studio. It displays the current game, system, and thumbnails via a web interface and supports optional file output for text sources.

### Features
- Real-time monitoring of RetroArch to display the current game and system.
- Web interface with endpoints: `/game`, `/system`, `/all`, `/thumbnails`, and `/settings`.
- Game thumbnails with customizable sizes (e.g., `200x200`, `200x`, `x200`, or original).
- Optional text file output (`game.txt`, `console.txt`, or `output.txt`) for application integration.
- Configurable settings: RetroArch path, thumbnails folder, refresh interval, theme, and more.
- System tray icon showing game/system info with an autorun option.
- Donation support via [DonationAlerts](https://www.donationalerts.com/r/ork8bit) with a cute heart icon ❤️.

### Requirements
- Windows operating system.
- RetroArch installed with a valid `content_history.lpl` file.

### Installation
1. Download the latest installer from the [Releases](https://github.com/DmitriusFalse/TrackGameName/releases) page (`trackgamename-vX.X.X-setup.exe`).
2. Run the installer and follow the on-screen instructions.
3. Launch the program and open `http://localhost:3489/` in your browser to access the web interface.

### Building from Source (Optional)
If you prefer to build the executable yourself:
1. Clone the repository:
   ```cmd
   git clone https://github.com/DmitriusFalse/TrackGameName.git
   cd trackgamename
2. Install Go 1.16+.
3. Build the executable:
   ```cmd
   go build -ldflags="-H=windowsgui" -o trackgamename.exe
4. Run the program: trackgamename.exe
   
### Configuration
- Edit `config.ini` to set your RetroArch path, thumbnails folder, and other preferences.
- Place thumbnail images in the specified folder (e.g., `C:\thumbnails\<system>\<game>.png`).
- Alternatively, configure settings via the web interface at `/settings`.

### Support
If you enjoy this tool, please consider supporting me on [DonationAlerts](https://www.donationalerts.com/r/ork8bit) ❤️.

### License
This project is licensed under the [MIT License](LICENSE) - see the file for details.

---

## Русский

### Описание
TrackGameName — это лёгкий инструмент для Windows, который отслеживает и выводит данные игр из RetroArch для использования в различных приложениях, таких как OBS Studio. Программа отображает текущую игру, систему и миниатюры через веб-интерфейс и поддерживает вывод в текстовые файлы.

### Возможности
- Отслеживание RetroArch в реальном времени для показа текущей игры и системы.
- Веб-интерфейс с эндпоинтами: `/game`, `/system`, `/all`, `/thumbnails` и `/settings`.
- Отображение миниатюр игр с настраиваемым размером (например, `200x200`, `200x`, `x200` или оригинал).
- Опциональный вывод в текстовые файлы (`game.txt`, `console.txt` или `output.txt`) для интеграции с приложениями.
- Настраиваемые параметры: путь к RetroArch, папка с миниатюрами, интервал обновления, тема и др.
- Иконка в системном трее с информацией об игре/системе и опцией автозапуска.
- Поддержка через [DonationAlerts](https://www.donationalerts.com/r/ork8bit) с милым сердечком ❤️.

### Требования
- Операционная система Windows.
- Установленный RetroArch с файлом `content_history.lpl`.

### Установка
1. Скачайте последнюю версию установщика со страницы [Releases](https://github.com/DmitriusFalse/TrackGameName/releases) (`trackgamename-vX.X.X-setup.exe`).
2. Запустите установщик и следуйте инструкциям на экране.
3. Запустите программу и откройте `http://localhost:3489/` в браузере для доступа к веб-интерфейсу.

### Сборка из исходников (по желанию)
Если хотите собрать исполняемый файл самостоятельно:
1. Склонируйте репозиторий:
   ```cmd
   git clone https://github.com/DmitriusFalse/TrackGameName.git
   cd trackgamename
2. Установите Go 1.16+.
3. Скомпилируйте программу:
   ```cmd
   go build -ldflags="-H=windowsgui" -o trackgamename.exe
4. Запустите: trackgamename.exe

### Настройка
- Отредактируйте `config.ini`, указав путь к RetroArch, папку с миниатюрами и другие настройки.
- Разместите изображения миниатюр в указанной папке (например, `C:\thumbnails\<система>\<игра>.png`).
- Также можно настроить параметры через веб-интерфейс на странице `/settings`.

### Поддержка
Если вам понравился этот инструмент, поддержите меня на [DonationAlerts](https://www.donationalerts.com/r/ork8bit) ❤️.

### Лицензия
Проект распространяется под лицензией [MIT License](LICENSE) — подробности в файле.

## Web Interface Settings (`/settings`)

Below is a guide to all settings available on the `/settings` page and how they correspond to parameters in `config.ini`.

- **RetroArch Path**
  - Description: Specifies the path to the RetroArch installation folder. Used to locate `content_history.lpl` for detecting the current game and system.
  - Example value: `C:\RetroArch-Win64`
  - Matches in `config.ini`: `retroarch_path`
  - Default: `C:\RetroArch-Win64`
  - Note: Ensure the path is correct, or game tracking won’t work.

- **Save Path**
  - Description: Path to the folder where text files (`game.txt`, `console.txt`, or `output.txt`) with game/system data are saved.
  - Example value: `D:\Folder`
  - Matches in `config.ini`: `save_path`
  - Default: Empty string (program’s current directory)
  - Note: Leave empty to use the program’s folder.

- **Save to One File**
  - Description: Enables saving game and system data to a single file (`output.txt`) instead of two separate files (`game.txt` and `console.txt`).
  - Values: Checkbox (on/off)
  - Matches in `config.ini`: `save_to_one_file` (`true`/`false`)
  - Default: Off (`false`)
  - Note: When enabled, `output.txt` contains `<system>: <game>`.

- **Autorun**
  - Description: Enables the program to start automatically with Windows.
  - Values: Checkbox (on/off)
  - Matches in `config.ini`: `autorun` (`true`/`false`)
  - Default: Off (`false`)
  - Note: Adds the program to the Windows registry (`Run`).

- **Output to Files**
  - Description: Enables writing game and system data to text files for use in OBS or other applications.
  - Values: Checkbox (on/off)
  - Matches in `config.ini`: `output_to_files` (`true`/`false`)
  - Default: On (`true`)
  - Note: If disabled, files aren’t created, but data is available via the web interface.

- **Web Port**
  - Description: Port for the program’s web server, used to access the interface in a browser.
  - Example value: `3489`
  - Matches in `config.ini`: `web_port`
  - Default: `3489`
  - Note: Restart the program after changing. Ensure the port is free.

- **System Icon**
  - Description: Controls the display of system icons on `/system` and `/all` pages.
  - Values: Number (`0` — no icons, `1` — icon with text, `2` — icon only)
  - Matches in `config.ini`: `system_icon`
  - Default: `0`
  - Note: Icons are defined in the `[systems]` section of `config.ini`.

- **Refresh Interval**
  - Description: Refresh interval (in seconds) for web interface pages, except the main page.
  - Example value: `20`
  - Matches in `config.ini`: `refresh_interval`
  - Default: `20`
  - Note: The main page requires manual refresh.

- **Theme**
  - Description: Selects the web interface theme.
  - Example value: `default`, `dark`
  - Matches in `config.ini`: `theme`
  - Default: `default`
  - Note: Themes are stored in the `Theme` folder.

- **Language**
  - Description: Selects the interface language of the program.
  - Example value: `en` (English), `ru` (Russian)
  - Matches in `config.ini`: `language`
  - Default: `en`
  - Note: Translations are stored in the `lang` folder.

- **Thumbnails Path**
  - Description: Path to the folder with game thumbnails for display on `/thumbnails` and the main page.
  - Example value: `D:\thumbnails`
  - Matches in `config.ini`: `thumbnails_path`
  - Default: Empty string (thumbnails disabled)
  - Note: Format: `<system>\<game>.png`.

- **Enable Thumbnails**
  - Description: Enables searching and displaying game thumbnails.
  - Values: Checkbox (on/off)
  - Matches in `config.ini`: `enable_thumbnails` (`true`/`false`)
  - Default: Off (`false`)
  - Note: Requires a valid `thumbnails_path`.

- **Thumbnail Size**
  - Description: Size of thumbnails in pixels on `/thumbnails` and the main page.
  - Example value:
    - `200x200` — fixed size of 200 pixels wide and 200 pixels tall.
    - `200x` — width of 200 pixels, height proportional to the original.
    - `x200` — height of 200 pixels, width proportional to the original.
    - `0` — original image size without scaling.
  - Matches in `config.ini`: `thumbnail_size`
  - Default: `0`
  - Note: Invalid values are ignored.
 
    ## Настройки через веб-интерфейс (`/settings`)

Ниже описаны все поля настроек на странице `/settings` и их соответствие параметрам в файле `config.ini`.

- **Путь к RetroArch (`RetroArch Path`)**
  - Описание: Указывает путь к папке, где установлен RetroArch. Программа использует этот путь для поиска файла `content_history.lpl`, чтобы определять текущую игру и систему.
  - Пример значения: `C:\RetroArch-Win64`
  - Соответствие в `config.ini`: `retroarch_path`
  - По умолчанию: `C:\RetroArch-Win64`
  - Примечание: Убедитесь, что путь указан верно, иначе отслеживание игр не будет работать.

- **Путь сохранения (`Save Path`)**
  - Описание: Путь к папке, куда сохраняются текстовые файлы (`game.txt`, `console.txt` или `output.txt`) с данными о текущей игре и системе.
  - Пример значения: `D:\Папка`
  - Соответствие в `config.ini`: `save_path`
  - По умолчанию: Пустая строка (текущая директория программы)
  - Примечание: Оставьте пустым, чтобы использовать папку программы.

- **Сохранять в один файл (`Save to One File`)**
  - Описание: Включает сохранение данных об игре и системе в один файл (`output.txt`) вместо двух отдельных (`game.txt` и `console.txt`).
  - Значения: Чекбокс (вкл/выкл)
  - Соответствие в `config.ini`: `save_to_one_file` (`true`/`false`)
  - По умолчанию: Выключено (`false`)
  - Примечание: При включении в `output.txt` будет строка вида `<система>: <игра>`.

- **Автозапуск (`Autorun`)**
  - Описание: Включает автозапуск программы при старте Windows.
  - Значения: Чекбокс (вкл/выкл)
  - Соответствие в `config.ini`: `autorun` (`true`/`false`)
  - По умолчанию: Выключено (`false`)
  - Примечание: Добавляет программу в реестр Windows (`Run`).

- **Вывод в файлы (`Output to Files`)**
  - Описание: Включает запись данных об игре и системе в текстовые файлы для использования в OBS или других программах.
  - Значения: Чекбокс (вкл/выкл)
  - Соответствие в `config.ini`: `output_to_files` (`true`/`false`)
  - По умолчанию: Включено (`true`)
  - Примечание: Если выключено, файлы не создаются, но данные доступны в веб-интерфейсе.

- **Веб-порт (`Web Port`)**
  - Описание: Порт для веб-сервера программы, используемый для доступа через браузер.
  - Пример значения: `3489`
  - Соответствие в `config.ini`: `web_port`
  - По умолчанию: `3489`
  - Примечание: После изменения нужно перезапустить программу. Порт должен быть свободен.

- **Иконка системы (`System Icon`)**
  - Описание: Управляет отображением иконок систем на страницах `/system` и `/all`.
  - Значения: Число (`0` — без иконок, `1` — иконка с текстом, `2` — только иконка)
  - Соответствие в `config.ini`: `system_icon`
  - По умолчанию: `0`
  - Примечание: Иконки задаются в секции `[systems]` в `config.ini`.

- **Интервал обновления (`Refresh Interval`)**
  - Описание: Интервал обновления страниц веб-интерфейса (кроме главной) в секундах.
  - Пример значения: `20`
  - Соответствие в `config.ini`: `refresh_interval`
  - По умолчанию: `20`
  - Примечание: Главная страница обновляется вручную.

- **Тема (`Theme`)**
  - Описание: Выбор темы оформления веб-интерфейса.
  - Пример значения: `default`, `dark`
  - Соответствие в `config.ini`: `theme`
  - По умолчанию: `default`
  - Примечание: Темы находятся в папке `Theme`.

- **Язык (`Language`)**
  - Описание: Выбор языка интерфейса программы.
  - Пример значения: `en` (английский), `ru` (русский)
  - Соответствие в `config.ini`: `language`
  - По умолчанию: `en`
  - Примечание: Переводы хранятся в папке `lang`.

- **Путь к миниатюрам (`Thumbnails Path`)**
  - Описание: Путь к папке с миниатюрами игр для отображения на `/thumbnails` и главной странице.
  - Пример значения: `D:\thumbnails`
  - Соответствие в `config.ini`: `thumbnails_path`
  - По умолчанию: Пустая строка (миниатюры отключены)
  - Примечание: Формат: `<система>\<игра>.png`.

- **Включить миниатюры (`Enable Thumbnails`)**
  - Описание: Включает поиск и отображение миниатюр игр.
  - Значения: Чекбокс (вкл/выкл)
  - Соответствие в `config.ini`: `enable_thumbnails` (`true`/`false`)
  - По умолчанию: Выключено (`false`)
  - Примечание: Требуется указать `thumbnails_path`.

- **Размер миниатюр (`Thumbnail Size`)**
  - Описание: Размер миниатюр в пикселях на страницах `/thumbnails` и главной.
  - Пример значения:
    - `200x200` — фиксированный размер 200 пикселей по ширине и высоте.
    - `200x` — ширина 200 пикселей, высота пропорциональна оригиналу.
    - `x200` — высота 200 пикселей, ширина пропорциональна оригиналу.
    - `0` — оригинальный размер изображения без изменений.
  - Соответствие в `config.ini`: `thumbnail_size`
  - По умолчанию: `0`
  - Примечание: Некорректные значения игнорируются.
