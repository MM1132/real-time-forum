$(document).ready(() => {
    // Event listener for the like
    $(".shown-like-button").click(function() {
        let postID = $(this).attr("postid");
        $.ajax({
            url: "/like?id=" + postID,
            method: "post",
            dataType: "json", // We are getting back json from the server
            success: (response) => {
                if (response.RedirectPath) {
                    $(location).attr("href", response.RedirectPath)
                    return;
                }

                $("#likes-" + postID).text(response.LikeCount);
                $(this).toggleClass("like-dislike-highlight");
                $(".shown-dislike-button[postid='" + postID + "']").removeClass("like-dislike-highlight")
            }
        });
    });

    // Event listener for the dislike
    $(".shown-dislike-button").click(function() {
        let postID = $(this).attr("postid");
        $.ajax({
            url: "/dislike?id=" + postID,
            method: "post",
            dataType: "json", // We are getting back json from the server
            success: (response) => {
                if (response.RedirectPath) {
                    $(location).attr("href", response.RedirectPath)
                    return;
                }

                $("#likes-" + postID).text(response.LikeCount);
                $(this).toggleClass("like-dislike-highlight");
                $(".shown-like-button[postid='" + postID + "']").removeClass("like-dislike-highlight")
            }
        });
    });
});