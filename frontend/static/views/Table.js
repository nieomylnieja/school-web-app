import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Table");
    }

    requireLogin() {
        return true;
    }

    getHtmlFileName() {
        return 'table.html';
    }
}
