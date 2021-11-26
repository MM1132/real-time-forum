package pages

import (
	"errors"
	"forum/internal/forumEnv"
	"log"
	"net/http"
)

type Settings struct {
	forumEnv.Env
}

type SettingsData struct {
	forumEnv.GenericData
	Error   string
	Success string
}

func (env Settings) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := SettingsData{}
	if err := data.InitData(env.Env, r); err != nil {
		return
	}

	// Check if the user has even loged in
	if data.Session.UserID == 0 {
		sendErr(errors.New("Please log in!"), w, http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		err := env.checkForm(w, r, data)
		if err != nil {
			data.Error = err.Error()
			log.Println(err)
		} else {
			data.Success = "Settings updated successfully!"
		}
	}

	template := env.Templates["settings"]
	if err := template.ExecuteTemplate(w, "layout", data); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}
}

func (env Settings) checkForm(w http.ResponseWriter, r *http.Request, data SettingsData) error {
	// If we are loged in though, get the information from the form, about the password fields
	password_1 := r.FormValue("new-password-first")
	password_2 := r.FormValue("new-password-second")

	// Checks for password length
	if len(password_1) < 6 || len(password_1) > 128 {
		return errors.New("Password length must be between 6 and 128 characters")
	}

	// Check if the passwords match
	if password_1 != password_2 {
		return errors.New("Passwords do not match")
	}

	// Set the new password for the user
	err := env.Users.SetPassword(password_1, data.Session.UserID)
	if err != nil {
		return err
	}
	return nil
}
