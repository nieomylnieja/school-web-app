import { Post } from "./http.js";
import { SignUpObject, SavePasswordObject } from "./models.js";

let form = document.getElementById("register");
let modal = document.getElementById("m-modal");

form.addEventListener("submit", (event) => {
    event.preventDefault()
    if (validateForm(event)) {
        const data = Object.fromEntries(new FormData(event.target));
        Post('/users', new SignUpObject(data))
            .then(resp => Post('/auth/signup', new SavePasswordObject(resp.id, data.pass)))
            .then(resp => {
                window.location.pathname = '/signin'
            })
            .catch(err => {
                errorStore.push(err.toString())
                modal.style.display = "block";
                document.getElementById("modal-errors").innerHTML = errorStore.join("<br>");
                document.getElementById("modal-errors-h").innerText = "Wystąpił błąd podczas rejestracji!"
            });
    }
});

// error modal display and discard
window.onclick = event => ((event.target === modal) ? clearErrorDisplay() : null);
document.getElementsByClassName("close")[0].onclick = () => clearErrorDisplay();

function clearErrorDisplay() {
    modal.style.display = "none";
    errorStore = [];
}

let errorStore = [];

function validateForm(e) {
    const fs = [validateCredentials, validateEmail, validateDate, validatePassword, validatePasswordRe, validateDescription];
    let invalidCtr = 0;
    fs.forEach(f => f(e) ? null : invalidCtr++);
    if (invalidCtr > 0) {
        modal.style.display = "block";
        document.getElementById("modal-errors").innerHTML = errorStore.join("<br>");
        return false;
    }
    return true;
}

function validateCredentials(e) {
    let input = form.elements["credentials"];
    const cond = (input.value.split(" ").length >= 2);
    return validate(cond, e, input, "Imię oraz nazwisko muszą zostać podane");
}

form.elements["credentials"].addEventListener("blur", validateCredentials);

function validateEmail(e) {
    const re = /^(([^<>()[\].,;:\s@"]+(\.[^<>()[\].,;:\s@"]+)*)|(".+"))@(([^<>()[\].,;:\s@"]+\.)+[^<>()[\].,;:\s@"]{2,})$/i;
    let input = form.elements["email"];
    const cond = (re.test(input.value.trim()));
    return validate(cond, e, input, "Niepoprawny adres email");
}

form.elements["email"].addEventListener("blur", validateEmail);

function validateDate(e) {
    let input = form.elements["date"];
    const cond = input.value;
    return validate(cond, e, input, "Podaj datę urodzenia");
}

form.elements["date"].addEventListener("blur", validateDate)

function validatePassword(e) {
    let input = form.elements["pass"];
    let v = input.value;
    const cond = (v.length >= 8 && /\d/.test(v));
    return validate(cond, e, input, "Hasło musi składać się z 8 znaków oraz zawierać min. 1 cyfrę");
}

form.elements["pass"].addEventListener("blur", validatePassword)

function validatePasswordRe(e) {
    let pass = form.elements["pass"].value;
    let input = form.elements["pass-re"];
    let v = input.value;
    const cond = (v && pass === v);
    return validate(cond, e, input, "Wprowadzono niepasujące do siebie hasła")
}

form.elements["pass-re"].addEventListener("blur", validatePasswordRe)

function validateDescription(e) {
    let input = form.elements["description"];
    const cond = (input.value.length <= 20);
    return validate(cond, e, input, "Opis może zawierać co najwyżej 20 znaków")
}

form.elements["description"].addEventListener("blur", validateDescription)

function validate(cond, e, input, msg) {
    return cond ? success(e, input) : error(e, input, msg)
}

function success(e, input) {
    if (e.type === "blur") input.setAttribute("style", "border: 5px solid #a3be8c;");
    return true;
}

function error(e, input, msg) {
    if (e.type === "blur") input.setAttribute("style", "border: 5px solid #bf616a;");
    if (e.type === "submit") errorStore.push(msg);
    return false;
}
