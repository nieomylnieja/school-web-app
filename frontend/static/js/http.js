const methodPost = 'POST';
const methodGet = 'GET';
const backendBaseUrl = 'http://localhost:' + 9000;

async function doRequest(pattern, method, body) {
    let params = {
        headers: {
            'Access-Control-Request-Method': method,
            'Authorization': `Bearer ${window.localStorage.getItem('token')}`
        },
        method: method,
    }
    if (body) {
        params.body = JSON.stringify(body);
    }
    const url = backendBaseUrl + pattern
    return await fetch(url, params)
        .then(async data => {
            const code = data.status;
            if (code < 200 || code >= 300) {
                throw new Error(await data.text());
            }
            if (method === methodPost) return await data.json();
            return data;
        })
}

export function Get(pattern) {
    return doRequest(pattern, methodGet, null)
}

export function Post(pattern, body) {
    return doRequest(pattern, methodPost, body)
}