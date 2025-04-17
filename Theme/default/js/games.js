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
            socket.send(JSON.stringify({ type: "register", screen: "game" }));
        };

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.type === "update" && data.screen === "game") {
                if (data.payload.game !== lastGame) {
                    document.getElementById('game').textContent = data.payload.game;
                    lastGame = data.payload.game;
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