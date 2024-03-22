"use strict"

let avatars = document.querySelectorAll('.avatar-choice');

function selectAvatar(e) {
    let choice = document.getElementById(`${e.currentTarget.dataset.idValue}`);
    for (let avatar of avatars) {
        avatar.classList.remove('selected');
    }
    e.currentTarget.classList.add('selected');
    choice.checked = true;
}

for (let avatar of avatars) {
    avatar.addEventListener('click', selectAvatar);
}