// flag to check if password was read
let passwordRead = false;

async function asyncRecords() {
    // clean-up previously shown content
    var rec = document.getElementById('records');
    rec.innerHTML = "";
    // unlock records part
    Unlock();
    // obtain server creds
    var server = document.getElementById("server").value;
    var vault = document.getElementById("vault").value;
    var cipher = document.getElementById("cipher").value;
    var pattern = document.getElementById("search").value;
    var x = document.getElementById("password");
    var password = x.value;
    x.value = "";
    if(password == "" && passwordSet == false) {
        Lock();
        var doc = document.getElementById('records');
        doc.setAttribute('class', 'alert is-error');
        doc.innerHTML = 'Invalid password';
        return
    }
    passwordSet=true;
    try {
        const response = await records(server, vault, cipher, password, pattern);
        const data = await response.json();
        for (const key of data) {
            console.log("array key", key)
            var button = document.getElementById('autofill-'+key)
            button.addEventListener('click', fillFormHelper, false);
        }
    } catch (err) {
        console.error('Caught exception', err)
    }
}

// fetch our records from remote server
async function asyncRecords_Original() {
    // clean-up previously shown content
    var rec = document.getElementById('records');
    rec.innerHTML = "";
    // unlock records part
    Unlock();
    // obtain server creds
    var server = document.getElementById("server").value;
    var vault = document.getElementById("vault").value;
    var cipher = document.getElementById("cipher").value;
    var pattern = document.getElementById("search").value;
    var x = document.getElementById("password");
    var password = x.value;
    x.value = "";
    if(password=="") {
        var doc = document.getElementById('records');
        doc.setAttribute('class', 'alert is-error');
        doc.innerHTML = 'Invalid password';
        document.getElementById('records').appendChild(doc);
        return
    }
    try {
        const response = await records(server, vault, cipher, password, pattern);
        const data = await response.json();
        // add records list
        var ul = document.createElement('ul');
        ul.setAttribute('class','records');
        document.getElementById('records').appendChild(ul);
        for (let key in data) {
            let rec = data[key];
            var li = document.createElement('li');
            li.setAttribute('class','item');
            ul.appendChild(li);
            var name = document.createElement('div');
            name.innerHTML = 'Name: ' + rec.Name;
            //var tags = document.createElement('div');
            //tags.innerHTML = 'Tags: ' + rec.Tags;
            var login = document.createElement('div');
            login.innerHTML = 'Login: ' + rec.Login;
            var pass = document.createElement('div');
            var pid = 'pid-'+key;
            pass.setAttribute('id', pid);
            pass.setAttribute('class', 'hide');
            //var length = rec.Password.length;
            //var hide = '*'.repeat(length);
            //pass.innerHTML = 'Password: ' + hide;

            // add buttons
            var buttons = document.createElement('div');
            buttons.setAttribute('class', 'button-right');

            // add show button
            var button = document.createElement('button');
            var bid = 'bid-'+key;
            button.type = "button";
            button.setAttribute('id', bid);
            button.setAttribute('class', 'label');
            button.innerHTML = "Show password";
            button.addEventListener('click', showPassword, false);
            button.ButtonID = bid;
            button.PassID=pid;
            button.Password=rec.Password;
            buttons.appendChild(button);

            // add autofill button
            var button = document.createElement('button');
            button.type = "button";
            button.setAttribute('id', 'autofill-'+key);
            button.setAttribute('class', 'label autofill is-bold');
            button.innerHTML = "Autofill";
            button.addEventListener('click', fillFormInTab, false);
            button.Login=rec.Login;
            button.Password=rec.Password;
            buttons.appendChild(button);

            var site = document.createElement('div');
            site.innerHTML = 'URL: ' + rec.URL;
            var copy = document.createElement('div');
            li.append(name);
            li.append(site);
            li.append(login);
            li.append(pass);
            //li.append(tags);
            li.append(buttons);
        }
        // allow autofill button to execute click function
        //document.getElementById("autofill").addEventListener("click", fillForm);
        // console.log(data);
        //output.value = JSON.stringify(data)
    } catch (err) {
        console.error('Caught exception', err)
    }
}
// functino to fill form on given web page
// https://stackoverflow.com/questions/5897122/accessing-elements-by-type-in-javascript
function fillForm() {
    var inputs = document.getElementsByTagName('input');
    for(var i = 0; i < inputs.length; i++) {
        if(inputs[i].type.toLowerCase() == 'text') {
            inputs[i].value = 'login';
        }
        if(inputs[i].type.toLowerCase() == 'password') {
            inputs[i].value = 'password';
        }
    }
}
// fetch form POST request from a web page
// https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/API/webRequest/onBeforeRequest
// https://developer.chrome.com/docs/extensions/reference/webRequest/
// https://spin.atomicobject.com/2017/08/18/chrome-extension-form-data/
/*
chrome.webRequest.onBeforeRequest.addListener(
  function(details) {
    if(details.method == "POST") {
      var inputs = getPageInputs();
      console.log('inputs', inputs)
      console.log('details', details)
      console.log(details.url);
      console.log(details.requestBody.formData);
      let formData = details.requestBody.formData;
      if(formData) {
        Object.keys(formData).forEach(key => {
          formData[key].forEach(value => {
			console.log(key, value, inputType(key, inputs));
          });
        });
      }
    }
  },
  {urls: ["<all_urls>"]},
  ["requestBody"]
);
*/
