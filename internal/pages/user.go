package pages

import (
	"errors"
	"fmt"
	fdb "forum/internal/forumDB"
	"forum/internal/forumEnv"
	"forum/internal/search"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type User struct {
	forumEnv.Env
}

type UserData struct {
	forumEnv.GenericData
	User         fdb.User
	UserPosts    *search.PostSearch
	UserLikes    *search.PostSearch
	UserDislikes *search.PostSearch
	Error        string
}

func (env User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := UserData{}
	data.InitData(env.Env, r)

	// Get the id of the user
	idInt, err := GetQueryInt("id", r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Get all the information about the user by its ID
	if user, err := env.Users.Get(idInt); err != nil {
		sendErr(err, w, http.StatusBadRequest)
		return
	} else if err := env.Users.SetExtras(&user); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	} else {
		data.User = user
	}

	if r.Method == "POST" {
		// Handle the profile picture uploading
		if err := env.profilePictureUpload(w, r, data); err != nil {
			log.Println(err.Error())
			data.Error = err.Error()
		} else {

			http.Redirect(w, r, r.URL.RequestURI(), http.StatusSeeOther)
			return
		}
	}

	// Newest posts by this user
	data.UserPosts = search.NewPostSearch(env.Env, data.CurrentURL, data.GenericData.User.UserID)
	config(data.UserPosts)
	data.UserPosts.AuthorID.Int64 = int64(data.User.UserID)
	data.UserPosts.AuthorID.Valid = true
	if err := data.UserPosts.DoSearch(env.Env); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	// Posts recently liked by this user
	data.UserLikes = search.NewPostSearch(env.Env, data.CurrentURL, data.GenericData.User.UserID)
	config(data.UserLikes)
	data.UserLikes.ProcessOrder("likeDate")
	data.UserLikes.LikedByID.Int64 = int64(data.User.UserID)
	data.UserLikes.LikedByID.Valid = true
	if err := data.UserLikes.DoSearch(env.Env); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	// Posts recently disliked by this user
	data.UserDislikes = search.NewPostSearch(env.Env, data.CurrentURL, data.GenericData.User.UserID)
	config(data.UserDislikes)
	data.UserDislikes.ProcessOrder("likeDate")
	data.UserDislikes.DislikedByID.Int64 = int64(data.User.UserID)
	data.UserDislikes.DislikedByID.Valid = true
	if err := data.UserDislikes.DoSearch(env.Env); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	// And finally we are executing the template with the data we got
	template := env.Templates["user"]
	if err := template.ExecuteTemplate(w, "layout", data); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}
}

// Common configs for all searches
func config(s *search.PostSearch) {
	s.Breadcrumbs = true
	s.Page = -1
	s.PageLen = 3
}

func (env User) profilePictureUpload(w http.ResponseWriter, r *http.Request, data UserData) error {
	// CHeck if the user has loged in or not
	if data.Session.UserID == 0 {
		return errors.New("Please log in to upload picture")
	}

	// Set the limit to upload size
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 Mb

	// Here we can use ParseMultipartForm to check for the size limit, or any other possible form errors
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		return errors.New("File size over 1MB")
	}

	// Just read the data from the form
	file, _, err := r.FormFile("profile-picture")
	if err != nil {
		return errors.New("Choose a file please!")
	}

	// Create a buffer to store the header of the file in
	fileHeader := make([]byte, 512)

	// Copy the headers into the FileHeader buffer
	if _, err := file.Read(fileHeader); err != nil {
		return err
	}

	// Set position in the file back to start.
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	// Checking the file content type
	if http.DetectContentType(fileHeader) != "image/png" {
		return errors.New("Uploaded file is not image/png")
	}

	// Also, check the dimensions of the image
	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	if img.Bounds().Max.X < 256 || img.Bounds().Max.Y < 256 {
		return errors.New("Image must be at least 256 pixels wide and high")
	}

	if img.Bounds().Max.X != img.Bounds().Max.Y {
		return errors.New("Image height does not match it's width")
	}

	// Get the information about user's current image
	// The first element is the user's ID, the second one is the counter of the image itself
	currentImageName := strings.SplitN(data.User.Image, "-", 2)
	/* imageUserID, err := strconv.Atoi(currentImageName[0])
	if err != nil {
		return err
	} */
	imageID, err := strconv.Atoi(strings.TrimSuffix((currentImageName[1]), ".png"))
	if err != nil {
		return err
	}

	// Just creating the user image filename for greater readability purposes
	newImageFilename := strconv.Itoa(data.User.UserID) + "-" + strconv.Itoa(imageID+1) + ".png"

	// Create the image file locally
	localImage, err := os.Create("./server/static/profile-pictures/" + newImageFilename)
	if err != nil {
		return errors.New("Failed to create local profile picture for the user")
	}
	defer localImage.Close()

	// Write the image data to the file
	err = png.Encode(localImage, img)
	if err != nil {
		return errors.New("Failed to Encode the png image profile-picture")
	}

	// Also update the Image field of the user in the database
	err = env.Users.SetImage(newImageFilename, data.User.UserID)
	if err != nil {
		log.Println(fmt.Errorf("error removing old profile picture: %w", err))
	}

	// Here, delete the old image file, unless the previous image ID was 0
	if imageID != 0 {
		err = os.Remove("./server/static/profile-pictures/" + data.User.Image)
		if err != nil {
			return err
		}
	}

	return nil
}
