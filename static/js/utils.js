function HideTag(tag) {
    var id=document.getElementById(tag);
    if (id) {
        id.className="hide";
    }
}
function ShowTag(tag) {
    var id=document.getElementById(tag);
    if (id) {
        id.className="show";
    }
}
function FlipTag(tag) {
    var id=document.getElementById(tag);
    if (id) {
        if  (id.className == "show") {
            id.className="hide";
        } else {
            id.className="show";
        }
    }
}
function readFile(input) {
    let file = input.files[0];
    let fileReader = new FileReader();
    // TODO: replace readAsText to reading bytes
    fileReader.readAsText(file);
    fileReader.onload = function() {
        let data = fileReader.result
        console.log(file.name);
        console.log(data);
        // call wasm uploadFile function
        uploadFile(file.name, file.size, file.type, data);
    };
    fileReader.onerror = function() {
      alert(fileReader.error);
    };
}
