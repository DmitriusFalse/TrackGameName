// ВНИМАНИЕ!
// НИ В КОЕМ СЛУЧАЕ НЕ МЕНЯЙТЕ КОД!
// ТУТ НЕЧЕГО МЕНЯТЬ!
// ЭТОТ КОД НУЖЕН ДЛЯ РАБОТЫ СИСТЕМЫ!
// ANTENTION!
// DO NOT CHANGE THE CODE IN ANY WAY!
// THERE'S NOTHING TO CHANGE!
// THIS CODE IS NEEDED FOR THE SYSTEM TO WORK!
window.onload = function() {
    socket = new WebSocket(`ws://localhost:${port}/startport`);
    socket.onopen = () => {
        socket.send(JSON.stringify({ type: "register", screen: "system" }));
    };
    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.type === "update" && data.screen === "system") {
            if (data.payload.console !== lastSystem) {
                document.getElementById('system').textContent = data.payload.console;
                lastSystem = data.payload.console;
            }
            if(IconFile !==data.payload.icon){ //Если иконка меняется
                if (data.payload.icon !== "") { //Если иконка не пустая
                    IconFile = data.payload.icon;
                    const imgSystem = document.querySelector('img#icon-system');
                    if (imgSystem) { // проверяем, есть ли вывод картинки да - меняем картинку
                        imgSystem.src =  '/systems/'+data.payload.icon;
                    } else { // нет - создаем вывод и добавляем
                        const img = document.createElement('img');
                        img.src = '/systems/'+data.payload.icon;
                        img.alt = data.payload.console;
                        img.id = "icon-system";
                        const spanSystem = document.getElementById('system');
                        spanSystem.before(img);
                    }
                }else{
                    const imgSystem = document.querySelector('img#icon-system');
                    if (imgSystem) {
                        imgSystem.remove();
                    }
                }
            }

        }
    };
};