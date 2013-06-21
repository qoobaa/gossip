jQuery(function ($) {
    var eventSource = new EventSource("/stream");

    eventSource.onmessage = function (messageEvent) {
        $("<ul>").appendTo("#chat").text(messageEvent.data);
    }

    $("form").submit(function (event) {
        event.preventDefault();

        var message = $("#message").val();

        $("#message").val("");

        $.ajax("/messages", { type: "POST", data: { name: "guest", message: message } });
    });
});
