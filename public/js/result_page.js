var filter = ["", "", ""];

function load_result(){
    let url_get_res = new URL(window.location.href);
    url_get_res.search = new URLSearchParams({result:"get"})
    fetch(url_get_res, {
        method: 'POST',
        headers: {
            'Content-type': 'application/json',
        },
        body: JSON.stringify(
            {
                source_filter: filter[0],
                key_filter: filter[1],
                content_filter: filter[2]
            }
        )
    })
    .then( res => res.json())
    .then( data => {
        console.log(data);
        result_list.innerHTML = data;
    })
}
window.onload = () =>{
    let result_list = document.getElementById('result_list');
    load_result();
    let butts = document.getElementsByClassName('img_butt');
    let input_fields = document.getElementsByTagName('input');
    butts[0].onclick = () => {
        filter[0] = input_fields[0].value;
        filter[1] = input_fields[1].value;
        filter[2] = input_fields[2].value;
        load_result();

    }
    butts[1].onclick = () => {
        load_result();
    }
}