# Groupie Tracker – MangaDex API – Antoine de Barbarin

[MangaDex API documentation](https://api.mangadex.org/docs/).

---

## Endpoints

### Search endpoint:
https://api.mangadex.org/manga

![alt text](https://github.com/deBarbarinAntoine/Livrables-projet-groupie-tracker/blob/main/assets/img/readme/search-endpoint.png?raw=true)
 
With this endpoint, any kind of advanced search is possible:
-	by title (title= :title)
-	by including tags (includedTags[]= :ids)
-	by excluding tags (excludedTags[]= :ids)
-	by status (status[]= :status)
-	by targeted public (publicationDemographic[]= :targetType)

---

### get manga’s cover endpoint:
https://api.mangadex.org/cover/:cover-id
 
![alt text](https://github.com/deBarbarinAntoine/Livrables-projet-groupie-tracker/blob/main/assets/img/readme/cover-endpoint.png?raw=true)

---

### get chapter id list endpoint:
https://api.mangadex.org/manga/:manga-id/feed
 
![alt text](https://github.com/deBarbarinAntoine/Livrables-projet-groupie-tracker/blob/main/assets/img/readme/list-chapter-id-endpoint.png?raw=true)

---

### get chapter images endpoint:
https://api.mangadex.org/at-home/server/:chapter-id

![alt text](https://github.com/deBarbarinAntoine/Livrables-projet-groupie-tracker/blob/main/assets/img/readme/chapter-imgs-endpoint.png?raw=true)
 
---

### get tag list enpoint:
https://api.mangadex.org/manga/tag
 
![alt text](https://github.com/deBarbarinAntoine/Livrables-projet-groupie-tracker/blob/main/assets/img/readme/tag-endpoint.png?raw=true)

---

## Category page:
Tags, status and targeted public will be available categories.

---

## Favorites:
The favorites will be a personal list of mangas’ ids and/or chapters’ ids.
