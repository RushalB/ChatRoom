var display_username = ""
var socket = ""
let AUTH_STATUS = false

function addMessage(displayname, message, type) {
    var ul = document.getElementById("messages");
    var li = document.createElement("li");
    h3 = document.createElement("h3");
    p = document.createElement("p");
    h3.appendChild(document.createTextNode(displayname));
    p.appendChild(document.createTextNode(message));
    li.appendChild(h3); li.appendChild(p)
    if (type === "right") { li.setAttribute("class", "right") };
    if (type === "left") { li.setAttribute("class", "") };
    ul.appendChild(li);
    (ul.parentElement).scrollTop = (ul.parentElement).scrollHeight;


}


function addError(message,color) {
    code = 0
    var ul = document.getElementById("messages");
    var li = document.createElement("li");
    li.appendChild(document.createTextNode(message));
    li.setAttribute("class", "error")
    li.style.color = color;

    ul.appendChild(li);
    (ul.parentElement).scrollTop = (ul.parentElement).scrollHeight;


}




function sendMessage() {
    var input = document.getElementById("input");
    message = input.value
    if (message != "") {
        input.value = ""
        sendMsg(message)
    }
    else {
        alert("Enter Message to send")
    }
}




async function authenticate() {


    var password = document.getElementById("inputPassword")
    var username = document.getElementById("username")
    console.log(password.value, username.value)

    try {
        let response = await fetch("http://localhost:8000/login", {
            method: "POST",
            body:
                JSON.stringify({
                    "username": username.value,
                    "password": password.value
                }),
        });
   



    let user = (await response.json());
    AUTH_STATUS = user.Status;
    display_username = user.Username


    console.log(AUTH_STATUS);
    if (AUTH_STATUS) {
        var main = document.getElementById("main")
        var title = document.getElementById("title")
        var loginDiv = document.getElementById("login")
        main.style.display = 'block'
        loginDiv.remove()
        title.innerHTML += "- " + display_username

        socket = new WebSocket("ws:///localhost:8000/ws");
        connect()




    }
    else {
        password.value = ""
        username.value = ""
        alert("Incorrect Password")
    }

} catch (error) {
    alert("Server Down")
    return
}

}










var CONN_STATUS = 0

function sendMsg(msg) {
    console.log("sending msg: ", msg);
    if (CONN_STATUS == 1) {
        socket.send(JSON.stringify({"type":0,"body":msg}));
    }
    else {
        addError("Unable to send message","red")
    }
};

let connect = () => {


    socket.onopen = () => {

        console.log("Successfully Connected");
        addError("Successfully Connected to chat Room","grey")
        CONN_STATUS = 1
        //Confirmation to server 
        socket.send(JSON.stringify({"type":1,"body":display_username}));   
    };

    

    socket.onmessage = msg => {

        message = JSON.parse(msg.data)
        console.log(message.body[0]);
        if (message) {

            switch (message.type) {
                case 1:
                    chat = message.body[0]
                    console.log(chat);
                    addMessage("You", chat.text, "right")
                    break;
                case 0:
                    chat = message.body[0]
                    console.log(chat);
                    addMessage(chat.sender, chat.text, "left")
                    break;
                case -1:
                    if (message.body) {
                        chats = message.body
                        chats.forEach(function (msg) {
                            addMessage(msg.sender, msg.text, "left");
                        });
                    }
                    else {
                        addError("Unable to Retrieve the history","red")
                    }
                    break;
                case 2:
                    chat = message.body[0]
                    addError(chat.text,"grey")
                    break;
            }
        } else {
            addError("Unable to load message, Please refresh your page","red")
        }
    };

    socket.onclose = event => {
        console.log("Socket Closed Connection: ", event);
        addError("Connection to chat Room was closed","red")
        CONN_STATUS = 0
    };

    socket.onerror = error => {
        console.log("Socket Error: ", error);
        addError("Chat Room Server Error","red")
    };
};








