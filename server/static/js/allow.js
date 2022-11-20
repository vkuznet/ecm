// add listener functions to our input fields
var elem = document.getElementById("password");
if (elem) {
    elem.addEventListener("keypress", clickPress);
}
var elem = document.getElementById("search");
if (elem) {
    elem.addEventListener("keypress", clickPress);
}
var elem = document.getElementById("unlock");
if (elem) {
    elem.addEventListener("click", asyncRecords);
}
var elem = document.getElementById("lock");
if (elem) {
    elem.addEventListener("click", Lock);
}
var elem = document.getElementById("menu");
if (elem) {
    elem.addEventListener("click", Config);
}
var elem = document.getElementById("exit");
if (elem) {
    elem.addEventListener("click", Exit);
}

// fill out pageUrl input with actual url of the tab
tabUrl();
