package pages

import (
	"forum/internal/forumDB"
	"forum/internal/forumEnv"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Register struct {
	forumEnv.Env
}

// Contains things that are generated for every request and passed on to the template
type registerData struct {
	forumEnv.GenericData
	Errors map[string]string
}

func (env Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We must create a new indexData struct because it can't be shared between requests
	data := registerData{}
	if err := data.InitData(env.Env, r); err != nil {
		return
	}
	if data.User.UserID != 0 { // access denied if logged in
		http.Redirect(w, r, "/board", http.StatusTemporaryRedirect)
		return
	}
	data.Errors = make(map[string]string) // errors map for validate function

	if r.Method == "POST" {
		if env.validate(r, data) { // checks for errors in form before calling register function.
			env.register(w, r)
		}
	}
	data.AddTitle("Register")

	// Finally execute the template with the data we got
	tmpl := env.Templates["register"]
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}
}

func (env Register) register(w http.ResponseWriter, r *http.Request) { // Creates a new user from POST request. Only usable in POST request.
	// Let's convert the password into a hash
	passwordHash, err := generateHash(r.FormValue("password"))
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	// username always capitalized, lowercase. Email always lowercase.
	newUser := forumDB.User{
		Name:     strings.Title(strings.ToLower(r.FormValue("username"))),
		Email:    strings.ToLower(r.FormValue("email")),
		Password: passwordHash,
	}

	_, err = env.Users.Insert(newUser)
	if err != nil {
		log.Println(err)
	}
	log.Printf("New user registered: %s\n", newUser.Name)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (env Register) validate(r *http.Request, data registerData) bool { // Checks form for errors and logs them, returns true if no errors found. Usable only in POST request.
	usernameFormat := regexp.MustCompile(`^[a-zA-Z0-9]*$`) // alphanumerical only
	emailFormat := regexp.MustCompile(`.+@+.+\..+`)        // x@x.x format

	// username errors
	r.ParseForm() // Works without this line too, but my anti-virus detecs forum.exe as trojan then. Keeping it for the sake of my sanity.

	if len(r.FormValue("username")) < 4 || len(r.FormValue("username")) > 12 { // checks if length is between 4 and 12 characters.
		if len(r.FormValue("username")) == 0 { // checks for empty username
			data.Errors["Username"] = "Username can't be empty."
		} else {
			data.Errors["Username"] = "Length must be between 4 and 12 characters."
		}
	} else if !usernameFormat.Match([]byte(r.FormValue("username"))) { // checks username format. Must be alphanumerical
		data.Errors["Username"] = "Username can contain only alphanumerical characters. Please enter a valid username."
	} else if _, err := env.Users.GetByName(r.FormValue("username")); err == nil { // checks if username is already taken
		data.Errors["Username"] = "This username has already been taken. Please choose another username and try again."
	}

	// password errors

	if len(r.FormValue("password")) < 6 || len(r.FormValue("password")) > 128 { // checks if password is between 6 and 20 characters
		if len(r.FormValue("password")) == 0 { // checks if password is empty
			data.Errors["Password"] = "Password can't be empty."
		} else {
			data.Errors["Password"] = "Length must be between 6 and 128 characters."
		}
	}

	// password2 error
	if r.FormValue("password") != r.FormValue("confirm-password") { // checks if passwords match.
		data.Errors["Password2"] = "The passwords do not match. Please try again."
	}

	// email errors

	if len(r.FormValue("email")) == 0 { // checks if email is empty
		data.Errors["Email"] = "Email can't be empty."
	} else if !emailFormat.Match([]byte(r.FormValue("email"))) { // checks email format. has to be x@x.x
		data.Errors["Email"] = "Invalid email address. Please enter a valid email adress."
	} else if _, err := env.Users.GetByEmail(r.FormValue("email")); err == nil { // checks if email is already taken
		data.Errors["Email"] = "This email address has already been registered for a different user. Please try again."
	}

	// email2 error

	if r.FormValue("email") != r.FormValue("confirm-email") { // checks if email field matches confirm email field
		data.Errors["Email2"] = "The email addresses do not match. Please try again."
	}

	return len(data.Errors) == 0 // if no errors found, this is true, otherwise false.
}
