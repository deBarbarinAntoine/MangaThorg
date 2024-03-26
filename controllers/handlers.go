package controllers

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"log/slog"
	"mangathorg/internal/api"
	"mangathorg/internal/middlewares"
	"mangathorg/internal/models"
	"mangathorg/internal/utils"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
)

// rootHandlerGet
//
//	@Description: manages the / route and redirects to the principal page.
func rootHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())

	http.Redirect(w, r, "/principal", http.StatusSeeOther)
}

// errorHandler
//
//	@Description: displays the Error404 page.
func errorHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	log.Println("HTTP ApiErr", http.StatusNotFound)
	w.WriteHeader(http.StatusNotFound)
	utils.Logger.Warn("errorHandler", slog.Int("req_id", middlewares.LogId), slog.String("req_url", r.URL.String()), slog.Int("http_status", http.StatusNotFound))

	var data struct {
		IsConnected bool
		Username    string
		AvatarImg   string
	}

	session, sessionId := utils.GetSession(r)
	if sessionId != "" {
		data.IsConnected = true
		data.Username = session.Username
	}

	user, ok := utils.SelectUser(session.Username)
	if !ok {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
		data.IsConnected = false
	}

	data.AvatarImg = user.Avatar

	tmpl, err := template.ParseFiles(utils.Path+"templates/base.gohtml", utils.Path+"templates/header-line2.gohtml", utils.Path+"templates/error404.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

// loginHandlerGet
//
//	@Description: displays the login form and possible messages according to the
//	`err` or `status` query keys.
func loginHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	var message template.HTML
	if r.URL.Query().Has("err") {
		switch r.URL.Query().Get("err") {
		case "login":
			message = "<div class=\"message\">Wrong username or password!</div>"
		case "restricted":
			message = "<div class=\"message\">You need to login to access that area!</div>"
		}
	} else if r.URL.Query().Has("status") {
		switch r.URL.Query().Get("status") {
		case "update-pwd":
			message = `<div class="message">Your password has been updated!</div>`
		case "signed-up":
			message = `<div class="message">We've sent you a message to confirm your account!</div>`
		}
	}
	var data = struct {
		Message template.HTML
	}{
		Message: message,
	}
	tmpl, err := template.ParseFiles(utils.Path+"templates/base.gohtml", utils.Path+"templates/login.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

// loginHandlerPost
//
//	@Description: login treatment handler
func loginHandlerPost(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	credentials := models.Credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	if utils.CheckPwd(credentials) {
		utils.OpenSession(&w, credentials.Username, r)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login?err=login", http.StatusSeeOther)
	}
}

// registerHandlerGet
//
//	@Description: displays the register form and possible messages according to the
//	`err` query key.
func registerHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	var message template.HTML
	if r.URL.Query().Has("err") {
		switch r.URL.Query().Get("err") {
		case "username":
			message = "<div class=\"message\">Username must be at least 3 characters long!</div>"
		case "user":
			message = "<div class=\"message\">Username or email already used!</div>"
		case "match":
			message = "<div class=\"message\">Both passwords need to be equal!</div>"
		case "email":
			message = "<div class=\"message\">Wrong email value!</div>"
		case "password":
			message = "<div class=\"message\">Password needs 8 characters min, 1 digit, 1 lowercase, 1 uppercase and 1 symbol.</div>"
		}
	}
	var data = struct {
		Message template.HTML
	}{
		Message: message,
	}
	tmpl, err := template.ParseFiles(utils.Path+"templates/base.gohtml", utils.Path+"templates/register.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

// registerHandlerPost
//
//	@Description: register treatment handler.
func registerHandlerPost(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	formValues := struct {
		username  string
		email     string
		password1 string
		password2 string
	}{
		username:  r.FormValue("username"),
		email:     strings.TrimSpace(strings.ToLower(r.FormValue("email"))),
		password1: r.FormValue("password"),
		password2: r.FormValue("confirm-password"),
	}
	switch {
	case len(formValues.username) < 3:
		http.Redirect(w, r, "register?err=username", http.StatusSeeOther)
		return
	case !utils.CheckUser(models.User{
		Username: formValues.username,
		Email:    formValues.email,
	}):
		http.Redirect(w, r, "register?err=user", http.StatusSeeOther)
		return
	case formValues.password1 != formValues.password2:
		http.Redirect(w, r, "register?err=match", http.StatusSeeOther)
		return
	case !utils.CheckEmail(formValues.email):
		http.Redirect(w, r, "register?err=email", http.StatusSeeOther)
		return
	case !utils.CheckPasswd(formValues.password1):
		http.Redirect(w, r, "register?err=password", http.StatusSeeOther)
		return
	}
	hash, salt := utils.NewPwd(formValues.password1)
	newTempUser := models.TempUser{
		ConfirmID:    "",
		CreationTime: time.Now(),
		User: models.User{
			Id:        0,
			Username:  formValues.username,
			HashedPwd: hash,
			Salt:      salt,
			Email:     formValues.email,
		},
	}
	utils.SendMail(&newTempUser, "creation")
	utils.TempUsers = append(utils.TempUsers, newTempUser)
	http.Redirect(w, r, "/login?status=signed-up", http.StatusSeeOther)
}

// forgotPasswordHandlerGet
//
//	@Description: displays the forgot password's form and the possible
//	confirmation message.
func forgotPasswordHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	var message template.HTML
	if r.URL.Query().Has("confirm") {
		message = "<div class=\"message\">A mail has been sent to set a new password!</div>"
	}
	var data = struct {
		Message template.HTML
	}{
		Message: message,
	}
	tmpl, err := template.ParseFiles(utils.Path+"templates/base.gohtml", utils.Path+"templates/forgot-passwd.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

// forgotPasswordHandlerPost
//
//	@Description: forgot password's form treatment handler.
func forgotPasswordHandlerPost(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())

	email := r.FormValue("email")
	if exists, user := utils.EmailExists(email); exists {
		var temp = models.TempUser{
			CreationTime: time.Now(),
			User:         user,
		}
		utils.SendMail(&temp, "lost")
		utils.LostUsers = append(utils.LostUsers, temp)
	}

	http.Redirect(w, r, "/forgot-password?confirm=true", http.StatusSeeOther)
}

// updateCredentialsHandlerGet
//
//	@Description: displays the update credentials form (for those who forgot their
//	password and accessed it through the URL sent by mail).
func updateCredentialsHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	var id string
	if r.URL.Query().Has("id") {
		id = r.URL.Query().Get("id")
	}
	var message template.HTML
	if r.URL.Query().Has("err") {
		switch r.URL.Query().Get("err") {
		case "match":
			message = "<div class=\"message\">Both passwords need to be equal!</div>"
		case "password":
			message = "<div class=\"message\">Password needs 8 characters min, 1 digit, 1 lowercase, 1 uppercase and 1 symbol.</div>"
		default:
			message = "<div class=\"message\">An error has occured!</div>"
		}
	}
	var data = struct {
		Message template.HTML
		Id      string
	}{
		Message: message,
		Id:      id,
	}
	tmpl, err := template.ParseFiles(utils.Path+"templates/base.gohtml", utils.Path+"templates/update-credentials.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

// updateCredentialsHandlerPost
//
//	@Description: update credentials form's treatment handler.
func updateCredentialsHandlerPost(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	id := r.PathValue("id")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")
	var lostUser models.TempUser
	for _, user := range utils.LostUsers {
		if user.ConfirmID == id {
			lostUser = user
			break
		}
	}
	if password != confirmPassword || !utils.CheckPasswd(password) {
		http.Redirect(w, r, "/update-credentials?id="+id, http.StatusSeeOther)
		return
	}
	lostUser.User.HashedPwd, lostUser.User.Salt = utils.NewPwd(password)
	utils.UpdateLostUser(lostUser)
	http.Redirect(w, r, "/login?status=update-pwd", http.StatusSeeOther)
}

// profileHandlerGet
//
//	@Description: displays the profile page with its update user form and possible
//	messages according to the `err` and `status` query keys.
func profileHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())

	var message template.HTML
	if r.URL.Query().Has("err") {
		switch r.URL.Query().Get("err") {
		case "match":
			message = "<div class=\"message\">Both passwords need to be equal!</div>"
		case "password":
			message = "<div class=\"message\">Password needs 8 characters min, 1 digit, 1 lowercase, 1 uppercase and 1 symbol.</div>"
		case "current-pwd":
			message = "<div class=\"message\">Incorrect password!</div>"
		default:
			message = "<div class=\"message\">An error has occured!</div>"
		}
	} else if r.URL.Query().Has("status") {
		if r.URL.Query().Get("status") == "updated" {
			message = "<div class=\"message\">Your information has been successfully updated!</div>"
		} else if r.URL.Query().Get("status") == "nothing" {
			message = "<div class=\"message\">Nothing has been changed!</div>"
		}
	}

	session, _ := utils.GetSession(r)
	user, ok := utils.SelectUser(session.Username)
	if !ok {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
	}

	var data = struct {
		IsConnected bool
		Username    string
		Message     template.HTML
		AvatarImg   string
		Avatars     []string
	}{
		IsConnected: true,
		Username:    user.Username,
		Message:     message,
		AvatarImg:   user.Avatar,
	}

	for i := range 86 {
		var nb string
		if i+1 < 10 {
			nb = "00" + strconv.Itoa(i+1)
		} else {
			nb = "0" + strconv.Itoa(i+1)
		}
		data.Avatars = append(data.Avatars, "profile-avatar-"+nb+".jpg")
	}

	tmpl, err := template.ParseFiles(utils.Path+"templates/base.gohtml", utils.Path+"templates/header-line2.gohtml", utils.Path+"templates/profile.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

// profileHandlerPost
//
//	@Description: profile form's treatment handler.
func profileHandlerPost(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())

	avatar := r.FormValue("avatar")
	password := r.FormValue("password")
	newPassword := r.FormValue("new-password")
	confirmPassword := r.FormValue("confirm-password")

	session, _ := utils.GetSession(r)

	user, ok := utils.SelectUser(session.Username)
	if !ok {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
		http.Redirect(w, r, "/profile?err=internal-error", http.StatusSeeOther)
		return
	}

	if password != "" && newPassword != "" && confirmPassword != "" {
		if !utils.CheckPwd(models.Credentials{Username: session.Username, Password: password}) {
			http.Redirect(w, r, "/profile?err=current-pwd", http.StatusSeeOther)
			return
		}
		if newPassword != confirmPassword {
			http.Redirect(w, r, "/profile?err=match", http.StatusSeeOther)
			return
		} else if !utils.CheckPasswd(newPassword) {
			http.Redirect(w, r, "/profile?err=password", http.StatusSeeOther)
			return
		}
		user.HashedPwd, user.Salt = utils.NewPwd(newPassword)
	} else if user.Avatar == avatar {
		http.Redirect(w, r, "/profile?status=nothing", http.StatusSeeOther)
		return
	}

	user.Avatar = avatar
	utils.UpdateUser(user)

	http.Redirect(w, r, "/profile?status=updated", http.StatusSeeOther)
}

// homeHandlerGet
//
//	@Description: display the user's home page.
func homeHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())

	session, _ := utils.GetSession(r)
	user, ok := utils.SelectUser(session.Username)

	if !ok {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
		http.Redirect(w, r, "/login?err=restricted", http.StatusSeeOther)
		return
	}

	var data = struct {
		Order        string
		Path         string
		CreationTime time.Time
		Username     string
		Email        string
		AvatarImg    string
		HasBanner    bool
		Banner       models.MangaUsefullData
		HasFavorites bool
		Favorites    []models.MangaUsefullData
	}{
		Order:        "desc",
		Path:         "static",
		CreationTime: user.CreationTime,
		Username:     user.Username,
		Email:        user.Email,
		AvatarImg:    user.Avatar,
		Banner:       api.FetchMangaById(user.MangaBanner.Id, "desc", 0),
		Favorites:    api.FetchMangasById(user.Favorites, "desc", 0),
	}

	data.HasFavorites = data.Favorites != nil && len(data.Favorites) > 0
	data.HasBanner = !reflect.DeepEqual(data.Banner, models.MangaUsefullData{})

	tmpl, err := template.ParseFiles(utils.Path+"templates/home.gohtml", utils.Path+"templates/base.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

// TESTING only.
//
// logHandlerGet
//
//	@Description: fetches the logs according to the query sent.
func logHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Query().Has("level") {
		err := json.NewEncoder(w).Encode(utils.FetchAttrLogs("level", r.URL.Query().Get("level")))
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
			return
		}
		return
	} else if r.URL.Query().Has("user") {
		err := json.NewEncoder(w).Encode(utils.FetchAttrLogs("user", r.URL.Query().Get("user")))
		if err != nil {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
			return
		}
		return
	}
	err := json.NewEncoder(w).Encode(utils.RetrieveLogs())
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		return
	}
}

