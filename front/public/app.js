const Http = new XMLHttpRequest();
const url_test='http://localhost:3333/test';
const url_hello='http://localhost:3333/hello';

function getHello(){
    Http.open("GET", url_hello);
    Http.onreadystatechange = (e) => {
        console.log(Http.responseText)
    }
    Http.send();
}

function getTest(){
    Http.open("GET", url_test);
    Http.onreadystatechange = (e) => {
        console.log(Http.responseText)
    }
    Http.send();
}

function getRoles(){
    Http.open("GET", "http://localhost:60000/api/roles")
    Http.onreadystatechange = (e) => {
        console.log(Http.responseText)
    }
    Http.send()
}
