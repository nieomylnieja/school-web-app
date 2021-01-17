export default class {
    constructor(params) {
        this.params = params;
    }

    setTitle(title) {
        document.title = title;
    }

    requireLogin() {
        return false;
    }

    async getHtmlFileName() {
        return "";
    }
}
