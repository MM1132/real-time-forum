function atLeastOneFilled(fields) {
    for (i in fields) {
        if (fields[i].length > 0) {
            return true;
        }
    }
    return false;
}

$(document).ready(function() {
    // Listener for the submit button
    $("#save-settings").click(function() {
        data = {};

        // Password
        password = {
            current_password: $("#current-password").val(),
            new_password_first: $("#new-password-first").val(),
            new_password_second: $("#new-password-second").val()
        };
        if (atLeastOneFilled(password)) {
            data.password = password;
        }

        // Description
        description = $("#new-description").val();
        if (description.length > 0) {
            data.description = description;
        }

        $.ajax({
            url: "/settings",
            method: "post",
            data: JSON.stringify(data),
            contentType: "application/json",
            dataType: "json",
            success: (response) => {
                // If we get a redirect as response
                if (response.RedirectPath) {
                    $(location).attr("href", response.RedirectPath)
                }

                // Clear all the fields of the settings form
                $(".settings-input-text>input").val("");
                // Also clear both the message fields
                $("#change-password").text("");
                $("#change-description").text("");

                // Handle the output we got back
                for (v of response) {
                    $("#" + v.About).text(v.Message);
                    $("#" + v.About).attr("class", v.Type);
                }
            }
        });
    });
});