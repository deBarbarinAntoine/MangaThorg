package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"log/slog"
	"mangathorg/internal/api"
	"mangathorg/internal/middlewares"
	"mangathorg/internal/models"
	"mangathorg/internal/utils"
	"net/http"
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

	tmpl, err := template.ParseFiles(utils.Path+"templates/base.gohtml", utils.Path+"templates/error404.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func indexHandlerOther(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	log.Println("HTTP ApiErr", http.StatusNotFound)
	w.WriteHeader(http.StatusNotFound)
	utils.Logger.Warn("indexHandlerOther", slog.Int("req_id", middlewares.LogId), slog.String("req_url", r.URL.String()), slog.Int("http_status", http.StatusNotFound))

	tmpl, err := template.ParseFiles(utils.Path+"templates/base.gohtml", utils.Path+"templates/error404.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "base", nil)
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
	}
	tmpl, err := template.ParseFiles(utils.Path + "templates/login.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "login", message)
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
	tmpl, err := template.ParseFiles(utils.Path + "templates/register.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "register", message)
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
	utils.SendMail(&newTempUser)
	utils.TempUsers = append(utils.TempUsers, newTempUser)
	log.Printf("newTempUser: %#v\n", newTempUser)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
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
	tmpl, err := template.ParseFiles(utils.Path+"templates/principal.gohtml", utils.Path+"templates/base.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	var data = struct {
		Banner         models.MangaUsefullData
		LatestUploaded []models.MangaUsefullData
		Popular        []models.MangaUsefullData
	}{
		Banner:         api.FetchMangaById("cb676e05-8e6e-4ec4-8ba0-d3cb4f033cfa", "asc", 1),
		LatestUploaded: api.FetchManga(api.TopLatestUploadedRequest).Mangas,
		Popular:        api.FetchManga(api.TopPopularRequest).Mangas,
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
	tmpl, err := template.ParseFiles(utils.Path+"templates/manga.gohtml", utils.Path+"templates/base.gohtml")
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
	log.Println("pages:", data.Pages)
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

func chapterHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	mangaId := r.PathValue("mangaId")
	chapterNb := r.PathValue("chapterNb")
	chapterId := r.PathValue("chapterId")
	if chapterId == "" || mangaId == "" || chapterNb == "" {
		http.Redirect(w, r, "/error404", http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles(utils.Path+"templates/chapter.gohtml", utils.Path+"templates/base.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	scan := api.ScanRequest(chapterId)
	var data = struct {
		Manga     string
		ChapterNb string
		Id        string
		Quality   string
		Alt       string
		Scan      struct {
			Hash      string
			Data      []string
			DataSaver []string
		}
	}{
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
	data.Alt = data.Manga + " - Ch. " + data.ChapterNb
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}

func tagsHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	tmpl, err := template.ParseFiles(utils.Path+"templates/tags.gohtml", utils.Path+"templates/base.gohtml")
	if err != nil {
		log.Fatalln(err)
	}

	err = tmpl.ExecuteTemplate(w, "base", api.TagsRequest().Data)
	if err != nil {
		log.Fatalln(err)
	}
}

func categoryHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	tmpl, err := template.ParseFiles(utils.Path+"templates/category.gohtml", utils.Path+"templates/base.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
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
		Tag         models.ApiTag
		Response    models.MangasInBulk
		CurrentPage int
		TotalPages  int
		Order       string
		Previous    int
		Next        int
	}{
		Tag:         api.TagSelect(tagId),
		Response:    api.FetchManga(request),
		CurrentPage: pag,
		Order:       order,
		Previous:    pag - 1,
		Next:        pag + 1,
	}
	data.TotalPages = data.Response.NbMangas / 18
	if data.Response.NbMangas%18 > 0 {
		data.TotalPages++
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Fatalln(err)
	}
}
