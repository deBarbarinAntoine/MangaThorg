"use strict"

let tags = document.querySelectorAll('.search-tag');
let simpleTags = document.querySelectorAll('.search-tag-simple');
let checkboxes = document.querySelectorAll('.three-way-checkbox');
let simpleCheckboxes = document.querySelectorAll('.two-way-checkbox');
let resetBtn = document.querySelector('.reset-btn');

function ToggleCheckboxes(tag, status) {
    let id = tag.id;
    switch (status) {
        case 'selected':
            for (let checkbox of checkboxes) {
                if (checkbox.attributes.getNamedItem('value') === id) {
                    checkbox.toggleAttribute('checked');
                }
            }
            return;
        case 'unselected':
            for (let checkbox of checkboxes) {
                if (checkbox.getAttribute('value') === id && checkbox.getAttribute('name') === 'excludedTags[]') {
                    checkbox.toggleAttribute('checked');
                }
            }
            return;
        case 'none':
            for (let checkbox of checkboxes) {
                if (checkbox.getAttribute('value') === id && checkbox.getAttribute('name') === 'includedTags[]') {
                    checkbox.toggleAttribute('checked');
                }
            }
            return;
    }
}

function ClickSimpleTag(e) {
    e.currentTarget.classList.toggle('selected');
    for (let checkbox of simpleCheckboxes) {
        if (checkbox.getAttribute('name') === e.currentTarget.id) {
            checkbox.toggleAttribute('checked');
        }
    }
}

function ClickTag(e) {
    if (e.currentTarget.classList.contains('selected')) {
        e.currentTarget.classList.remove('selected');
        e.currentTarget.classList.add('unselected');
        ToggleCheckboxes(e.currentTarget, 'selected');
    } else if (e.currentTarget.classList.contains('unselected')) {
        e.currentTarget.classList.remove('unselected');
        ToggleCheckboxes(e.currentTarget, 'unselected');
    } else {
        e.currentTarget.classList.add('selected');
        ToggleCheckboxes(e.currentTarget, 'none');
    }
}

function Reset() {
    for (let tag of tags) {
        tag.classList.remove('selected', 'unselected');
    }
    for (let tag of simpleTags) {
        tag.classList.remove('selected', 'unselected');
    }
}

for (let tag of tags) {
    tag.addEventListener('click', ClickTag);
}

for (let simpleTag of simpleTags) {
    simpleTag.addEventListener('click', ClickSimpleTag);
}

resetBtn.addEventListener('click', Reset);