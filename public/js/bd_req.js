function get_dbs(cur_url){
    let get_conf = new URLSearchParams({
        bd:"get_conf"
    });
    cur_url.search = get_conf.toString();
    fetch(cur_url, {
        method: 'GET',
        headers: {
            'Content-type': 'application/json',
        }
    })
    .then((res) => res.json())
    .then((data) => {
        let tmp = document.getElementById('dbs_list');
        let keys = Object.keys(data);
        tmp.innerHTML = '';
        for (let i = 0; i < Object.keys(data).length; i++){
            let div = document.createElement('div');
            let del_butt = document.createElement('img');
            div.innerHTML += "<p>" + data[keys[i]].Name + " :</p>" +
            "<p> Адрес сервера СУБД: " + data[keys[i]].Addr +  ":" + data[keys[i]].Port + "</p>" + 
            "<p> Категория БД: " + data[keys[i]].DB + "</p>";
            div.dataset.id = keys[i];
            del_butt.src = "public/img/del_icon.png";
            del_butt.onclick = () => {
                div.style.display = 'none';
                let tmp_url = new URL(window.location.href);
                let del_conf = new URLSearchParams({
                    bd:"del_conf"
                });
                tmp_url.search = del_conf.toString();
                fetch(tmp_url, {
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
            }
            div.appendChild(del_butt);
            tmp.appendChild(div);
        }

    })
}
window.onload = () => {
    let url = new URL(window.location.href);
    let form = document.getElementById('db_config_form');
    let butt = document.getElementsByClassName('butt');
    let fields = document.getElementsByTagName('input');
    let slctr = document.getElementsByTagName('select');
    get_dbs(url);
    butt[0].onclick = () => {
        form.style.display = 'flex';
    }
    butt[1].onclick = () => {
        get_dbs(url);
    };
    butt[2].onclick = () => {
        let get_conf = new URLSearchParams({
            bd:"add_conf"
        });
        url.search = get_conf.toString();
        fetch(url, {
            method: 'POST',
            headers: {
                'Content-type': 'application/json',
            },
            body: JSON.stringify(
                {
                    Name    : fields[0].value,
                    Type    : "mongodb",
                    Addr    : fields[1].value,
                    Port    : fields[2].value,
                    Login   : fields[3].value,
                    Password: fields[4].value,
                    DB      : slctr[0].value
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
    }
    butt[3].onclick = () =>{
        form.style.display = 'none';
    }
}