const router = new Navigo(null, true, '#');

router.on('form', {
    'form': () => {
        loadHTML('./templates/form.html', 'app')
    },
});
router.notFound(() => {
    $id('app').innerHTML = '<h3>Strona nie istnieje...</h3>';
});


router.resolve();

function $id(id) {
    return document.getElementById(id);
}

function loadHTML(url, id) {
    let req = new XMLHttpRequest();
    req.open('GET', url);
    req.send();
    req.onload = () => {
        $id(id).innerHTML = req.responseText;
    };
}
