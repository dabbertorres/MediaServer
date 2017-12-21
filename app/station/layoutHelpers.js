"use strict";

function listItem(lines, ...content)
{
    let el = document.createElement("li");
    el.className = "mdl-list__item";

    switch(lines)
    {
        case 1:
            break;

        case 2:
            el.className += " mdl-list__item--two-line";
            break;

        case 3:
            el.className += " mdl-list__item--three-line";
            break;

        default:
            return null;
    }

    let primary = document.createElement("span");
    primary.className = "mdl-list__item-primary-content";

    for(let c of content)
        primary.appendChild(c);

    el.appendChild(primary);

    return el;
}

function listItemTitle(text)
{
    let el = document.createElement("span");
    el.textContent = text;
    return el;
}

function listItemSubTitle(text)
{
    let el = document.createElement("span");
    el.className = "mdl-list__item-sub-title";
    el.textContent = text;
    return el;
}

function listItemTextBody(text)
{
    let el = document.createElement("p");
    el.className = "mdl-list__item-text-body";
    el.textContent = text;
    return el;
}
