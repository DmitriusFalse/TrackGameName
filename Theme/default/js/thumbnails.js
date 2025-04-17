// ВНИМАНИЕ!
// НИ В КОЕМ СЛУЧАЕ НЕ МЕНЯЙТЕ КОД!
// ТУТ НЕЧЕГО МЕНЯТЬ!
// ЭТОТ КОД НУЖЕН ДЛЯ РАБОТЫ СИСТЕМЫ!
// ANTENTION!
// DO NOT CHANGE THE CODE IN ANY WAY!
// THERE'S NOTHING TO CHANGE!
// THIS CODE IS NEEDED FOR THE SYSTEM TO WORK!
window.onload = function() {

    const container = document.querySelector('.thumbnail-container');

    lastThumbnails = container.querySelectorAll("img");
    startThumbnailInterval(container)
    setSizeContainer(container, width, height, container.querySelector("img"))
    connectWebSocket()

    function startThumbnailInterval(container) {
        if (intervalId) {
            clearInterval(intervalId);
        }
        // Запускаем новый интервал
        intervalId = setInterval(() => {
            switchToNextThumbnail(container);
        }, intervalDelay);
    }
    function setSizeContainer(container, width, height, img) {
        let finalWidth = 0;
        let finalHeight = 0;

        let aspectRatio = img.naturalWidth / img.naturalHeight;
        if (height >0 && width === 0) {
            finalHeight = height + 'px';
            finalWidth = Math.round(aspectRatio * height)+'px';
        }else if(height === 0 && width > 0){
            finalHeight = Math.round(aspectRatio*width)+'px';
            finalWidth  = width + 'px';
        }else if(height > 0 && width > 0){
            finalHeight = height + 'px';
            finalWidth  = width + 'px';
        }else{
            finalHeight = img.naturalHeight + 'px';
            finalWidth  = img.naturalWidth + 'px';
        }

        container.style.minHeight = finalHeight;
        container.style.minWidth = finalWidth;
        container.style.height = finalHeight;
        container.style.width = finalWidth;

    }
    document.head.insertAdjacentHTML("beforeend", `
    <style>
        .container.thumbnail-container {
            position: relative;
            display: grid;
            justify-items: center;
            overflow: hidden;
            transition: height ${fadeDuration / 1000}s ${fadeType};
        }
        .thumbnail {
            width: 100%;
            height: 100%;
            position: absolute;
            object-fit: contain;
            transition: opacity ${fadeDuration / 1000}s ${fadeType};
        }
        .thumbnail.hidden {
            opacity: 0;
        }
        .thumbnail.visible {
            opacity: 1;
        }
    </style>
    `);

    function switchToNextThumbnail(container) {
        thumbnails = container.querySelectorAll('img');
        thumbnails[currentThumbnailIndex].className = 'thumbnail hidden';
        thumbnails[currentThumbnailIndex].style = "opacity: 0;";
        let nextIndex = (currentThumbnailIndex + 1) % thumbnails.length;
        while (nextIndex !== currentThumbnailIndex && thumbnails[nextIndex].classList.contains('error')) {
            nextIndex = (nextIndex + 1) % thumbnails.length;
        }
        currentThumbnailIndex = nextIndex;
        thumbnails[currentThumbnailIndex].className = 'thumbnail visible';
        thumbnails[currentThumbnailIndex].style = "opacity: 1;";
        setSizeContainer(container, width, height, thumbnails[currentThumbnailIndex]);
    }


    function connectWebSocket() {
        socket = new WebSocket(`ws://localhost:${port}/startport`);

        socket.onopen = () => {
            console.log("✅ WebSocket подключён");

            // Пример регистрации страницы
            socket.send(JSON.stringify({ type: "register", screen: "thumbnails" }));
        };

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.type === "update" && data.screen === "thumbnails") {
                if (lastGame !==data.payload.game){
                    lastThumbnails = data.payload.paths || [];
                    currentThumbnailIndex = 0;
                    update(lastThumbnails,currentThumbnailIndex, data.payload.width, data.payload.height);
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
    function update(lastThumbnails,currentThumbnailIndex, width, height){
        container.innerHTML = '';
        if (!lastThumbnails.length) {
            console.warn('No thumbnails found, displaying placeholder.');
            const img = document.createElement('img');
            img.src = '/theme/default/noimage.png?cache=' + Date.now();
            img.className = 'thumbnail visible';
            img.alt = 'No Thumbnail';
            container.appendChild(img);
            setSizeContainer(container, width, height, img);
        } else {
            const thumbnails = [];
            lastThumbnails.forEach((path, index) => {
                const img = document.createElement('img');
                img.src = path + '?cache=' + Date.now();
                img.className = 'thumbnail ' + (index === 0 ? 'visible' : 'hidden');
                img.alt = 'Thumbnail';
                img.onerror = () => {
                    img.classList.add('error');
                    console.warn(`Failed to load image: ${path}`);
                    if (index === 0) switchToNextThumbnail(thumbnails);
                };
                container.appendChild(img);
                thumbnails.push(img);
            });
        }
        startThumbnailInterval(container);
    }

};