// confirmHandlerGet
//
//	@Description: displays the new account's confirmation page.
func confirmHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	if r.URL.Query().Has("id") {
		id := r.URL.Query().Get("id")
		utils.PushTempUser(id)

		tmpl, err := template.ParseFiles(utils.Path+"templates/confirm.gohtml", utils.Path+"templates/base.gohtml")
		if err != nil {
			log.Fatalln(err)
		}
		err = tmpl.ExecuteTemplate(w, "base", nil)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// logoutHandlerGet
//
//	@Description: logs the user out.
func logoutHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	utils.Logout(&w, r)
	http.Redirect(w, r, "/principal", http.StatusSeeOther)
}

// principalHandlerGet
//
//	@Description: displays the principal page.
func principalHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	tmpl, err := template.ParseFiles(utils.Path+"templates/principal.gohtml", utils.Path+"templates/header-line2.gohtml", utils.Path+"templates/base.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	var data = struct {
		IsConnected    bool
		Username       string
		AvatarImg      string
		Banner         models.MangaUsefullData
		LatestUploaded []models.MangaUsefullData
		Popular        []models.MangaUsefullData
	}{
		Banner:         api.FetchMangaById("cb676e05-8e6e-4ec4-8ba0-d3cb4f033cfa", "asc", 1),
		LatestUploaded: api.FetchManga(api.TopLatestUploadedRequest).Mangas,
		Popular:        api.FetchManga(api.TopPopularRequest).Mangas,
	}

	session, _ := utils.GetSession(r)
	data.Username = session.Username
	data.IsConnected = api.AddSingleFavoriteInfo(r, &data.Banner)
	_ = api.AddFavoriteInfo(r, &data.LatestUploaded)
	_ = api.AddFavoriteInfo(r, &data.Popular)

	user, ok := utils.SelectUser(session.Username)
	if !ok {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
	}
	data.AvatarImg = user.Avatar

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

