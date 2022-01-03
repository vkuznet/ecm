function Lock() {
    var rec = document.getElementById('records');
    rec.setAttribute('class', 'hide');
    rec.innerHTML = "";
    var config = document.getElementById("config");
    config.setAttribute('class', 'hide');

    var password = document.getElementById("password");
    password.setAttribute('class', 'show-inline');
    var search = document.getElementById("search");
    search.setAttribute('class', 'hide');
    var lock = document.getElementById("lock");
    lock.setAttribute('class', 'is-warning hide');
    var unlock = document.getElementById("unlock");
    unlock.setAttribute('class', 'is-focus show');
}
function Unlock() {
    var config = document.getElementById("config");
    config.setAttribute('class', 'hide');
    var rec = document.getElementById("records");
    rec.setAttribute('class', 'show');

    var password = document.getElementById("password");
    password.setAttribute('class', 'hide');
    var search = document.getElementById("search");
    search.setAttribute('class', 'show-inline');
    var lock = document.getElementById("lock");
    lock.setAttribute('class', 'is-warning show');
    var unlock = document.getElementById("unlock");
    unlock.setAttribute('class', 'is-focus hide');
}
function Config() {
    var config = document.getElementById("config");
    config.setAttribute('class', 'show');
    var rec = document.getElementById("records");
    rec.setAttribute('class', 'hide');
}
function Exit() {
    var config = document.getElementById("config");
    config.setAttribute('class', 'hide');
    var rec = document.getElementById("records");
    rec.setAttribute('class', 'show');
}
// helper function to invoke asyncRecords on key enter in password field
function clickPress(event) {
    if (event.keyCode == 13) {
        asyncRecords();
    }
}
// helper function to show password of the record
function showPassword(evt) {
    var bid = evt.currentTarget.ButtonID;
    var pid = evt.currentTarget.PassID;
    var password = evt.currentTarget.Password;
    var doc = document.getElementById(pid);
    var button = document.getElementById(bid);
    if(button.innerHTML == "Show password") {
        doc.setAttribute('class', 'show-inline');
        doc.innerHTML = "Password: " + password;
        button.innerHTML = "Hide password";
    } else {
        doc.setAttribute('class', 'show-inline');
        doc.innerHTML = "";
        button.innerHTML = "Show password";
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
// helper function which we will use in go wasm
function fillFormHelper(evt) {
    var rid = evt.currentTarget.RecordID;
    var login = getLogin(rid);
    var password = getPassword(rid);
    console.log("login: ", login, " password: ", password)
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

