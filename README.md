# TrackGameName

## English

### Overview
TrackGameName is a lightweight Windows tool designed to track and output RetroArch game data for use in various applications, such as OBS Studio. It displays the current game, system, and thumbnails via a web interface and supports optional file output for text sources.

### Features
- Real-time monitoring of RetroArch to display the current game and system.
- Endpoints(Web address for Widgets) (`/game`, `/system`, `/all`, `/thumbnails`) for seamless integration with applications like OBS Studio.
- Game thumbnails with customizable sizes (e.g., `200x200`, `200x`, `x200`, or original).
- Optional text file output (`game.txt`, `console.txt`, or `output.txt`) for application integration.
- Configurable settings: RetroArch path, thumbnails folder, refresh interval, theme, and more.
- System tray icon showing game/system info with an autorun option.

### Requirements
- Windows operating system.
- RetroArch installed with a valid `content_history.lpl` file.

### Installation
1. Download the latest installer from the [Releases](https://github.com/DmitriusFalse/TrackGameName/releases) page (`trackgamename-vX.X.X-setup.exe`).
2. Run the installer and follow the on-screen instructions.
3. Launch the program and open `http://localhost:3489/` in your browser to access the web interface.
   
### Configuration
Configure the program via the web interface at `/settings`, accessible at `http://localhost:<web_port>/settings` (default port: `3489`). Here you can set:

- **RetroArch Path** (`retroarch_path`):  
  Path to your RetroArch installation (e.g., `C:\RetroArch-Win64`). This is where the program looks for `content_history.lpl` to track the current game.
- **Save Path** (`save_path`):  
  Directory where output files (e.g., `game.txt`, `console.txt`) and theme/system folders are stored. If left empty, defaults to the current working directory.
- **Save to One File** (`save_to_one_file`):  
  Toggle to save game and console info into a single `output.txt` file (e.g., `Nintendo: Super Mario Bros`) instead of separate `game.txt` and `console.txt` files.
- **Autorun** (`autorun`):  
  Enable or disable automatic startup with Windows. When enabled, the program is added to the Windows registry autorun list.
- **Output to Files** (`output_to_files`):  
  Toggle to enable or disable writing game and console info to text files in the save path.
- **Web Port** (`web_port`):  
  Port for the web server (e.g., `3489`). Must be a number between 1 and 65535. Restart the program for changes to take effect.
- **System Icon** (`system_icon`):  
  Choose how system icons are displayed (0 = disabled, 1 = small, 2 = large). Icons are sourced from the `systems` section of the configuration.
- **Refresh Interval** (`refresh_interval`):  
  Time in seconds between updates of game and console info (e.g., `20`). Must be at least 1 second.
- **Theme** (`theme`):  
  Select the visual theme for the web interface (e.g., `default`). Available themes are detected from the `Theme` folder in the save path.
- **Language** (`language`):  
  Choose the interface language (e.g., `en` for English). Available languages are detected from `.json` files in the `lang` folder.
- **Thumbnails Path** (`thumbnails_path`):  
  Directory containing your thumbnails (e.g., `C:\RetroArch-Win64\thumbnails`). Thumbnails must follow the RetroArch structure.
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
- Эндпоинты (Веб-адрес для Виджетов) (`/game`, `/system`, `/all`, `/thumbnails`) для удобной интеграции с приложениями, такими как OBS Studio.
- Отображение миниатюр игр с настраиваемым размером (например, `200x200`, `200x`, `x200` или оригинал).
- Опциональный вывод в текстовые файлы (`game.txt`, `console.txt` или `output.txt`) для интеграции с приложениями.
- Настраиваемые параметры: путь к RetroArch, папка с миниатюрами, интервал обновления, тема и др.
- Иконка в системном трее с информацией об игре/системе и опцией автозапуска.

### Требования
- Операционная система Windows.
- Установленный RetroArch с файлом `content_history.lpl`.

### Установка
1. Скачайте последнюю версию установщика со страницы [Releases](https://github.com/DmitriusFalse/TrackGameName/releases) (`trackgamename-vX.X.X-setup.exe`).
2. Запустите установщик и следуйте инструкциям на экране.
3. Запустите программу и откройте `http://localhost:3489/` в браузере для доступа к веб-интерфейсу.

## Настройка

Настройте программу через веб-интерфейс на странице `/settings`, доступной по адресу `http://localhost:<web_port>/settings` (порт по умолчанию: `3489`). Здесь вы можете указать:

- **Путь к RetroArch** (`retroarch_path`):  
  Путь к установленному RetroArch (например, `C:\RetroArch-Win64`). Здесь программа ищет файл `content_history.lpl` для отслеживания текущей игры.
- **Путь для сохранения** (`save_path`):  
  Директория, где сохраняются выходные файлы (например, `game.txt`, `console.txt`) и папки тем/систем. Если оставить пустым, используется текущая рабочая директория.
- **Сохранять в один файл** (`save_to_one_file`):  
  Переключатель для сохранения информации об игре и консоли в один файл `output.txt` (например, `Nintendo: Super Mario Bros`) вместо отдельных файлов `game.txt` и `console.txt`.
- **Автозапуск** (`autorun`):  
  Включение или отключение автоматического запуска вместе с Windows. При включении программа добавляется в список автозагрузки в реестре Windows.
- **Вывод в файлы** (`output_to_files`):  
  Переключатель для включения или отключения записи информации об игре и консоли в текстовые файлы в указанной директории.
- **Веб-порт** (`web_port`):  
  Порт веб-сервера (например, `3489`). Должен быть числом от 1 до 65535. Для применения изменений требуется перезапуск программы.
- **Иконка системы** (`system_icon`):  
  Выбор отображения иконок систем (0 = отключено, 1 = маленькие, 2 = большие). Иконки берутся из секции `systems` конфигурации.
- **Интервал обновления** (`refresh_interval`):  
  Время в секундах между обновлениями информации об игре и консоли (например, `20`). Минимальное значение — 1 секунда.
- **Тема** (`theme`):  
  Выбор визуальной темы для веб-интерфейса (например, `default`). Доступные темы определяются из папки `Theme` в директории сохранения.
- **Язык** (`language`):  
  Выбор языка интерфейса (например, `en` для английского). Доступные языки определяются из файлов `.json` в папке `lang`.
- **Путь к миниатюрам** (`thumbnails_path`):  
  Директория с миниатюрами (например, `C:\RetroArch-Win64\thumbnails`). Миниатюры должны соответствовать структуре RetroArch.
- **Включить миниатюры** (`enable_thumbnails`):  
  Переключатель для включения или отключения отображения миниатюр в веб-интерфейсе.
- **Размер миниатюр** (`thumbnail_size`):  
  Установка размера отображения миниатюр (например, `200x200`). Формат — `ширинаxвысота` в пикселях. Также можно указать:  
  - Только ширину (например, `200x`) для установки ширины с пропорциональной высотой.  
  - Только высоту (например, `x200`) для установки высоты с пропорциональной шириной.  
  - Оставить пустым или установить `0` для размера по умолчанию.

Разместите изображения миниатюр в указанной директории в структуре RetroArch:  
`<thumbnails_path>\<system>\Named_Titles\<game>.png`.  
Например:
`C:\RetroArch-Win64\thumbnails\Atari - 2600\Named_Titles\Q_bert's Qubes.png`
- Часть `<game>` должна совпадать с полным названием игры из `content_history.lpl`, включая регион и информацию о диске (например, `Armored Core - Master of Arena (USA) (Disc 1)`).

Если миниатюра не найдена, программа отобразит `noimage.png` из папки темы (например, `Theme\default\noimage.png`). Убедитесь, что этот файл присутствует в директории выбранной темы.

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
  - Note: Format: `<thumbnails_path>\<system>\Named_Titles\<game>.png`. The `<game>` part must match the full game name from `content_history.lpl`, including region and disc info (e.g., `Armored Core - Master of Arena (USA) (Disc 1)`)

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
  - Примечание: Формат: `<system>\<game>.png`. Часть `<game>` должна совпадать с полным названием игры из `content_history.lpl`, включая регион и информацию о диске (например, `Armored Core - Master of Arena (USA) (Disc 1)`). Если миниатюра не найдена, программа отобразит `noimage.png` из папки темы (например, `Theme\default\noimage.png`). Убедитесь, что этот файл присутствует в директории выбранной темы.

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
