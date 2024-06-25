let ref_func;
let run = false;
function startInterval() {
    run = true;
    ref_func = setInterval(function() {
        if (run) {
            data_update();
        }
    }, 2000);
}

function pauseInterval() {
    run = false;
    console.log('Функция на паузе...');
}

function resumeInterval() {
    run = true;
    console.log('Функция возобновлена...');
}

function stopInterval() {
    clearInterval(ref_func);
    console.log('Функция остановлена навсегда...');
}

startInterval();
function data_update(){
    let list = document.getElementById('spyders_list');
    let url_ref = new URL(window.location.href);
    url_ref.search = new URLSearchParams({spdr:"ref"})
    fetch(url_ref, {
        method: 'POST',
        headers: {
            'Content-type': 'application/json',
        },
    })
    .then(res => res.json())
    .then(data => {
        console.log(data);
        let keys = Object.keys(data);
            list.innerHTML = "";
        for (let i = 0; i < keys.length; i++){
            let div = document.createElement('div');
            let p1 = document.createElement('p');
            let p2 = document.createElement('p');
            let img = document.createElement('img');
            p1.textContent = keys[i];
            p2.textContent = data[keys[i]].Addr + ";" + data[keys[i]].Config.fetcher["xmlrpc-port"] + ";" + data[keys[i]].Config.puppeteer.port + ";" + data[keys[i]].Config.scheduler["xmlrpc-port"];
            img.src = "public/img/del_icon.png";
            img.classList.add("del_spdr");
            img.onclick = () => {
                div.style.display = "none";
                let url_del = new URL(window.location.href);
                url_del.search = new URLSearchParams({spdr:"del"})
                stopInterval();
                fetch(url_del, {
                    method: 'POST',
                    headers: {
                        'Content-type': 'application/json',
                    },
                    body:JSON.stringify({
                        "to_del":keys[i]
                    })
                })
                .then( res => {
                    if (res.ok) {
                        console.log("Даннные успешно удалены!");
                        return res.text()
                    }
                })
                data_update();
                startInterval();
            };
            div.appendChild(p1);
            div.appendChild(p2);
            div.appendChild(img);
            list.appendChild(div);
        }
    })
    .catch(error => console.error('Ошибка при получении данных:', error));
}
window.onload = () => {
    data_update();
    let url_add = new URL(window.location.href);
    let url_get = new URL(window.location.href);
    let fields = document.getElementsByClassName('spdr_stngs');
    url_get.search = new URLSearchParams({spdr:"gt_opt"});
    fetch(url_get, {
        method: 'POST',
        headers: {
            'Content-type': 'application/json',
        },
    })
    .then(res => res.json())
    .then(data => {
        fields[1].innerHTML = '';
        fields[2].innerHTML = '';
        fields[3].innerHTML = '';
        for (let i = 0; i < Object.keys(data['taskdb']).length; i++){
            let str = JSON.stringify(data['taskdb'][String(i)].Name) + " " + JSON.stringify(data['taskdb'][String(i)].Type) + " " + JSON.stringify(data['taskdb'][String(i)].Addr);
            fields[1].options.add(new Option(str, JSON.stringify(data['taskdb'][String(i)])));
        }
        for (let i = 0; i < Object.keys(data['projectdb']).length; i++){
            let str = JSON.stringify(data['projectdb'][String(i)].Name) + " " + JSON.stringify(data['projectdb'][String(i)].Type) + " " + JSON.stringify(data['projectdb'][String(i)].Addr);
             fields[2].options.add(new Option(str, JSON.stringify(data['projectdb'][String(i)])));
         }
        for (let i = 0; i < Object.keys(data['resultdb']).length; i++){
            let str = JSON.stringify(data['resultdb'][String(i)].Name) + " " + JSON.stringify(data['resultdb'][String(i)].Type) + " " + JSON.stringify(data['resultdb'][String(i)].Addr);
            fields[3].options.add(new Option(str, JSON.stringify(data['resultdb'][String(i)])));
        }
    })
    let butt = document.getElementsByClassName('butt');
    let c = document.getElementById('add_spdr_menu');
    console.log(butt);
    url_add.search = new URLSearchParams({spdr:"add"});
    butt[0].onclick = () => { c.style.display = 'flex';}
    butt[1].onclick = () => {
        let taskdb = JSON.parse(fields[1].value);
        let resultdb = JSON.parse(fields[3].value);
        let projectdb = JSON.parse(fields[2].value);
        fetch(url_add, {
            method: 'POST',
            headers: {
                'Content-type': 'application/json',
            },
            body: JSON.stringify(
                {
                    Name          : fields[0].value,
                    TaskDB        : {
                        Name: taskdb.Name,
                        Type: taskdb.Type,
                        Addr: taskdb.Addr            
                    },
                    ResultDB      : {
                        Name: resultdb.Name,
                        Type: resultdb.Type,
                        Addr: resultdb.Addr            
                    },
                    ProjectDB     : {
                        Name: projectdb.Name,
                        Type: projectdb.Type,
                        Addr: projectdb.Addr            
                    },
                    Addr          : fields[4].value,
                }
            )
        })
        .then( res => {
            if (res.ok) {
                console.log("Даннные успешно отправлены!");
                return res.text()
            }
        })
        .then( data => {
            console.log(data);
        })
        data_update();
    }
    butt[2].onclick = () => {c.style.display = 'none';};
}