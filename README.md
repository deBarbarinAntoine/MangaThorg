# Groupie Tracker – MangaDex API – Antoine de Barbarin

[MangaDex API documentation](https://api.mangadex.org/docs/).

<div style="height: 3px; background-color: #EEEEEE; border-radius: 2px"></div>

## Introduction

This project is an assignment done for my studies in my first year of computer science. It is a website based on an API named MangaDex. It permits to search, consult and read mangas.

All resources are available for all users (search and advanced search, manga page and chapter page). Only the favorites feature are for registered users only.

The mangas available in the website are all in english (title, description, scans) and SFW (contentRating[]=safe).

<div style="height: 3px; background-color: #EEEEEE; border-radius: 2px"></div>

## Installation

To set the server up on your own computer, you can use one of those two options:

---

### + From the Release version (Windows users only)

Click on ``MangaThorg v1.0`` below the `Release` option and then on the `.exe` file to download it.

Then you can just execute the downloaded executable.

---

### + From the source code

**REQUIREMENT: you need Golang v1.22 or superior to be able to run this source code.**

Click on ``<> Code`` and then on ``Download ZIP``.

Unzip the downloaded file.

#### WARNING: you need a config folder to run the program properly

- Create a folder named ``config`` at the root of the project (at the same level as `assets`, `cache` or `cmd`).
- Create a JSON file named ``config.json`` inside the new `config` folder.
- In ``config.json``, you need to put an object of this type:
````json
{
  "email_addr": "<your-email-address>",
  "email_auth": "<your-email-password>",
  "host": "<your-email-smtp-hostname>",
  "port": 465 // <your-email-smtp-port>
}
````

Fill the mail parameters in ``config.json``, otherwise, the program will bug whenever it needs to send a mail (register and forgot password options).

<div style="height: 3px; background-color: #EEEEEE; border-radius: 2px"></div>

## Routes

- **GET /{$}**: root (redirects to **/principal**).


- **GET /login**: displays the login form.
- **POST /login**: login treatment (no display).


- **GET /register**: displays the register form.
- **POST /register**: register treatment (no display).


- **GET /forgot-password**: displays the forgot password form.
- **POST /forgot-password**: forgot password treatment (no display).


- **GET /update-credentials**: displays the update credentials form (takes the id sent by mail as a GET query: ``?id={id}``.
- **POST /update-credentials/{id}**: update credentials treatment (no display). Takes the id sent by mail as the ``{id}``.


- **GET /profile**: displays the profile form (to modify the user's avatar or/and password).
- **POST /profile**: profile treatment (no display).


- **GET /home**: displays the home page (user only) with all his favorites.
- **GET /confirm**: displays the mail confirm page (corresponds to the link sent to confirm the account just created).
- **GET /logout**: log the user out (no display and user only).
- **GET /principal**: displays the principal page.
- **GET /manga/{id}**: displays the manga page according to the id specified in the URL.
- **GET /categories**: displays the categories page.
- **GET /category/{tagId}**: displays the category page according to a specific tag (which id is specified in the URL).
- **GET /category/{group}/{name}**: displays the category page according to a specific name in a specific group (used for public and format tags).
- **GET /search**: displays the search page and the results according to the query params.
- **GET /chapter/{mangaId}/{offset}/{chapterId}**: displays a chapter according to a mangaId, an offset and a chapterId.


- **GET /covers/{manga}/{img}**: used to display images in the pages (cover image proxy).
- **GET /scan/{chapterId}/{quality}/{hash}/{img}**: used to display scan images in the pages (scan image proxy).
- **POST /favorite/{mangaId}**: adds a favorite to a user (no display and user only).
- **DELETE /favorite/{mangaId}**: removes a favorite from a user (no display and user only).
- **PUT /favorite/{mangaId}**: modifies the custom user banner (no display and user only).


- **GET /logs**: (*Testing handler*) sends the logs in JSON format. Accepts a filter with ``?level={info, warn, error}`` (one of the three, optional).

<div style="height: 3px; background-color: #EEEEEE; border-radius: 2px"></div>

## Endpoints

### Search endpoint:
https://api.mangadex.org/manga

![alt text](https://github.com/deBarbarinAntoine/Livrables-projet-groupie-tracker/blob/main/img_readme/search-endpoint.png?raw=true)
 
With this endpoint, any kind of advanced search is possible:
-	by title (title= :title)
-	by including tags (includedTags[]= :ids)
-	by excluding tags (excludedTags[]= :ids)
-	by status (status[]= :status)
-	by targeted public (publicationDemographic[]= :targetType)

---

### get manga’s cover endpoint:
https://api.mangadex.org/cover/:cover-id
 
![alt text](https://github.com/deBarbarinAntoine/Livrables-projet-groupie-tracker/blob/main/img_readme/cover-endpoint.png?raw=true)

---

### get chapter id list endpoint:
https://api.mangadex.org/manga/:manga-id/feed
 
![alt text](https://github.com/deBarbarinAntoine/Livrables-projet-groupie-tracker/blob/main/img_readme/list-chapter-id-endpoint.png?raw=true)

---

### get chapter images endpoint:
https://api.mangadex.org/at-home/server/:chapter-id

![alt text](https://github.com/deBarbarinAntoine/Livrables-projet-groupie-tracker/blob/main/img_readme/chapter-imgs-endpoint.png?raw=true)
 
---

### get tag list enpoint:
https://api.mangadex.org/manga/tag
 
![alt text](https://github.com/deBarbarinAntoine/Livrables-projet-groupie-tracker/blob/main/img_readme/tag-endpoint.png?raw=true)

---

## Category page:
Tags, status and targeted public will be available categories.

---

## Favorites:
The favorites will be a personal list of mangas’ ids and/or chapters’ ids.