// mangaHandlerGet
//
//	@Description: displays the manga page according to the id put in the URL.
func mangaHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	mangaId := r.PathValue("id")
	if mangaId == "" {
		http.Redirect(w, r, "/principal", http.StatusNotFound)
		return
	}
	var order, pagination string
	if r.URL.Query().Has("order") {
		order = r.URL.Query().Get("order")
	} else {
		order = "desc"
	}
	if r.URL.Query().Has("pag") {
		pagination = r.URL.Query().Get("pag")
	} else {
		pagination = "1"
	}
	pag, errAtoi := strconv.Atoi(pagination)
	if errAtoi != nil {
		pag = 1
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	if pag < 1 {
		pag = 1
	}
	tmpl, err := template.ParseFiles(utils.Path+"templates/manga.gohtml", utils.Path+"templates/header-line2.gohtml", utils.Path+"templates/base.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	offset := (pag - 1) * 15
	manga := api.FetchMangaById(mangaId, order, offset)
	var pages []int
	pageMax := manga.NbChapter / 15
	if manga.NbChapter%15 > 0 {
		pageMax++
	}
	for i := range pageMax {
		pages = append(pages, i+1)
	}
	var data = struct {
		IsConnected bool
		Username    string
		AvatarImg   string
		Manga       models.MangaUsefullData
		CurrentPage int
		Pages       []int
		Order       string
	}{
		Manga:       manga,
		CurrentPage: pag,
		Pages:       pages,
		Order:       order,
	}

	//  check if the manga was found, and if not, show the error404 page.
	if reflect.DeepEqual(manga, models.MangaUsefullData{}) {
		http.Redirect(w, r, "/error404", http.StatusSeeOther)
		return
	}

	session, _ := utils.GetSession(r)
	data.Username = session.Username
	data.IsConnected = api.AddSingleFavoriteInfo(r, &data.Manga)

	user, ok := utils.SelectUser(session.Username)
	if !ok {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
	}
	data.AvatarImg = user.Avatar

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
}

// coverHandlerGet
//
//	@Description: sends the manga's cover image according to the mangaId and the
//	img name.
func coverHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	mangaId := r.PathValue("manga")
	img := r.PathValue("img")
	if mangaId == "" || img == "" {
		http.Error(w, "invalid request: empty value", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	_, err := w.Write(api.ImageProxy(mangaId, img))
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		http.Error(w, "cover image not found", http.StatusNotFound)
	}
}

// scanHandlerGet
//
//	@Description: sends the scan's image according to the chapterId, quality, hash
//	and image name.
func scanHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	chapterId := r.PathValue("chapterId")
	quality := r.PathValue("quality")
	hash := r.PathValue("hash")
	img := r.PathValue("img")
	if chapterId == "" || quality == "" || hash == "" || img == "" {
		log.Println("empty value")
		http.Error(w, "invalid request: empty value", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	i, err := w.Write(api.ScanProxy(chapterId, quality, hash, img))
	if err != nil || i == 0 {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		http.Error(w, "scan image not found", http.StatusNotFound)
	}
}

// favoriteHandlerPost
//
//	@Description: adds a favorite to a user according to the mangaId sent in the
//	URL.
func favoriteHandlerPost(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	mangaId := r.PathValue("mangaId")

	session, sessionId := utils.GetSession(r)

	if sessionId == "" {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("sessionId not found")))
		http.Error(w, "restricted access: you need a valid session to proceed", http.StatusUnauthorized)
		return
	}

	if mangaId == "" {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("mangaId is null")))
		http.Error(w, "you need to provide a mangaId", http.StatusNotFound)
		return
	}

	user, ok := utils.SelectUser(session.Username)
	if !ok {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
		http.Error(w, "restricted access: you need a valid user to proceed", http.StatusUnauthorized)
		return
	}

	for _, favorite := range user.Favorites {
		if mangaId == favorite.Id {
			http.Error(w, "the manga is already present in the favorites", http.StatusConflict)
			return
		}
	}

	user.Favorites = append(user.Favorites, models.MangaUser{Id: mangaId})
	utils.UpdateUser(user)

	w.Header().Set("result", "Manga added successfully")
	w.WriteHeader(http.StatusOK)
}

