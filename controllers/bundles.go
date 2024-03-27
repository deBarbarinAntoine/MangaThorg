package controllers

import (
	"mangathorg/internal/middlewares"
)

// Generic Bundles (root and errors)

var RootHandlerGetBundle = middlewares.Join(rootHandlerGet, middlewares.Log, middlewares.UserCheck)
var ErrorHandlerBundle = middlewares.Join(errorHandler, middlewares.Log, middlewares.UserCheck)

// Only visitors Bundles

var LoginHandlerGetBundle = middlewares.Join(loginHandlerGet, middlewares.Log, middlewares.OnlyVisitors)
var LoginHandlerPostBundle = middlewares.Join(loginHandlerPost, middlewares.Log, middlewares.OnlyVisitors)

var RegisterHandlerGetBundle = middlewares.Join(registerHandlerGet, middlewares.Log, middlewares.OnlyVisitors)
var RegisterHandlerPostBundle = middlewares.Join(registerHandlerPost, middlewares.Log, middlewares.OnlyVisitors)

var ForgotPasswordHandlerGetBundle = middlewares.Join(forgotPasswordHandlerGet, middlewares.Log, middlewares.OnlyVisitors)
var ForgotPasswordHandlerPostBundle = middlewares.Join(forgotPasswordHandlerPost, middlewares.Log, middlewares.OnlyVisitors)

var UpdateCredentialsHandlerGetBundle = middlewares.Join(updateCredentialsHandlerGet, middlewares.Log, middlewares.OnlyVisitors)
var UpdateCredentialsHandlerPostBundle = middlewares.Join(updateCredentialsHandlerPost, middlewares.Log, middlewares.OnlyVisitors)

var ConfirmHandlerGetBundle = middlewares.Join(confirmHandlerGet, middlewares.Log, middlewares.OnlyVisitors)

// Only users Bundles

var ProfileHandlerGetBundle = middlewares.Join(profileHandlerGet, middlewares.Log, middlewares.Guard)
var ProfileHandlerPostBundle = middlewares.Join(profileHandlerPost, middlewares.Log, middlewares.Guard)

var LogoutHandlerGetBundle = middlewares.Join(logoutHandlerGet, middlewares.Log, middlewares.Guard)

var HomeHandlerGetBundle = middlewares.Join(homeHandlerGet, middlewares.Log, middlewares.Guard, middlewares.CheckApi) // API needed for this Bundle

// Image request Bundles

var CoversHandlerGetBundle = middlewares.Join(coverHandlerGet, middlewares.UserCheck)
var ScanHandlerGetBundle = middlewares.Join(scanHandlerGet, middlewares.UserCheck)

// User favorites' requests (accessed from javascript requests)

var FavoriteHandlerPostBundle = middlewares.Join(favoriteHandlerPost, middlewares.SimpleGuard)
var FavoriteHandlerDeleteBundle = middlewares.Join(favoriteHandlerDelete, middlewares.SimpleGuard)

var BannerHandlerPutBundle = middlewares.Join(bannerHandlerPut, middlewares.SimpleGuard)

// Bundles available for any clients: they all need MangaDex API to work

var AboutHandlerGetBundle = middlewares.Join(aboutHandlerGet, middlewares.Log, middlewares.UserCheck)
var PrincipalHandlerGetBundle = middlewares.Join(principalHandlerGet, middlewares.Log, middlewares.UserCheck, middlewares.CheckApi)
var MangaRequestHandlerGet = middlewares.Join(mangaHandlerGet, middlewares.Log, middlewares.UserCheck, middlewares.CheckApi)
var TagsHandlerGetBundle = middlewares.Join(tagsHandlerGet, middlewares.Log, middlewares.UserCheck, middlewares.CheckApi)
var CategoryHandlerGetBundle = middlewares.Join(categoryHandlerGet, middlewares.Log, middlewares.UserCheck, middlewares.CheckApi)
var CategoryNameHandlerGetBundle = middlewares.Join(categoryNameHandlerGet, middlewares.Log, middlewares.UserCheck, middlewares.CheckApi)
var SearchHandlerGetBundle = middlewares.Join(searchHandlerGet, middlewares.Log, middlewares.UserCheck, middlewares.CheckApi)
var ChapterHandlerGetBundle = middlewares.Join(chapterHandlerGet, middlewares.Log, middlewares.UserCheck, middlewares.CheckApi)

// LogHandlerGetBundle is a special Bundle that enables access to the logs: this is a testing and developping tool only (remove before deploying)
var LogHandlerGetBundle = middlewares.Join(logHandlerGet, middlewares.Log, middlewares.UserCheck)
