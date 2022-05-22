function loadDoc() {
  const xhttp = new XMLHttpRequest();
  xhttp.onload = function() {
    document.getElementById("demo").innerHTML =
    this.responseText;
  }
  xhttp.open("GET", "/records");
  xhttp.send();
}

/*
 * see https://www.w3schools.com/js/js_ajax_examples.asp
 *
<div id="demo">

<h2>The XMLHttpRequest Object</h2>

<button type="button" onclick="loadDoc('ajax_info.txt', myFunction)">Change Content
</button>
</div>

<script>
function loadDoc(url, xFunction) {
  const xhttp=new XMLHttpRequest();
  xhttp.onload = function() {xFunction(this);}
  xhttp.open("GET", url);
  xhttp.send();
}

function myFunction(xhttp) {
  document.getElementById("demo").innerHTML =  xhttp.responseText;
}
</script>
 */
