// ВНИМАНИЕ!
// НИ В КОЕМ СЛУЧАЕ НЕ МЕНЯЙТЕ КОД!
// ТУТ НЕЧЕГО МЕНЯТЬ!
// ЭТОТ КОД НУЖЕН ДЛЯ РАБОТЫ СИСТЕМЫ!
// ANTENTION!
// DO NOT CHANGE THE CODE IN ANY WAY!
// THERE'S NOTHING TO CHANGE!
// THIS CODE IS NEEDED FOR THE SYSTEM TO WORK!
function copyToClipboard(event) {
    const altText = event.target.alt; // Получаем текст из атрибута alt
    if (altText) {
        navigator.clipboard.writeText(altText)
            .then(() => {
                alert("Текст скопирован: " + altText);
            })
            .catch(err => {
                console.error("Ошибка копирования текста: ", err);
            });
    } else {
        alert("Атрибут alt пустой!");
    }
}