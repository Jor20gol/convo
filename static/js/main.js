$(document).ready(function(){

    if (!sessionStorage.username) {
        $("#name-save-modal").modal("show");
    }


    console.log("The page is loading");

    var socket = new WebSocket("ws://localhost:8080/socket");

    socket.onopen = function (event) {
        var test = {
            type: "test",
            text: "Is the connection open",
            user: "client",
            date: Date.now()
        };

        socket.send(JSON.stringify(test)); 
        console.log("Makeing connection to server...");
    };

    socket.onmessage = function(event) {
        var msg = JSON.parse(event.data);
        if (msg.type == "test" || msg.Type == "test") {

            console.log("Connection Successful...");

        } else if (msg.type == "message") {

            $("#messages").append("<li class=\"list-group-item message-card\"><span class=\"message-user\">" + msg.user + "</span><span class=\"message-date\"> - " + Date(msg.date).substring(4,24) + "</span><br /><span class=\"message-body\">" + msg.text + "</span></li>");

        } else {
            console.log(msg);
        }
    }

    $("#j-form").on("submit", function(e){
        e.preventDefault();
        console.log("FORM SUBMIT");
        var msgText = $("#msg").val();

        var fullMsg = {
            type: "message",
            text: msgText,
            user: sessionStorage.username,
            date: Date.now()
        };

        socket.send(JSON.stringify(fullMsg));

        console.log(fullMsg);
        $("#msg").val('');

    });

});

function saveName() {
    var name = $("#name-input").val();
    if (name.length > 30) {
        name = name.substring(0, 30);
    }
    sessionStorage.username = name;
    $("#name-save-modal").modal("hide");
}

