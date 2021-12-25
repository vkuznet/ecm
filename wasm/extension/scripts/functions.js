function Lock() {
    var lock = document.getElementById('records');
    lock.class = "hide";
    var records = document.getElementById('records');
    records.innerHTML = "";
}
async function asyncRecords() {
    try {
        var x = document.getElementById("password");
        var password = x.value;
        x.value = "";
        const response = await records(password);
        const data = await response.json();
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
            var link = document.createElement('div');
            var button = document.createElement('button');
            button.type = "button";
            button.setAttribute('id', 'autofill');
            button.setAttribute('class', 'button button-fill');
            button.innerHTML = "Autofill";
            button.addEventListener('click', fillFormInTab, false);
            button.Login=rec.Login;
            button.Password=rec.Password;
            link.appendChild(button);
            //var a = document.createElement('a');
            //a.setAttribute('href', 'javascript:fillForm()');
            //a.innerHTML = 'Autofill';
            //link.appendChild(a);
            //link.innerHTML = 'Autofill';
            //link.setAttribute('class', 'link');
            var site = document.createElement('div');
            site.innerHTML = 'URL: ' + rec.URL;
            var copy = document.createElement('div');
            li.append(name);
            li.append(site);
            li.append(login);
            li.append(pass);
            li.append(tags);
            li.append(link);
            //console.log(rec.Login);
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
