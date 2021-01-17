export function SignUpObject(data) {
    const credentials = data.credentials.split(' ');
    this.email = data.email;
    this.birthDate = data.date;
    this.name = credentials[0];
    this.surname = credentials[1];
}

export function SavePasswordObject(id, password) {
    this.userId = id;
    this.password = password;
}

export function SignInObject(email, password) {
    this.email = email;
    this.password = password;
}
