# TrackGameName

## English

### Overview
TrackGameName is a lightweight Windows tool designed to track and output RetroArch game data for use in various applications, such as OBS Studio. It displays the current game, system, and thumbnails via a web interface and supports optional file output for text sources.

### Features
- Real-time monitoring of RetroArch and Windows Game to display the current game and system.
- Widgets for seamless integration with applications like OBS Studio.
- Game thumbnails with customizable sizes (e.g., `200x200`, `200x`, `x200`, or original).
- Optional text file output (`game.txt`, `console.txt`, or `output.txt`) for application integration.
- System tray icon showing game/system info with an autorun option.

### Requirements
- Windows operating system.
- RetroArch installed with a valid `content_history.lpl` file. (No retroarch is required to track Windows games.)


### Installation
1. Download the latest installer from the [Releases](https://github.com/DmitriusFalse/TrackGameName/releases) page (`trackgamename-vX.X.X-setup.exe`).
2. Run the installer and follow the on-screen instructions.
3. Launch the program and open `http://localhost:3489/` in your browser to access the web interface.
   
### Configuration
Configure the program via the web interface at `/settings`, accessible at `http://localhost:<web_port>/settings` (default port: `3489`). Here you can set:
- **Autorun** (`autorun`):  
  Enable or disable automatic startup with Windows. When enabled, the program is added to the Windows registry autorun list.
- **Web Port** - Web server port `http://localhost:<web_port>/settings`
- **System Icon** - You can display the console icon next to the console name. 0 - no icon, 1 - icon with text, 2 - icon only
  - Open config.ini in the program folder
  - The [systems] section has already added an example for Nintendo - the Nintendo Entertainment System `Nintendo - Nintendo Entertainment System = nes.png`
  - Write THE NAME OF THE CONSOLE = THE NAME OF THE IMAGE.png
  - Place the file in the systems folder in the program folder.
  - And restart TrackGameName
- **RetroArch Path** (`retroarch_path`):  
  Path to your RetroArch installation (e.g., `C:\RetroArch-Win64`). This is where the program looks for `content_history.lpl` to track the current game.
- **Save Path** (`save_path`):  
  Directory where output files (e.g., `game.txt`, `console.txt`) and theme/system folders are stored. If left empty, defaults to the current working directory.
- **Save to One File** (`save_to_one_file`):  
  Toggle to save game and console info into a single `output.txt` file (e.g., `Nintendo: Super Mario Bros`) instead of separate `game.txt` and `console.txt` files.
- **Output to Files** (`output_to_files`):  
  Toggle to enable or disable writing game and console info to text files in the save path.
- **Theme** (`theme`):  
  Select the visual theme for the web interface (e.g., `default`). Available themes are detected from the `Theme` folder in the save path.
- **Language** (`language`):  
  Choose the interface language (e.g., `en` for English). Available languages are detected from `.json` files in the `lang` folder.
- **Thumbnails Path** (`thumbnails_path`):  
  Directory containing your thumbnails (e.g., `C:\RetroArch-Win64\thumbnails`). Thumbnails must follow the RetroArch structure.
  - Example: if the content in the playlist is named Q*bert's Qubes, then the thumbnails should be named Q_bert's Qubes.png and stored in the following ways:
  - thumbnails/ Atari — 2600/ Named_Boxarts/ Q_bert's Qubes.png
  - thumbnails/ Atari — 2600/ Named_Snaps/ Q_bert's Qubes.png
  - thumbnails/ Atari — 2600/ Named_Titles/ Q_bert's Qubes.png
    
- **Enable Thumbnails** (`enable_thumbnails`):  
  Toggle thumbnail display on or off in the web interface.
- **Thumbnail Size** (`thumbnail_size`):  
  Set the display size for thumbnails (e.g., `200x200`). Format is `widthxheight` in pixels. You can also specify:
  - Only width (e.g., `200x`) to set width while keeping height proportional.
  - Only height (e.g., `x200`) to set height while keeping width proportional.
  - Leave blank or set to `0` for default size.

Place thumbnail images in the specified thumbnails folder using the RetroArch structure:  
`<thumbnails_path>\<system>\Named_Titles\<game>.png`.  
For example:
`C:\RetroArch-Win64\thumbnails\Atari - 2600\Named_Titles\Q_bert's Qubes.png`
- The `<game>` part must match the full game name from `content_history.lpl`, including region and disc info (e.g., `Armored Core - Master of Arena (USA) (Disc 1)`).

If a thumbnail is not found, the program will display `noimage.png` from the theme folder (e.g., `Theme\default\noimage.png`). Ensure this file exists in your selected theme directory.

# Theming

To create your own visual theme:
## 1. Create a Theme Folder
Create a subfolder inside the `Theme` directory.  
You can name it anything you like.
## 2. Create a Stylesheet
Inside the new folder, create a `styles.css` file.  
It's recommended to copy `styles.css` from one of the existing themes and modify only the parts you want to change. Leave the rest untouched for compatibility.
## 3. (Optional) Modify Widget Templates
If you need to change the widget layout:
- Copy the required `.html` templates from the `default` theme folder.
- Edit them as needed.
> ⚠️ **Important:**  
> Modifying templates is not recommended unless you're confident in what you're doing.  
> Templates must retain certain class names and IDs that are necessary for scripts to function correctly.
## Recommendation
In most cases, modifying only the `styles.css` file is enough to achieve custom theming.

### Support
If you enjoy this tool, please consider supporting me on [DonationAlerts](https://www.donationalerts.com/r/ork8bit) ❤️.

### License
This project is licensed under the [MIT License](LICENSE) - see the file for details.

---

## Русский

### Обзор
TrackGameName — это легковесный инструмент для Windows, предназначенный для отслеживания и вывода данных игр RetroArch для использования в различных приложениях, таких как OBS Studio. Он отображает текущую игру, систему и миниатюры через веб-интерфейс и поддерживает опциональный вывод в текстовые файлы для источников текста.

### Особенности
- Мониторинг RetroArch и игр Windows в реальном времени для отображения текущей игры и системы.
- Виджеты для бесшовной интеграции с приложениями, такими как OBS Studio.
- Миниатюры игр с настраиваемыми размерами (например, `200x200`, `200x`, `x200` или оригинальный).
- Опциональный вывод в текстовые файлы (`game.txt`, `console.txt` или `output.txt`) для интеграции с приложениями.
- Иконка в системном трее, отображающая информацию об игре/системе, с опцией автозапуска.

### Требования
- Операционная система Windows.
- Установленный RetroArch с действительным файлом `content_history.lpl`. (RetroArch не требуется для отслеживания игр Windows.)

### Установка
1. Скачайте последнюю версию установщика со страницы [Releases](https://github.com/DmitriusFalse/TrackGameName/releases) (`trackgamename-vX.X.X-setup.exe`).
2. Запустите установщик и следуйте инструкциям на экране.
3. Запустите программу и откройте `http://localhost:3489/` в браузере для доступа к веб-интерфейсу.

### Настройка
Настройте программу через веб-интерфейс на странице `/settings`, доступной по адресу `http://localhost:<web_port>/settings` (порт по умолчанию: `3489`). Здесь можно настроить:

- **Автозапуск** (`autorun`):  
  Включить или отключить автоматический запуск с Windows. При включении программа добавляется в список автозагрузки реестра Windows.
- **Веб-порт** (`web_port`):  
  Порт веб-сервера `http://localhost:<web_port>/settings`.
- **Иконка системы** (`system_icon`):  
  Можно отображать иконку консоли рядом с её названием. 0 — без иконки, 1 — иконка с текстом, 2 — только иконка.
  - Откройте файл `config.ini` в папке программы.
  - В секции `[systems]` уже добавлен пример для Nintendo — Nintendo Entertainment System: `Nintendo - Nintendo Entertainment System = nes.png`.
  - Укажите ИМЯ КОНСОЛИ = ИМЯ ИЗОБРАЖЕНИЯ.png.
  - Поместите файл в папку `systems` в директории программы.
  - Перезапустите TrackGameName.
- **Путь к RetroArch** (`retroarch_path`):  
  Путь к установке RetroArch (например, `C:\RetroArch-Win64`). Здесь программа ищет файл `content_history.lpl` для отслеживания текущей игры.
- **Путь сохранения** (`save_path`):  
  Директория, куда сохраняются выходные файлы (например, `game.txt`, `consoleaparte.txt`) и папки тем/систем. Если оставить пустым, используется текущая рабочая директория.
- **Сохранение в один файл** (`save_to_one_file`):  
  Включить, чтобы сохранять информацию об игре и консоли в один файл `output.txt` (например, `Nintendo: Super Mario Bros`) вместо отдельных файлов `game.txt` и `console.txt`.
- **Вывод в файлы** (`output_to_files`):  
  Включить или отключить запись информации об игре и консоли в текстовые файлы в указанной директории.
- **Тема** (`theme`):  
  Выберите визуальную тему для веб-интерфейса (например, `default`). Доступные темы определяются из папки `Theme` в пути сохранения.
- **Язык** (`language`):  
  Выберите язык интерфейса (например, `en` для английского). Доступные языки определяются из файлов `.json` в папке `lang`.
- **Путь к миниатюрам** (`thumbnails_path`):  
  Директория, содержащая миниатюры (например, `C:\RetroArch-Win64\thumbnails`). Миниатюры должны следовать структуре RetroArch.
  - Пример: если в плейлисте игра называется Q*bert's Qubes, то миниатюры должны называться `Q_bert's Qubes.png` и храниться следующим образом:
    - `thumbnails/Atari — 2600/Named_Boxarts/Q_bert's Qubes.png`
    - `thumbnails/Atari — 2600/Named_Snaps/Q_bert's Qubes.png`
    - `thumbnails/Atari — 2600/Named_Titles/Q_bert's Qubes.png`
- **Включить миниатюры** (`enable_thumbnails`):  
  Включить или отключить отображение миниатюр в веб-интерфейсе.
- **Размер миниатюр** (`thumbnail_size`):  
  Установите размер отображения миниатюр (например, `200x200`). Формат: `ширинаxвысота` в пикселях. Также можно указать:
  - Только ширину (например, `200x`), чтобы задать ширину с пропорциональной высотой.
  - Только высоту (например, `x200`), чтобы задать высоту с пропорциональной шириной.
  - Оставить пустым или установить `0` для размера по умолчанию.

Поместите изображения миниатюр в указанную папку, следуя структуре RetroArch:  
`<thumbnails_path>\<system>\Named_Titles\<game>.png`.  
Пример:  
`C:\RetroArch-Win64\thumbnails\Atari - 2600\Named_Titles\Q_bert's Qubes.png`

- Часть `<game>` должна точно совпадать с полным названием игры из `content_history.lpl`, включая регион и информацию о диске (например, `Armored Core - Master of Arena (USA) (Disc 1)`).

Если миниатюра не найдена, программа отобразит `noimage.png` из папки темы (например, `Theme\default\noimage.png`). Убедитесь, что этот файл существует в директории выбранной темы.

# Темизация

Чтобы создать собственное визуальное оформление:

## 1. Создание папки темы
Создайте подпапку в директории `Theme`.  
Название может быть произвольным.
## 2. Создание файла стилей
Внутри новой папки создайте файл `styles.css`.  
Рекомендуется скопировать `styles.css` из одной из готовых тем и изменить только те параметры, которые вам нужны. Остальное лучше оставить без изменений.
## 3. (Необязательно) Изменение шаблонов виджетов
Если вам нужно изменить HTML-разметку виджетов:
- Скопируйте нужные `.html`-шаблоны из папки темы `default`.
- Отредактируйте их по своему усмотрению.
> ⚠️ **Важно:**  
> Изменение шаблонов не рекомендуется, если вы не уверены в своих действиях.  
> В шаблонах обязательно должны сохраняться определённые классы и идентификаторы — они необходимы для корректной работы скриптов.
## Рекомендация
В большинстве случаев достаточно изменения только `styles.css` для создания уникального оформления.


### Поддержка
Если вам понравился этот инструмент, пожалуйста, поддержите меня на [DonationAlerts](https://www.donationalerts.com/r/ork8bit) ❤️.

### Лицензия
Этот проект распространяется под [лицензией MIT](LICENSE) — подробности смотрите в файле.
