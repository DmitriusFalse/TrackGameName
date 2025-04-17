// ВНИМАНИЕ!
// НИ В КОЕМ СЛУЧАЕ НЕ МЕНЯЙТЕ КОД!
// ТУТ НЕЧЕГО МЕНЯТЬ!
// ЭТОТ КОД НУЖЕН ДЛЯ РАБОТЫ СИСТЕМЫ!
// ANTENTION!
// DO NOT CHANGE THE CODE IN ANY WAY!
// THERE'S NOTHING TO CHANGE!
// THIS CODE IS NEEDED FOR THE SYSTEM TO WORK!

// Обновление информации о процессе
function updateProcessInfo() {
    let pid = document.getElementById('process-select').value;
    if (pid) {
        socket.send(JSON.stringify({ type: "get_data", screen: "settings-games", dataType: "infoProcess", pid: pid }));
    }
}
// Функция удаления (заглушка, нужно реализовать серверную часть)
function deleteTemplate(processName) {
    if (confirm(confirmText)) {
        const message = {
            type: "delete",
            dataType: "deleteGameTemplate",
            screen: "settings-games",
            processName: processName,
        };
        socket.send(JSON.stringify(message));
        console.log('Delete template:', processName);
    }
}
window.onload = function() {

    connectWebSocket()

    function connectWebSocket() {
        socket = new WebSocket(`ws://localhost:`+port+`/startport`);

        socket.onopen = () => {
            console.log("✅ WebSocket подключён");

            socket.send(JSON.stringify({ type: "register", screen: "settings-games" }));
            socket.send(JSON.stringify({ type: "get_data", screen: "settings-games", dataType: "gameTemplates" }));
        };

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            switch (data.type) {
                case "gameTemplates":
                    templates = data.payload
                    loadTemplates(templates);
                    break;
                case "infoProcess":
                    document.getElementById('process-name-display').value = data.payload["name"];
                    document.getElementById('window-title-display').value = data.payload["title"] || '';
                    break;

                case "processes":
                    processes = data.payload
                    loadProcesses(processes);
                    break;
                case "refresh":
                    socket.send(JSON.stringify({ type: "get_data", screen: "settings-games", dataType: "gameTemplates" }));
                    break;
            }
        };

        socket.onerror = (err) => {
            console.warn("⚠️ WebSocket ошибка:", err);
        };

        socket.onclose = () => {
            console.warn("❌ WebSocket отключён. Попытка переподключения через", reconnectDelay / 1000, "сек.");
            setTimeout(connectWebSocket, reconnectDelay);
        };
    }



    // Показ формы добавления
    document.getElementById('add-template-btn').addEventListener('click', function() {
        document.getElementById('add-template-form').style.display = 'block';
        socket.send(JSON.stringify({ type: "get_data", screen: "settings-games", dataType: "processes" }));
    });

    document.getElementById('closed').addEventListener('click', function() {
        document.getElementById('add-template-form').style.display = 'none';
    });
    // Загрузка списка процессов
    function loadProcesses(processes){
        let select = document.getElementById('process-select');
        select.innerHTML = '<option value="">'+chooseProcessText+'</option>';
        processes.forEach(proc => {
            let option = document.createElement('option');
            option.value = proc.pid;
            option.text = proc.name;
            select.appendChild(option);
        });
        new SlimSelect({
            select: '#process-select'
        })
    }


    // Функция загрузки шаблонов
    function loadTemplates(templates) {
        const tbody = document.getElementById('templates-list');
        tbody.innerHTML = ''; // Очищаем таблицу

        const start = (currentPage - 1) * rowsPerPage;
        const end = start + rowsPerPage;
        const pageTemplates = templates.slice(start, end);

        pageTemplates.forEach(tmpl => {
            const row = document.createElement('tr');
            row.innerHTML = `
      <td>${tmpl.process_name}</td>
      <td>${tmpl.window_title}</td>
      <td>${tmpl.named_titles ? '<img src="/thumbnails/' + tmpl.named_titles + '" width="50">' : ''}</td>
      <td>${tmpl.named_boxarts ? '<img src="/thumbnails/' + tmpl.named_boxarts + '" width="50">' : ''}</td>
      <td class="last-td"><span class="submit-button" onclick="deleteTemplate('${tmpl.process_name}')">${buttonDeleteText}</span></td>
      `;
            tbody.appendChild(row);
        });

        updatePagination(templates.length);
    }
    // Функция обновления пагинации
    function updatePagination(totalRows) {
        const totalPages = Math.ceil(totalRows / rowsPerPage);
        let pagination = document.getElementById('pagination');
        if (pagination) {
            pagination.innerHTML = '';
        }else{
            pagination = document.createElement('div');
            pagination.id = 'pagination';
            document.getElementById('templates-list').appendChild(pagination); // Добавляем в конец body или в нужный контейнер

        }

        for (let i = 1; i <= totalPages; i++) {
            const button = document.createElement('button');
            button.textContent = i;
            button.disabled = i === currentPage;
            button.addEventListener('click', () => {
                currentPage = i;
                loadTemplates(templates);
            });
            pagination.appendChild(button);
        }
    }

    document.getElementById('template-form').addEventListener('submit', function(e) {
        e.preventDefault(); // Отменяем стандартное поведение формы

        let formData = new FormData(this);
        let formObject = {};
        let filePromises = [];
        formData.forEach((value, key) => {
            if (!(value instanceof File)) {
                formObject[key] = value; // Добавляем обычные данные формы
            }
        });
        // Преобразуем данные формы в объект
        formData.forEach((value, key) => {
            if (value instanceof File && value.size > 0) { // Проверяем, что это файл и он не пустой
                let reader = new FileReader();
                let filePromise = new Promise((resolve, reject) => {
                    reader.onload = function(event) {
                        console.log(value);
                        socket.send(JSON.stringify({
                            type: "saveData",
                            screen: "settings-games",
                            dataType: "saveFile",
                            name: formObject["window_title"],
                            imgType: key,
                            fileName: value.name,
                            fileType: value.type,
                            fileData: event.target.result // Base64-данные файла
                        }));
                        resolve(); // Уведомляем о завершении обработки файла
                    };
                    reader.onerror = reject; // Обрабатываем ошибку чтения файла
                });
                reader.readAsDataURL(value); // Преобразуем файл в Base64
                filePromises.push(filePromise); // Добавляем промис в массив
            }
        });

        // Ждем завершения всех операций с файлами
        Promise.all(filePromises).then(() => {

            delete formObject["process_name"];
            const message = {
                type: "saveData",
                screen: "settings-games",
                dataType: "saveProcess",
                dataForm: formObject,
            };
            socket.send(JSON.stringify(message));
            document.getElementById('add-template-form').style.display = 'none';  // Закрываем форму
        }).catch(error => {
            console.error("Ошибка обработки файлов:", error);
        });
    });
    // Загрузка шаблонов при загрузке страницы


};