$(document).ready(function() {

    $("#chat-close-button").click(function() {
        $('#chat').hide("slow", "swing");
    });
    $("#open-chat").click(function() {
        $('#chat').show("slow", "swing");
    });
});