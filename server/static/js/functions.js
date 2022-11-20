async function asyncRecords() {
    // clean-up previously shown content
    var rec = document.getElementById('records');
    rec.innerHTML = "";
    // obtain server creds
    var pageUrl = "";
    var server = document.getElementById("server").value;
    var vault = document.getElementById("vault").value;
    var cipher = document.getElementById("cipher").value;
    var pattern = document.getElementById("search").value;
    var x = document.getElementById("password");
    var password = x.value;
    try {
        const response = await records(server, vault, cipher, password, pattern, pageUrl);
        const data = await response.json();
        for (const key of data) {
            var button = document.getElementById('autofill-'+key)
            if (button) {
                button.addEventListener('click', fillFormHelper, false);
            }
        }
    } catch (err) {
        console.log(err)
        console.error('Caught exception', err)
    }
}
