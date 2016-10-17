$(document).ready(function(){

    console.log("The page is loading");

    var socket = new WebSocket("ws://192.168.1.218:8080/socket", "protocolOne");

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
            user: "Jordan",
            date: Date.now()
        };

        socket.send(JSON.stringify(fullMsg));

        console.log(fullMsg);
        $("#msg").val('');

    });

});