// favoriteHandlerDelete
//
//	@Description: removes a favorite from a user according to the mangaId sent in
//	the URL.
func favoriteHandlerDelete(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	mangaId := r.PathValue("mangaId")

	session, sessionId := utils.GetSession(r)

	if sessionId == "" {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("sessionId not found")))
		http.Error(w, "restricted access: you need a valid session to proceed", http.StatusUnauthorized)
		return
	}

	if mangaId == "" {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("mangaId is null")))
		http.Error(w, "you need to provide a mangaId", http.StatusNotFound)
		return
	}

	user, ok := utils.SelectUser(session.Username)
	if !ok {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
		http.Error(w, "restricted access: you need a valid user to proceed", http.StatusUnauthorized)
		return
	}

	for i, favorite := range user.Favorites {
		if mangaId == favorite.Id {
			user.Favorites = append(user.Favorites[:i], user.Favorites[i+1:]...)
			utils.UpdateUser(user)
			w.Header().Set("result", "Manga deleted successfully")
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("error", "The manga was not found in the favorites")
}

// bannerHandlerPut
//
//	@Description: modifies a user's banner according to the mangaId sent in the
//	URL.
func bannerHandlerPut(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	mangaId := r.PathValue("mangaId")

	session, sessionId := utils.GetSession(r)

	if sessionId == "" {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("sessionId not found")))
		http.Error(w, "restricted access: you need a valid session to proceed", http.StatusUnauthorized)
		return
	}

	if mangaId == "" {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("mangaId is null")))
		http.Error(w, "you need to provide a mangaId", http.StatusNotFound)
		return
	}

	user, ok := utils.SelectUser(session.Username)
	if !ok {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
		http.Error(w, "restricted access: you need a valid user to proceed", http.StatusUnauthorized)
		return
	}

	for _, favorite := range user.Favorites {
		if mangaId == favorite.Id {
			user.MangaBanner = models.MangaUser{
				Id:              mangaId,
				LastChapterRead: favorite.LastChapterRead,
			}
			utils.UpdateUser(user)
			w.Header().Set("result", "Banner updated successfully")
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("error", "The manga was not found in the favorites")
}

// chapterHandlerGet
//
//	@Description: displays the requested chapter's scans to read it, according to
//	the mangaId, chapterNb and chapterId sent in the URL.
func chapterHandlerGet(w http.ResponseWriter, r *http.Request) {

	log.Println(utils.GetCurrentFuncName())

	mangaId := r.PathValue("mangaId")
	offsetString := r.PathValue("offset")
	chapterId := r.PathValue("chapterId")

	if chapterId == "" || mangaId == "" || offsetString == "" {
		http.Redirect(w, r, "/error404", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles(utils.Path+"templates/chapter.gohtml", utils.Path+"templates/header-line2.gohtml", utils.Path+"templates/base.gohtml")
	if err != nil {
		log.Fatalln(err)
	}

	scan := api.ScanRequest(chapterId)

	var apiMangaFeed models.ApiMangaFeed
	offset, errConv := strconv.Atoi(offsetString)
	if errConv != nil {
		offset = 0
	}
	limit := 10
	reqOffset := offset - 5
	if reqOffset < 0 {
		limit = 10 + reqOffset
		reqOffset = 0
	}

	var query = make(url.Values)
	query.Add("order[chapter]", "asc")
	query.Add("translatedLanguage[]", "en")
	query.Add("contentRating[]", "safe")
	query.Add("includes[]", "scanlation_group")
	query.Add("limit", strconv.Itoa(limit))
	query.Add("offset", strconv.Itoa(reqOffset))

	err = apiMangaFeed.SendRequest(api.BaseApiURL, "manga/"+mangaId+"/feed", query)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}

	currentInd := offset - reqOffset
	chapters := apiMangaFeed.Format()

	chapterNb := chapters[currentInd].Chapter
	var previous, next int
	for i, chapter := range chapters {
		if chapter.Chapter == chapterNb {
			previous = i - 1
			break
		}
	}
	for i := len(chapters) - 1; i >= 0; i-- {
		if chapters[i].Chapter == chapterNb {
			next = i + 1
			break
		}
	}

	var isPrevious, isNext bool = true, true

	if previous < 0 {
		isPrevious = false
	}

	if (offset + (next - currentInd)) >= apiMangaFeed.Total {
		isNext = false
	}

	var previousLink, nextLink string

	if isPrevious {
		previousLink = "/chapter/" + mangaId + "/" + strconv.Itoa(offset-(currentInd-previous)) + "/" + chapters[previous].Id
	}

	if isNext {
		nextLink = "/chapter/" + mangaId + "/" + strconv.Itoa(offset+(next-currentInd)) + "/" + chapters[next].Id
	}

	var data = struct {
		IsConnected     bool
		Username        string
		AvatarImg       string
		IsOk            bool
		ToDataSaver     bool
		MangaId         string
		Manga           string
		ChapterNb       string
		IsPrevious      bool
		IsNext          bool
		Previous        string
		Next            string
		ScanlationGroup string
		Id              string
		Quality         string
		Alt             string
		Scan            struct {
			Hash      string
			Data      []string
			DataSaver []string
		}
	}{
		MangaId:         mangaId,
		Manga:           api.FetchMangaById(mangaId, "desc", 0).Title,
		ChapterNb:       chapterNb,
		IsPrevious:      isPrevious,
		IsNext:          isNext,
		Previous:        previousLink,
		Next:            nextLink,
		ScanlationGroup: chapters[currentInd].ScanlationGroup,
		Id:              chapterId,
		Quality:         "data",
		Alt:             "",
		Scan: struct {
			Hash      string
			Data      []string
			DataSaver []string
		}{
			Hash:      scan.Chapter.Hash,
			Data:      scan.Chapter.Data,
			DataSaver: scan.Chapter.DataSaver,
		},
	}

	session, sessionId := utils.GetSession(r)
	if sessionId != "" {
		data.IsConnected = true
		data.Username = session.Username
	}

	user, ok := utils.SelectUser(session.Username)
	if !ok {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
	}
	data.AvatarImg = user.Avatar

	if len(data.Scan.Data) == 0 {
		if len(data.Scan.DataSaver) == 0 {
			data.IsOk = false
		} else {
			data.ToDataSaver = true
		}
	} else {
		data.IsOk = true
		data.ToDataSaver = false
	}
	data.Alt = data.Manga + " - Ch. " + data.ChapterNb
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

// tagsHandlerGet
//
//	@Description: displays all the available Tag sorted by type.
func tagsHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())

	sortedTags := api.FetchSortedTags()

	var data = struct {
		IsConnected bool
		Username    string
		AvatarImg   string
		FormatTags  []models.ApiTag
		GenreTags   []models.ApiTag
		ThemeTags   []models.ApiTag
		PublicTags  []string
		StatusTags  []string
	}{
		FormatTags: sortedTags.FormatTags,
		GenreTags:  sortedTags.GenreTags,
		ThemeTags:  sortedTags.ThemeTags,
		PublicTags: sortedTags.PublicTags,
		StatusTags: sortedTags.StatusTags,
	}

	session, sessionId := utils.GetSession(r)
	if sessionId != "" {
		data.IsConnected = true
		data.Username = session.Username
	}

	user, ok := utils.SelectUser(session.Username)
	if !ok {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
	}
	data.AvatarImg = user.Avatar

	tmpl, err := template.ParseFiles(utils.Path+"templates/tags.gohtml", utils.Path+"templates/header-line2.gohtml", utils.Path+"templates/base.gohtml")
	if err != nil {
		log.Fatalln(err)
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

// categoryHandlerGet
//
//	@Description: displays the mangas matching a specific Tag (Format, Genre or
//	Theme only) which id is sent in the URL.
func categoryHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	tagId := r.PathValue("tagId")
	if tagId == "" {
		http.Redirect(w, r, "/error404", http.StatusNotFound)
		return
	}
	var order, pagination string
	if r.URL.Query().Has("order") {
		order = r.URL.Query().Get("order")
	} else {
		order = "desc"
	}
	if r.URL.Query().Has("pag") {
		pagination = r.URL.Query().Get("pag")
	} else {
		pagination = "1"
	}
	pag, errAtoi := strconv.Atoi(pagination)
	if errAtoi != nil {
		pag = 1
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	if pag < 1 {
		pag = 1
	}
	offset := (pag - 1) * 18

	var request = models.MangaRequest{
		OrderType:    "rating",
		OrderValue:   order,
		IncludedTags: []string{tagId},
		ExcludedTags: nil,
		Limit:        18,
		Offset:       offset,
	}

	var data = struct {
		IsConnected bool
		Username    string
		AvatarImg   string
		Path        string
		Name        string
		Response    models.MangasInBulk
		CurrentPage int
		TotalPages  int
		Order       string
		Previous    int
		Next        int
	}{
		AvatarImg:   "avatar.jpg",
		Path:        "../static",
		Name:        api.TagSelect(tagId).Attributes.Name.En,
		Response:    api.FetchManga(request),
		CurrentPage: pag,
		Order:       order,
		Previous:    pag - 1,
		Next:        pag + 1,
	}

	session, _ := utils.GetSession(r)
	data.Username = session.Username
	data.IsConnected = api.AddFavoriteInfo(r, &data.Response.Mangas)

	user, ok := utils.SelectUser(session.Username)
	if !ok {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
	}
	data.AvatarImg = user.Avatar

	data.TotalPages = data.Response.NbMangas / 18
	if data.Response.NbMangas%18 > 0 {
		data.TotalPages++
	}

	tmpl, err := template.ParseFiles(utils.Path+"templates/category.gohtml", utils.Path+"templates/header-line2.gohtml", utils.Path+"templates/base.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

// categoryNameHandlerGet
//
//	@Description: displays the mangas matching a specific Tag (Public, Status, or
//	a special request) which group and name is sent in the URL.
func categoryNameHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	group := r.PathValue("group")
	name := r.PathValue("name")
	if name == "" || group == "" {
		http.Redirect(w, r, "/error404", http.StatusNotFound)
		return
	}
	var order, pagination string
	if r.URL.Query().Has("order") {
		order = r.URL.Query().Get("order")
	} else {
		order = "desc"
	}
	if r.URL.Query().Has("pag") {
		pagination = r.URL.Query().Get("pag")
	} else {
		pagination = "1"
	}
	pag, errAtoi := strconv.Atoi(pagination)
	if errAtoi != nil {
		pag = 1
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	if pag < 1 {
		pag = 1
	}
	offset := (pag - 1) * 18

	var request models.MangaRequest

	if group == "public" && slices.Contains(models.MangaPublic, name) {
		request = models.MangaRequest{
			OrderType:  "rating",
			OrderValue: order,
			Public:     []string{name},
			Limit:      18,
			Offset:     offset,
		}
	} else if group == "status" && slices.Contains(models.MangaStatus, name) {
		request = models.MangaRequest{
			OrderType:  "rating",
			OrderValue: order,
			Status:     []string{name},
			Limit:      18,
			Offset:     offset,
		}
	} else if group == "special" {
		if name == "latest-updates" {
			request = models.MangaRequest{
				OrderType:  "latestUploadedChapter",
				OrderValue: order,
				Limit:      18,
				Offset:     offset,
			}
			name = "latest uploaded"
		} else if name == "popular" {
			request = models.MangaRequest{
				OrderType:  "rating",
				OrderValue: order,
				Limit:      18,
				Offset:     offset,
			}
		}
	} else {
		http.Redirect(w, r, "/error404", http.StatusNotFound)
		return
	}

	var data = struct {
		IsConnected bool
		Username    string
		AvatarImg   string
		Path        string
		Name        string
		Response    models.MangasInBulk
		CurrentPage int
		TotalPages  int
		Order       string
		Previous    int
		Next        int
	}{
		AvatarImg:   "avatar.jpg",
		Path:        "../../static",
		Name:        strings.ToTitle(group) + ": " + name,
		Response:    api.FetchManga(request),
		CurrentPage: pag,
		Order:       order,
		Previous:    pag - 1,
		Next:        pag + 1,
	}

	session, _ := utils.GetSession(r)
	data.Username = session.Username
	data.IsConnected = api.AddFavoriteInfo(r, &data.Response.Mangas)

	user, ok := utils.SelectUser(session.Username)
	if !ok {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
	}
	data.AvatarImg = user.Avatar

	data.TotalPages = data.Response.NbMangas / 18
	if data.Response.NbMangas%18 > 0 {
		data.TotalPages++
	}

	tmpl, err := template.ParseFiles(utils.Path+"templates/category.gohtml", utils.Path+"templates/header-line2.gohtml", utils.Path+"templates/base.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

// searchHandlerGet
//
//	@Description: displays all mangas matching search and advanced search requests
//	and the search and advanced search form.
func searchHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())

	var data struct {
		ExpandedFilters bool
		IsConnected     bool
		Username        string
		AvatarImg       string
		Tags            models.OrderedTags
		Path            string
		IsResponse      bool
		Response        models.MangasInBulk
		CurrentPage     int
		TotalPages      int
		Order           string
		Previous        int
		Next            int
		Req             string
	}

	if r.URL.Query().Has("q") {

		var pagination string
		if r.URL.Query().Has("pag") {
			pagination = r.URL.Query().Get("pag")
		} else {
			pagination = "1"
		}
		pag, errAtoi := strconv.Atoi(pagination)
		if errAtoi != nil {
			pag = 1
		}
		if pag < 1 {
			pag = 1
		}
		offset := (pag - 1) * 18

		var request models.MangaRequest

		request = models.MangaRequest{
			OrderType:      "rating",
			OrderValue:     r.URL.Query().Get("order"),
			IncludedTags:   r.URL.Query()["includedTags[]"],
			ExcludedTags:   r.URL.Query()["excludedTags[]"],
			Title:          r.URL.Query().Get("title"),
			Author:         r.URL.Query().Get("author"),
			AuthorOrArtist: r.URL.Query().Get("authorOrArtist"),
			Status:         r.URL.Query()["status[]"],
			Public:         r.URL.Query()["public[]"],
			Limit:          18,
			Offset:         offset,
		}
		log.Printf("request: %#v\n", request)

		if request.OrderValue != "asc" && request.OrderValue != "desc" {
			request.OrderValue = "desc"
		}

		var query string
		query += "?q=Search"
		query += "&title=" + request.Title
		query += "&author=" + request.Author
		query += "&authorOrArtist=" + request.AuthorOrArtist
		for _, tag := range request.IncludedTags {
			query += "&includedTags[]=" + tag
		}
		for _, tag := range request.ExcludedTags {
			query += "&excludedTags[]=" + tag
		}
		for _, status := range request.Status {
			query += "&status[]=" + status
		}
		for _, public := range request.Public {
			query += "&public[]=" + public
		}
		query += "&order[" + request.OrderType + "]=" + request.OrderValue

		data = struct {
			ExpandedFilters bool
			IsConnected     bool
			Username        string
			AvatarImg       string
			Tags            models.OrderedTags
			Path            string
			IsResponse      bool
			Response        models.MangasInBulk
			CurrentPage     int
			TotalPages      int
			Order           string
			Previous        int
			Next            int
			Req             string
		}{
			ExpandedFilters: r.URL.Query().Has("option"),
			IsConnected:     data.IsConnected,
			Username:        data.Username,
			AvatarImg:       "avatar.jpg",
			Tags:            api.FetchSortedTags(),
			Path:            "../static",
			Response:        api.FetchManga(request),
			CurrentPage:     pag,
			Order:           request.OrderValue,
			Previous:        pag - 1,
			Next:            pag + 1,
			Req:             query,
		}

		session, _ := utils.GetSession(r)
		data.Username = session.Username
		data.IsConnected = api.AddFavoriteInfo(r, &data.Response.Mangas)

		user, ok := utils.SelectUser(session.Username)
		if !ok {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
		}
		data.AvatarImg = user.Avatar

		data.IsResponse = data.Response.Mangas != nil
		data.TotalPages = data.Response.NbMangas / 18
		if data.Response.NbMangas%18 > 0 {
			data.TotalPages++
		}
	} else {
		data = struct {
			ExpandedFilters bool
			IsConnected     bool
			Username        string
			AvatarImg       string
			Tags            models.OrderedTags
			Path            string
			IsResponse      bool
			Response        models.MangasInBulk
			CurrentPage     int
			TotalPages      int
			Order           string
			Previous        int
			Next            int
			Req             string
		}{
			ExpandedFilters: r.URL.Query().Has("option"),
			IsConnected:     data.IsConnected,
			Username:        data.Username,
			AvatarImg:       "avatar.jpg",
			Tags:            api.FetchSortedTags(),
			Path:            "../static",
			IsResponse:      false,
			Response:        models.MangasInBulk{},
			CurrentPage:     1,
			TotalPages:      1,
			Order:           "desc",
			Previous:        1,
			Next:            1,
			Req:             "",
		}

		session, _ := utils.GetSession(r)
		data.Username = session.Username
		data.IsConnected = api.AddFavoriteInfo(r, &data.Response.Mangas)

		user, ok := utils.SelectUser(session.Username)
		if !ok {
			utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", errors.New("user not found")))
		}
		data.AvatarImg = user.Avatar

	}

	tmpl, err := template.ParseFiles(utils.Path+"templates/search.gohtml", utils.Path+"templates/header-line2.gohtml", utils.Path+"templates/base.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}
