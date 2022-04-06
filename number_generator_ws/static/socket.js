let socket = new WebSocket("ws://127.0.0.1:4000/socket");
console.log("Attempting Connection...");

socket.onopen = () => {
    console.log("Successfully Connected");
};

socket.onmessage = function (event) {
    console.log(event.data)
    let el = document.getElementById("nums");
    if (event.data == "0") {
        el.innerText = "Invalid input";
    } else {
        el.append(event.data + " ");
    }
};

socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
};

socket.onerror = error => {
    console.log("Socket Error: ", error);
};

document.forms.form.onsubmit = function () {
    let el = document.getElementById("nums");
    el.innerText = "";
    let range_num = this.range_num.value;
    let limit = this.limit.value;
    let go = this.go.value;

    socket.send([limit, range_num, go]);
    return false;
};

