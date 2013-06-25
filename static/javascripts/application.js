jQuery(function ($) {
    var eventSource = new EventSource("/stream");

    eventSource.onmessage = function (messageEvent) {
        var message = JSON.parse(messageEvent.data);

        $("<li>").prependTo("#chat").text(message.name + ": " + message.body);
    };

    $("form").submit(function (event) {
        event.preventDefault();

        var message = $("#message").val();

        $("#message").val("");

        $.ajax("/messages", { type: "POST", data: { name: "guest", message: message } });
    });
});
