const Http = new XMLHttpRequest();
const url_test='http://localhost:60000/api/testt';
const url_role='http://localhost:60000/api/roles/';

function getRoles(){
    Http.open("GET", url_role);
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
