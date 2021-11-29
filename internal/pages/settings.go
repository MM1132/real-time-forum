package pages

import (
	"encoding/json"
	"forum/internal/forumEnv"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Settings struct {
	forumEnv.Env
}

type SettingsData struct {
	forumEnv.GenericData
	Error   string
	Success string
}

type Response struct {
	Type    string // Whether the request was a success or not
	Message string // Something descirptive
	About   string // Where the message should be placed on the page
}

type Redirect struct {
	RedirectPath string
}

type Password struct {
	Current_password    string `json:"current_password"`
	New_password_first  string `json:"new_password_first"`
	New_password_second string `json:"new_password_second"`
}

type FormData struct {
	Password    Password `json:"password"`
	Description string   `json:"description"`
}

func (env Settings) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Initialize data
	data := SettingsData{}
	data.InitData(env.Env, r)

	// If we are using the right http method
	switch r.Method {
	case "POST":
		// Check if the user has even logged in
		// And send a response of json to the jquery to redirect them to the login page
		if data.User.UserID == 0 {
			sendInterfaceAsJson(w, Redirect{RedirectPath: "/login"})
			return
		}

		response, err := env.checkSettingsForm(w, r, data)
		if err != nil {
			log.Println(err)
			return
		}
		// But in case of no errors, send the results to the client
		err = sendInterfaceAsJson(w, response)
		if err != nil {
			log.Println(err)
			return
		}
	case "GET":
		// Check if the user has even logged in
		if data.User.UserID == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		template := env.Templates["settings"]
		if err := template.ExecuteTemplate(w, "layout", data); err != nil {
			sendErr(err, w, http.StatusInternalServerError)
		}
	}
}

func (env Settings) checkSettingsForm(w http.ResponseWriter, r *http.Request, data SettingsData) ([]Response, error) {
	formData := FormData{}

	// Limit the bytes of the body. 5000 should be enough for our small form
	r.Body = http.MaxBytesReader(w, r.Body, 5000)
	// Read in the body as bytes
	formDataBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []Response{}, err
	}
	// And get the passwords from the bytes into the struct
	err = json.Unmarshal(formDataBytes, &formData)
	if err != nil {
		return []Response{}, err
	}

	// This is to hold all the response information
	responseMessages := []Response{}

	// Password checking
	resMsg, err := env.checkPassword(data, formData.Password)
	if err != nil {
		return []Response{}, err
	}
	responseMessages = append(responseMessages, resMsg)

	// Description checking
	resMsg, err = env.checkDescription(data, formData.Description)
	if err != nil {
		return []Response{}, err
	}
	responseMessages = append(responseMessages, resMsg)

	return responseMessages, nil
}

func (env Settings) checkPassword(data SettingsData, passwordData Password) (Response, error) {
	// Password field was empty
	if len(passwordData.Current_password) < 1 {
		return Response{Type: "error", Message: "", About: "change-password"}, nil
	}

	// Currect password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(data.User.Password), []byte(passwordData.Current_password)); err != nil {
		return Response{Type: "error", Message: "Wrong current password.", About: "change-password"}, nil
	}

	// New entered password is of correct length
	if len(passwordData.New_password_first) < 6 || len(passwordData.New_password_first) > 128 {
		return Response{Type: "error", Message: "New password length must be between 6 and 128 characters.", About: "change-password"}, nil
	}

	// If the new password matches the repeated new password
	if passwordData.New_password_first != passwordData.New_password_second {
		return Response{Type: "error", Message: "Passwords did not match.", About: "change-password"}, nil
	}

	// Turn the new wanted password into a hash
	hashedPassword, err := generateHash(passwordData.New_password_first)
	if err != nil {
		return Response{}, err
	}
	err = env.Users.SetPassword(hashedPassword, data.User.UserID)
	if err != nil {
		return Response{}, err
	}

	// If we passed all the error checks, it must mean that the password has been changed successfully
	return Response{Type: "success", Message: "Password changed successfully.", About: "change-password"}, nil
}

func (env Settings) checkDescription(data SettingsData, description string) (Response, error) {
	if len(description) < 1 {
		return Response{Type: "error", Message: "", About: "change-description"}, nil
	}

	if len(description) < 10 || len(description) > 256 {
		return Response{Type: "error", Message: "Description length must be between 8 and 256 characters.", About: "change-description"}, nil
	}

	err := env.Users.SetDescription(description, data.User.UserID)
	if err != nil {
		return Response{}, err
	}
	return Response{Type: "success", Message: "Description changed successfully.", About: "change-description"}, nil
}
