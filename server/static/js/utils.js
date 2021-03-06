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
    fileReader.readAsArrayBuffer(file);
    fileReader.onload = function() {
        let data = buf2hex(fileReader.result);
//        console.log(file.name);
//        console.log(data);
        // call wasm uploadFile function
        uploadFile(file.name, file.size, file.type, data);
    };
    fileReader.onerror = function() {
      alert(fileReader.error);
    };
}
// https://stackoverflow.com/questions/40031688/javascript-arraybuffer-to-hex
function buf2hex(buffer) { // buffer is an ArrayBuffer
  return [...new Uint8Array(buffer)]
      .map(x => x.toString(16).padStart(2, '0'))
      .join('');
}
