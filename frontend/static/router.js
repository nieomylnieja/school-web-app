import Home from './views/Home.js'
import SignUp from './views/SignUp.js';
import SignIn from './views/SignIn.js';
import Table from './views/Table.js';
import './js/jquery-3.5.1.min.js';
import { Get } from "./js/http.js";

const pathToRegex = path => new RegExp('^' + path.replace(/\//g, '\\/').replace(/:\w+/g, '(.+)') + '$');

const getParams = match => {
    const values = match.result.slice(1);
    const keys = Array.from(match.route.path.matchAll(/:(\w+)/g)).map(result => result[1]);

    return Object.fromEntries(keys.map((key, i) => {
        return [key, values[i]];
    }));
};

const navigateTo = url => {
    history.pushState(null, null, url);
    router().then(_ => _);
};

const router = async () => {
    const routes = [
        { path: '/', view: Home },
        { path: '/signup', view: SignUp },
        { path: '/signin', view: SignIn },
        { path: '/table', view: Table },
    ];

    // Test each route for potential match
    const potentialMatches = routes.map(route => {
        return {
            route: route,
            result: location.pathname.match(pathToRegex(route.path))
        };
    });

    let match = potentialMatches.find(potentialMatch => potentialMatch.result !== null);

    if (!match) {
        match = {
            route: routes[0],
            result: [location.pathname]
        };
    }

    const view = new match.route.view(getParams(match));

    // use jQuery to easily insert the template alongside with running eval() on the scripts within
    if (view.requireLogin()) {
        const token = window.localStorage.getItem('token');
        let authorized = false;
        if (!token) {
            window.location.pathname = '/signin';
        } else {
            await Get('/auth/verify')
                .then(() => authorized = true)
                .catch(() => {
                    window.localStorage.removeItem('token');
                    window.location.pathname = '/signin';
                })
        }
        if (!authorized) return;
    }
    $('#app').load('static/views/' + view.getHtmlFileName());
};

window.addEventListener('popstate', router);

document.addEventListener('DOMContentLoaded', () => {
    document.body.addEventListener('click', e => {
        if (e.target.matches('[data-link]')) {
            e.preventDefault();
            navigateTo(e.target.href);
        }
    });

    router().then(_ => _);
});
