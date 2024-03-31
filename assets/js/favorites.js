"use strict"

let addFavoriteBtns = document.querySelectorAll('.add-favorite');
let deleteFavoriteBtns = document.querySelectorAll('.delete-favorite');
let setBannerBtns = document.querySelectorAll('.set-banner');

async function sendRequest(method, manga) {
    const response = await fetch(`http://localhost:8080/favorite/${manga}`, {
        method: method, // *GET, POST, PUT, DELETE, etc.
        cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
        credentials: "same-origin", // include, *same-origin, omit
        redirect: "follow", // manual, *follow, error
        referrerPolicy: "no-referrer", // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
    });
    return response.ok;
}

function addToFavorites(e) {
    let Res = sendRequest('POST', e.currentTarget.id);
    if (Res) {
        console.log(`Manga ${e.currentTarget.id} has been added to your favorites!`);
        e.currentTarget.classList.remove('add-favorite');
        e.currentTarget.classList.add('delete-favorite');
        e.currentTarget.removeEventListener('click', addToFavorites);
        e.currentTarget.addEventListener('click', deleteFromFavorites);
        let imgFav = e.currentTarget.querySelector('img');
        imgFav.src = imgFav.src.split('/').slice(0,-1).join('/')+'/darkred-remove-favorite.png';
    } else {
        console.log(`An error occurred!`);
    }
}

for (let addFavoriteBtn of addFavoriteBtns) {
    addFavoriteBtn.addEventListener('click', addToFavorites);
}

function deleteFromFavorites(e) {
    let Res = sendRequest('DELETE', e.currentTarget.id);
    if (Res) {
        console.log(`Manga ${e.currentTarget.id} has been removed from your favorites!`);
        e.currentTarget.classList.remove('delete-favorite');
        e.currentTarget.classList.add('add-favorite');
        e.currentTarget.removeEventListener('click', deleteFromFavorites);
        e.currentTarget.addEventListener('click', addToFavorites);
        let imgFav = e.currentTarget.querySelector('img');
        imgFav.src = imgFav.src.split('/').slice(0,-1).join('/')+'/darkred-add-favorite.png';
    } else {
        console.log(`An error occurred!`);
    }
}

for (let deleteFavoriteBtn of deleteFavoriteBtns) {
    deleteFavoriteBtn.addEventListener('click', deleteFromFavorites);
}

function setBanner(e) {
    let Res = sendRequest('PUT', e.currentTarget.id);
    if (Res) {
        console.log(`Manga ${e.currentTarget.id} has been set as your banner!`);
        e.currentTarget.classList.remove('set-banner');
        e.currentTarget.removeEventListener('click', setBanner);
        location.reload();
    } else {
        console.log(`An error occurred!`);
    }
}

for (let bannerBtn of setBannerBtns) {
    bannerBtn.addEventListener('click', setBanner);
}