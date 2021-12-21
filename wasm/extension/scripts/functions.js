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
            //li.innerHTML=li.innerHTML + rec.Login + '<br/>' + rec.URL;
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
            link.innerHTML = 'Link';
            link.setAttribute('class', 'link');
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
