import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("SignIn");
    }

    getHtmlFileName() {
        return 'signin.html';
    }
}
