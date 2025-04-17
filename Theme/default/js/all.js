// ВНИМАНИЕ!
// НИ В КОЕМ СЛУЧАЕ НЕ МЕНЯЙТЕ КОД!
// ТУТ НЕЧЕГО МЕНЯТЬ!
// ЭТОТ КОД НУЖЕН ДЛЯ РАБОТЫ СИСТЕМЫ!
// ANTENTION!
// DO NOT CHANGE THE CODE IN ANY WAY!
// THERE'S NOTHING TO CHANGE!
// THIS CODE IS NEEDED FOR THE SYSTEM TO WORK!
window.onload = function() {
    connectWebSocket()

    function connectWebSocket() {
        socket = new WebSocket(`ws://localhost:${port}/startport`);

        socket.onopen = () => {
            console.log("✅ WebSocket подключён");

            // Пример регистрации страницы
            socket.send(JSON.stringify({ type: "register", screen: "all" }));
        };

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.type === "update" && data.screen === "all") {
                if (data.payload.game !== lastGame) {
                    document.getElementById('game').textContent = data.payload.game;
                    lastGame = data.payload.game;
                }
                if (data.payload.console !== lastSystem) {
                    document.getElementById('system').textContent = data.payload.console;
                    lastSystem = data.payload.console;
                }

                if(IconFile !==data.payload.icon){
                    if (data.payload.icon !== "") {
                        IconFile = data.payload.icon;
                        const imgSystem = document.querySelector('img#icon-system');
                        if (imgSystem) {
                            imgSystem.src =  '/systems/'+data.payload.icon;
                        } else {
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

        socket.onerror = (err) => {
            console.warn("⚠️ WebSocket ошибка:", err);
        };

        socket.onclose = () => {
            console.warn("❌ WebSocket отключён. Попытка переподключения через", reconnectDelay / 1000, "сек.");
            setTimeout(connectWebSocket, reconnectDelay);
        };
    }

};