import Navigo from 'navigo';
import './style.css';

const router = new Navigo(null, true, '#');

function getById(id) {
    return document.getElementById(id);
}

function loadHTML(url, id) {
    let req = new XMLHttpRequest();
    req.open('GET', url);
    req.send();
    req.onload = () => {
        getById(id).innerHTML = req.responseText;
    };
}

router.on('form', {
    'form': () => {
        loadHTML('./templates/form.html', 'app')
    },
});
router.notFound(() => {
    getById('app').innerHTML = '<h3>Strona nie istnieje...</h3>';
});

router.resolve();
