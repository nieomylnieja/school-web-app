let store = [
    new Row('Andrzej', 'Bargiel', 50),
    new Row('Jerzy', 'Kukuczka', 94),
    new Row('Leszek', 'Cichy', 80),
    new Row('Wanda', 'Rutkiewicz', 56),
    new Row('Maciej', 'Berbeka', 65),
    new Row('Adam', 'Bielecki', 38),
    new Row('Andrzej', 'Zawada', 41),
    new Row('Artur', 'Hajzer', 78),
    new Row('Tomek', 'Mackiewicz', 46),
]

renderTable();

function renderTable() {
    document.getElementById("t-button")
        .insertRow(0)
        .insertCell(0)
        .innerHTML =
        `
        <div id="counter" class="action">
          ${getCounterElement()}
        </div>
        `;
    for (let i = 1; i < store.length; i++) {
        insertRow(i);
    }
}

function getCounterElement() {
    return `<b>N: ${store.length}</b><a class="add-action" onclick="appendRow()"></a>`
}

function insertRow(i) {
    let table = document.getElementById("t-home");
    let buttons = document.getElementById("t-button");
    renderActionButtons(buttons.insertRow(i))
    let row = table.insertRow(i);
    let pi = 0;
    for (const v of Object.values(store[i])) {
        renderCell(row, pi, v);
        pi++;
    }
}

function moveUp(r) {
    const i = r.parentNode.parentNode.parentNode.rowIndex;
    if (i === 1) {
        return
    }
    [store[i - 1], store[i]] = [store[i], store[i - 1]];
    removeById(i);
    insertRow(i - 1);
}

function moveDown(r) {
    const i = r.parentNode.parentNode.parentNode.rowIndex;
    if (i === store.length - 1) {
        return
    }
    [store[i + 1], store[i]] = [store[i], store[i + 1]];
    removeById(i);
    insertRow(i + 1);
}

function remove(r) {
    const i = r.parentNode.parentNode.parentNode.rowIndex;
    store.splice(i, 1);
    removeById(i);
    document.getElementById("counter").innerHTML = getCounterElement();
}

function removeById(i) {
    document.getElementById("t-home").deleteRow(i);
    document.getElementById("t-button").deleteRow(i);
}

function renderActionButtons(row) {
    row.insertCell(0).innerHTML =
        `<div class="action">
        <a class="delete-action" onclick="remove(this)"></a>
        <a class="down-action" onclick="moveDown(this)"></a>
        <a class="up-action" onclick="moveUp(this)"></a>
     </div>`
}

function saveInput(r) {
    const cell = r.parentNode.cellIndex;
    const row = r.parentNode.parentNode.rowIndex;
    if (cell === 0) {
        store[row].name = r.value;
    } else if (cell === 1) {
        store[row].surname = r.value;
    } else {
        store[row].age = r.value;
    }
}

function appendRow() {
    const l = store.length;
    renderActionButtons(document.getElementById("t-button").insertRow(l));
    let table = document.getElementById("t-home");
    let row = table.insertRow(l);
    store[l] = new Row("", "", "");
    let pi = 0;
    for (const v of Object.values(store[l])) {
        renderCell(row, pi, v);
        pi++;
    }
    document.getElementById("counter").innerHTML = getCounterElement();
}

function renderCell(row, cell, value) {
    row.insertCell(cell).innerHTML = `<input type="text" value="${value}" onblur="saveInput(this)">`;
}

function Row(name, surname, age) {
    this.name = name;
    this.surname = surname;
    this.age = age;
}
