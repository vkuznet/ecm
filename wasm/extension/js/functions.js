function Lock() {
    var rec = document.getElementById('records');
    rec.setAttribute('class', 'hide');
    rec.innerHTML = "";
    var search = document.getElementById("search");
    search.setAttribute('class', 'hide');
    var gpm = document.getElementById("gpm");
    gpm.setAttribute('class', 'show');
}
function Unlock() {
    var gpm = document.getElementById("gpm");
    gpm.setAttribute('class', 'hide');
    var search = document.getElementById("search");
    search.setAttribute('class', 'show');
    var rec = document.getElementById("records");
    rec.setAttribute('class', 'show');
}

async function asyncRecords() {
    Unlock();
    try {
        var x = document.getElementById("password");
        var password = x.value;
        x.value = "";
        const response = await records(password);
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
            var tags = document.createElement('div');
            tags.innerHTML = 'Tags: ' + rec.Tags;
            var login = document.createElement('div');
            login.innerHTML = 'Login: ' + rec.Login;
            var pass = document.createElement('div');
            var length = rec.Password.length;
            var hide = '*'.repeat(length);
            pass.innerHTML = 'Password: ' + hide;

            // add buttons
            var buttons = document.createElement('div');
            buttons.setAttribute('class', 'button-right');

            // add show button
            var button = document.createElement('button');
            button.type = "button";
            button.setAttribute('id', 'autofill');
            button.setAttribute('class', 'label success');
            button.innerHTML = "Show";
            button.addEventListener('click', showRecord, false);
            buttons.appendChild(button);

            // add autofill button
            var button = document.createElement('button');
            button.type = "button";
            button.setAttribute('id', 'autofill');
            button.setAttribute('class', 'label status');
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
            li.append(tags);
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
async function asyncDecode(input) {
    try {
        var x = document.getElementById("password");
        var password = x.value;
        //var password = prompt("Enter the password")
        const response = await decode(input, password);
        const data = await response.json();
        // console.log(data);
        output.value = "login: " + data.Map.Login + "password:" + data.Map.Password;
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
function showRecord(evt) {
}

// for information how to pass argument values see
// https://stackoverflow.com/questions/256754/how-to-pass-arguments-to-addeventlistener-listener-function
function fillFormInTab(evt) {
    var login = evt.currentTarget.Login;
    var password = evt.currentTarget.Password;
    chrome.tabs.query({active: true}, function(tabs) {
      var tab = tabs[0];
      tab_title = tab.title;
      chrome.tabs.executeScript(tab.id, {
          code : 'var inputs = document.getElementsByTagName("input");for(var i = 0; i < inputs.length; i++) {if(inputs[i].type.toLowerCase() == "text") {inputs[i].value = "'+login+'";}if(inputs[i].type.toLowerCase() == "password") {inputs[i].value = "'+password+'";}}'
      }, function(response) {});
    });
}

// helper function to return main web page inputs
function getPageInputs(evt) {
    var inputs = Array();
    chrome.tabs.query({active: true, currentWindow: true}, function(tabs) {
      var tab = tabs[0];
      tab_title = tab.title;
      chrome.tabs.executeScript(tab.id, {
          code : '(function(){return document.getElementsByTagName("input");})();'
      }, function(result) {
          for (i = 0; i < result[0].length; i++) inputs[i]=result[0][i];
      });
    });
    return inputs;
}

// helper function to find input type for given key
function inputType(key, inputs) {
  for(var i = 0; i < inputs.length; i++) {
      console.log("inputType", inputs[i].type, inputs[i]);
    if(inputs[i].name == key) {
      if(inputs[i].type.toLowerCase() == 'text') {
        return 'login';
      }
      if(inputs[i].type.toLowerCase() == 'password') {
        return 'password';
      }
    }
  }
  return 'na'
}

// fetch form POST request from a web page
// https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/API/webRequest/onBeforeRequest
// https://developer.chrome.com/docs/extensions/reference/webRequest/
// https://spin.atomicobject.com/2017/08/18/chrome-extension-form-data/
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
			/*
            if(value.includes("foo")) {
              cancel = true;
            }
            */
          });
        });
      }
    }
  },
  {urls: ["<all_urls>"]},
  ["requestBody"]
);
