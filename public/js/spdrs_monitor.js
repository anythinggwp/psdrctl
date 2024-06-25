var ws;

function openSocket(obj,callback) {
    ws = new WebSocket("ws://" + window.location.hostname + ":8000/ws");
    
    ws.onopen = function(event) {
        console.log("Connection opened");
    };

    ws.onmessage = function(event) {
        const obj_data = JSON.parse(event.data)
        // console.log(event);
        // console.log(Object.keys(obj_data).length);
        callback(obj,obj_data);
    };

    ws.onclose = function(event) {
        console.log("Connection closed");
    };

    ws.onerror = function(error) {
        console.log("Error: " + error);
    };
}

function closeSocket() {
    ws.close();
}
function output_status(list, event){
    let info = document.getElementsByClassName('info');
    for (const conteiner of info) {
        if (conteiner.dataset.check == event.Addr && conteiner.dataset.status == event.Status){
            return
        }
        else if (conteiner.dataset.check == event.Addr && conteiner.dataset.status != event.Status){
            list.innerHTML = "";
        }
    }
    let base = document.createElement('div');
    let spdr = document.createElement('div');
    let tsk_db = document.createElement('div');
    let res_db = document.createElement('div');
    let proj_db = document.createElement('div');
    spdr.dataset.check = event.Addr;
    spdr.dataset.status = event.Status;
    spdr.className = "info";
    tsk_db.dataset.check = event.TSK.Addr;
    tsk_db.dataset.status = event.TSK.Status;
    tsk_db.className = "info";
    res_db.dataset.check = event.RES.Addr;
    res_db.dataset.status = event.RES.Status;
    res_db.className = "info";
    proj_db.dataset.check = event.PROJ.Addr;
    proj_db.dataset.status = event.PROJ.Status;
    proj_db.className = "info";
    let spdr_addr = document.createElement('a');
    let spdr_status = document.createElement('p');
    let tsk_db_type = document.createElement('p');
    let tsk_db_addr = document.createElement('p');
    let tsk_db_status = document.createElement('p');
    let res_db_type = document.createElement('p');
    let res_db_addr = document.createElement('p');
    let res_db_status = document.createElement('p');
    let proj_db_type = document.createElement('p');
    let proj_db_addr = document.createElement('p');
    let proj_db_status = document.createElement('p');
    spdr_addr.innerHTML = event.Addr;
    spdr_addr.href = "http://" + event.Addr;
    spdr_addr.target = "_blank";
    spdr_status.innerHTML = event.Status;
    tsk_db_type.innerHTML = event.TSK.Type;
    tsk_db_addr.innerHTML = event.TSK.Addr;
    tsk_db_status.innerHTML = event.TSK.Status;
    res_db_type.innerHTML = event.RES.Type;
    res_db_addr.innerHTML = event.RES.Addr;
    res_db_status.innerHTML = event.RES.Status;
    proj_db_type.innerHTML = event.PROJ.Type;
    proj_db_addr.innerHTML = event.PROJ.Addr;
    proj_db_status.innerHTML = event.PROJ.Status;
    spdr.appendChild(spdr_addr);
    spdr.appendChild(spdr_status);
    tsk_db.appendChild(tsk_db_type);
    tsk_db.appendChild(tsk_db_addr);
    tsk_db.appendChild(tsk_db_status);
    res_db.appendChild(res_db_type);
    res_db.appendChild(res_db_addr);
    res_db.appendChild(res_db_status);
    proj_db.appendChild(proj_db_type);
    proj_db.appendChild(proj_db_addr);
    proj_db.appendChild(proj_db_status);
    base.appendChild(spdr);
    base.appendChild(tsk_db);
    base.appendChild(res_db);
    base.appendChild(proj_db);
    list.appendChild(base);


}
window.onload = () =>{
    let spdr_list = document.getElementById('projects_list');
    openSocket(spdr_list, output_status);
}
