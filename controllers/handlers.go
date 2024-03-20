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
	"slices"
	"strconv"
	"strings"
	"time"
)

func indexHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	//tmpl, err := template.ParseFiles(utils.Path + "templates/index.gohtml")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//err = tmpl.ExecuteTemplate(w, "index", "indexHandlerGet")
	//if err != nil {
	//	log.Fatalln(err)
	//}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.StatRequest("d1a9fdeb-f713-407f-960c-8326b586e6fd"))
}

func indexHandlerPut(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	tmpl, err := template.ParseFiles(utils.Path + "templates/index.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	session, sessionID := utils.GetSession(r)
	err = tmpl.ExecuteTemplate(w, "index", "indexHandlerPut"+sessionID+"\nUsername: "+session.Username+"\nIP address: "+session.IpAddress)
	if err != nil {
		log.Fatalln(err)
	}
}

func indexHandlerDelete(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	tmpl, err := template.ParseFiles(utils.Path + "templates/index.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	session, sessionID := utils.GetSession(r)
	err = tmpl.ExecuteTemplate(w, "index", "indexHandlerPut"+sessionID+"\nUsername: "+session.Username+"\nIP address: "+session.IpAddress)
	if err != nil {
		log.Fatalln(err)
	}
}

func indexHandlerNoMeth(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	log.Println("HTTP ApiErr", http.StatusMethodNotAllowed)
	w.WriteHeader(http.StatusMethodNotAllowed)
	utils.Logger.Warn("indexHandlerNoMeth", slog.Int("req_id", middlewares.LogId), slog.String("req_url", r.URL.String()), slog.Int("http_status", http.StatusMethodNotAllowed))

	var data = struct {
		IsConnected bool
		Username    string
		AvatarImg   string
	}{
		AvatarImg: "avatar.jpg",
	}

	user, sessionId := utils.GetSession(r)
	if sessionId != "" {
		data.IsConnected = true
		data.Username = user.Username
	}

	tmpl, err := template.ParseFiles(utils.Path+"templates/base.gohtml", utils.Path+"templates/header-line2.gohtml", utils.Path+"templates/error404.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

func indexHandlerOther(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	log.Println("HTTP ApiErr", http.StatusNotFound)
	w.WriteHeader(http.StatusNotFound)
	utils.Logger.Warn("indexHandlerOther", slog.Int("req_id", middlewares.LogId), slog.String("req_url", r.URL.String()), slog.Int("http_status", http.StatusNotFound))

	var data = struct {
		IsConnected bool
		Username    string
		AvatarImg   string
	}{
		AvatarImg: "avatar.jpg",
	}

	user, sessionId := utils.GetSession(r)
	if sessionId != "" {
		data.IsConnected = true
		data.Username = user.Username
	}

	tmpl, err := template.ParseFiles(utils.Path+"templates/base.gohtml", utils.Path+"templates/header-line2.gohtml", utils.Path+"templates/error404.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

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
		if r.URL.Query().Get("status") == "update-pwd" {
			message = `<div class="message">Your password has been updated!</div>`
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
		password1: r.FormValue("password1"),
		password2: r.FormValue("password2"),
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
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

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
		http.Redirect(w, r, "update-credentials?id="+id, http.StatusSeeOther)
		return
	}
	lostUser.User.HashedPwd, lostUser.User.Salt = utils.NewPwd(password)
	utils.UpdateLostUser(lostUser)
	http.Redirect(w, r, "/login?status=update-pwd", http.StatusSeeOther)
}

func homeHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	tmpl, err := template.ParseFiles(utils.Path + "templates/index.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "index", "homeHandlerGet --- Restricted area! ---")
	if err != nil {
		log.Fatalln(err)
	}
}

func logHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Query().Has("level") {
		json.NewEncoder(w).Encode(utils.FetchAttrLogs("level", r.URL.Query().Get("level")))
		return
	} else if r.URL.Query().Has("user") {
		json.NewEncoder(w).Encode(utils.FetchAttrLogs("user", r.URL.Query().Get("user")))
		return
	}
	json.NewEncoder(w).Encode(utils.RetrieveLogs())
}

func confirmHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	if r.URL.Query().Has("id") {
		id := r.URL.Query().Get("id")
		utils.PushTempUser(id)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func logoutHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	utils.Logout(&w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

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
		AvatarImg:      "avatar.jpg",
		Banner:         api.FetchMangaById("cb676e05-8e6e-4ec4-8ba0-d3cb4f033cfa", "asc", 1),
		LatestUploaded: api.FetchManga(api.TopLatestUploadedRequest).Mangas,
		Popular:        api.FetchManga(api.TopPopularRequest).Mangas,
	}

	user, sessionId := utils.GetSession(r)
	if sessionId != "" {
		data.IsConnected = true
		data.Username = user.Username
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

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
		AvatarImg:   "avatar.jpg",
		Manga:       manga,
		CurrentPage: pag,
		Pages:       pages,
		Order:       order,
	}

	session, _ := utils.GetSession(r)
	data.Username = session.Username
	data.IsConnected = api.AddSingleFavoriteInfo(r, &data.Manga)

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
	}
}

func tagsRequestHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.TagsRequest())
}

func feedRequestHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.FeedRequest("d1a9fdeb-f713-407f-960c-8326b586e6fd", "desc", 1))
}

func chapterScanRequestHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.ScanRequest("444b113a-3705-4718-8f91-f46c640ab433"))
}

func mangaWholeRequestHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.FetchMangaById("d1a9fdeb-f713-407f-960c-8326b586e6fd", "desc", 1))
}

func coverHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	mangaId := r.PathValue("manga")
	img := r.PathValue("img")
	if mangaId == "" || img == "" {
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	_, err := w.Write(api.ImageProxy(mangaId, img))
	if err != nil {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		w.WriteHeader(http.StatusNotFound)
	}
}

func scanHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	chapterId := r.PathValue("chapterId")
	quality := r.PathValue("quality")
	hash := r.PathValue("hash")
	img := r.PathValue("img")
	if chapterId == "" || quality == "" || hash == "" || img == "" {
		log.Println("empty value")
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	i, err := w.Write(api.ScanProxy(chapterId, quality, hash, img))
	if err != nil || i == 0 {
		utils.Logger.Error(utils.GetCurrentFuncName(), slog.Any("output", err))
		w.WriteHeader(http.StatusNotFound)
	}
}

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

func chapterHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	mangaId := r.PathValue("mangaId")
	chapterNb := r.PathValue("chapterNb")
	chapterId := r.PathValue("chapterId")
	if chapterId == "" || mangaId == "" || chapterNb == "" {
		http.Redirect(w, r, "/error404", http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles(utils.Path+"templates/chapter.gohtml", utils.Path+"templates/header-line2.gohtml", utils.Path+"templates/base.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	scan := api.ScanRequest(chapterId)
	var data = struct {
		IsConnected bool
		Username    string
		AvatarImg   string
		IsOk        bool
		ToDataSaver bool
		Manga       string
		ChapterNb   string
		Id          string
		Quality     string
		Alt         string
		Scan        struct {
			Hash      string
			Data      []string
			DataSaver []string
		}
	}{
		AvatarImg: "avatar.jpg",
		Manga:     api.FetchMangaById(mangaId, "desc", 1).Title,
		ChapterNb: chapterNb,
		Id:        chapterId,
		Quality:   "data",
		Alt:       "",
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

	user, sessionId := utils.GetSession(r)
	if sessionId != "" {
		data.IsConnected = true
		data.Username = user.Username
	}

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
		AvatarImg:  "avatar.jpg",
		FormatTags: sortedTags.FormatTags,
		GenreTags:  sortedTags.GenreTags,
		ThemeTags:  sortedTags.ThemeTags,
		PublicTags: sortedTags.PublicTags,
		StatusTags: sortedTags.StatusTags,
	}

	user, sessionId := utils.GetSession(r)
	if sessionId != "" {
		data.IsConnected = true
		data.Username = user.Username
	}

	tmpl, err := template.ParseFiles(utils.Path+"templates/tags.gohtml", utils.Path+"templates/header-line2.gohtml", utils.Path+"templates/base.gohtml")
	if err != nil {
		log.Fatalln(err)
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

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

	user, _ := utils.GetSession(r)
	data.Username = user.Username
	data.IsConnected = api.AddFavoriteInfo(r, &data.Response.Mangas)

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

	user, _ := utils.GetSession(r)
	data.Username = user.Username
	data.IsConnected = api.AddFavoriteInfo(r, &data.Response.Mangas)

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

// searchHandlerGet is the handler that serves all search requests and the advanced search page.
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

	user, _ := utils.GetSession(r)
	data.Username = user.Username
	data.IsConnected = api.AddFavoriteInfo(r, &data.Response.Mangas)

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
