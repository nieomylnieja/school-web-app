import { Post } from "./http.js";
import { SignInObject } from "./models.js";

let form = document.getElementById("register");
let modal = document.getElementById("m-modal");

form.addEventListener("submit", (event) => {
    // TODO how to make it work
    event.preventDefault()
    if (validateForm(event)) {
        const data = Object.fromEntries(new FormData(event.target));
        Post('/auth/signin', new SignInObject(data.email, data.pass))
            .then(resp => {
                window.localStorage.setItem('token', resp.token);
                window.location.pathname = '/home'
            })
            .catch(err => {
                errorStore.push(err.toString())
                modal.style.display = "block";
                document.getElementById("modal-errors").innerHTML = errorStore.join("<br>");
                document.getElementById("modal-errors-h").innerText = "Wystąpił błąd podczas logowania!"
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
    const fs = [validateEmail];
    let invalidCtr = 0;
    fs.forEach(f => f(e) ? null : invalidCtr++);
    if (invalidCtr > 0) {
        modal.style.display = "block";
        document.getElementById("modal-errors").innerHTML = errorStore.join("<br>");
        return false;
    }
    return true;
}

function validateEmail(e) {
    const re = /^(([^<>()[\].,;:\s@"]+(\.[^<>()[\].,;:\s@"]+)*)|(".+"))@(([^<>()[\].,;:\s@"]+\.)+[^<>()[\].,;:\s@"]{2,})$/i;
    let input = form.elements["email"];
    const cond = (re.test(input.value.trim()));
    return validate(cond, e, input, "Niepoprawny adres email");
}

form.elements["email"].addEventListener("blur", validateEmail);

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
