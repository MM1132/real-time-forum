$(document).ready(function() {

    $("#close-chat").click(function() {
        $('#chat').hide("slow", "swing");
    });
    $("#open-chat").click(function() {
        $('#chat').show("slow", "swing");
    });
});