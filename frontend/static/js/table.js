import { Get, Put } from "./http.js";

document.getElementById("t-submit-button")
    .addEventListener("click", async event => {
        await Put('/students', store)
            .catch(e => console.log(e))
        event.preventDefault()
    });

let store = []

renderTable().then();

async function renderTable() {
    await Get('/students')
        .then(resp => resp.json())
        .then(d => d.forEach(s => store.push(new Row(s.name, s.surname, s.age, s.id))))
        .catch(e => console.log(e))
    document.getElementById("t-button")
        .insertRow(0)
        .insertCell(0)
        .innerHTML =
        `
        <div id="counter" class="action">
          ${getCounterElement()}
        </div>
        `;
    for (let i = 0; i < store.length; i++) {
        insertRow(i);
    }
}

document.getElementById("t-home")
    .addEventListener("keyup", (event) => {
        let r = event.target;
        const cell = r.parentNode.cellIndex;
        const row = r.parentNode.parentNode.rowIndex - 1;
        if (cell === 0) {
            store[row].name = r.value;
        } else if (cell === 1) {
            store[row].surname = r.value;
        } else {
            const parsed = parseInt(r.value);
            if (!isNaN(parsed)) store[row].age = parsed;
        }
    });

function getCounterElement() {
    return `<b>N: ${store.length}</b><a class="add-action"></a>`
}

function insertRow(i) {
    let table = document.getElementById("t-home");
    let buttons = document.getElementById("t-button");
    renderActionButtons(buttons.insertRow(i + 1))
    let row = table.insertRow(i + 1);
    let pi = 0;
    for (const v of Object.values(store[i])) {
        renderCell(row, pi, v);
        if (pi++ === 2) return;
    }
}

function renderActionButtons(row) {
    row.insertCell(0).innerHTML =
        `<div class="action">
        <a class="delete-action""></a>
        <a class="down-action""></a>
        <a class="up-action""></a>
     </div>`
}

document.getElementById("t-button")
    .addEventListener("click", async event => {
        const i = event.target.parentNode.parentNode.parentNode.rowIndex - 1;
        switch (event.target.className) {
            case 'delete-action':
                remove(i);
                break
            case 'down-action':
                moveDown(i);
                break
            case 'up-action':
                moveUp(i);
                break
            case 'add-action':
                addRow();
                break
        }
    });

function moveUp(i) {
    if (i === 0) {
        return
    }
    [store[i - 1], store[i]] = [store[i], store[i - 1]];
    removeById(i);
    insertRow(i - 1);
}

function moveDown(i) {
    if (i === store.length - 1) {
        return
    }
    [store[i + 1], store[i]] = [store[i], store[i + 1]];
    removeById(i);
    insertRow(i + 1);
}

function remove(i) {
    store.splice(i, 1);
    removeById(i);
    document.getElementById("counter").innerHTML = getCounterElement();
}

function addRow() {
    const l = store.length;
    renderActionButtons(document.getElementById("t-button").insertRow(l + 1));
    let table = document.getElementById("t-home");
    let row = table.insertRow(l + 1);
    store[l] = new Row("", "", 0);
    let pi = 0;
    for (const v of Object.values(store[l])) {
        renderCell(row, pi, v);
        if (pi++ === 2) return;
    }
    document.getElementById("counter").innerHTML = getCounterElement();
}

function removeById(i) {
    document.getElementById("t-home").deleteRow(i + 1);
    document.getElementById("t-button").deleteRow(i + 1);
}

function renderCell(row, cell, value) {
    if (cell === 2) {
        row.insertCell(cell).innerHTML = `<input type="number" value="${value}">`;
    } else {
        row.insertCell(cell).innerHTML = `<input type="text" value="${value}">`;
    }
}

function Row(name, surname, age, id) {
    const parsed = parseInt(age, 10);
    if (isNaN(parsed)) age = 0;

    this.name = name;
    this.surname = surname;
    this.age = age;
    this.id = id;
}

