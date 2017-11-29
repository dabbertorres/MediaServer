// this script will handle the WebSocket connection, and its communications

"use strict";

// websrv/msg/messaage.go
const MESSAGE_CHAT     = 0;
const MESSAGE_STREAM   = 1;
const MESSAGE_PLAYLIST = 2;
const MESSAGE_STATUS   = 3;

const SOCKET_URL = "ws://" + window.location.host + window.location.pathname.replace("station", "socket");

// elements

let icon            = document.getElementById("volumeIcon");
let toggle          = document.getElementById("volumeMuteToggle");
let chatMessageList = document.getElementById("chat-message-list");
let playlist        = document.getElementById("playlist");

// element events/actions

toggle.addEventListener("change", () => icon.textContent = toggle.checked ? "volume_off" : "volume_up");

// get user information

// TODO make nicer
let username = window.prompt("Name?");

// setup server socket connection

let socket = new WebSocket(SOCKET_URL);

if(socket !== null)
    document.addEventListener("unload", socket.close);

// TODO tell the user about the error
socket.addEventListener("error", err => console.error(err));
socket.addEventListener("message", socketMessage);

function chatInput(text)
{
    let msg = JSON.stringify(
        {
            type: MESSAGE_CHAT,
            chat:
            {
                from:    username,
                content: text,
            },
        });

    socket.send(msg);
}

function socketMessage(ev)
{
    let msg = JSON.parse(ev.data);

    switch(msg.type)
    {
        case MESSAGE_CHAT:
            chatMessageList.appendChild(
                listItem(2,
                    listItemTitle(msg.chat.from + ":"),
                    listItemTextBody(msg.chat.content)));
            break;

        case MESSAGE_STREAM:
            break;

        case MESSAGE_PLAYLIST:
            for(let song of msg.playlist.update)
            {
                playlist.appendChild(
                    listItem(2,
                        listItemTitle(song.title),
                        listItemSubTitle(song.artist)));
            }
            break;

        case MESSAGE_STATUS:
            break;

        default:
            console.warn("wtf, got a message with a weird id: " + msg.id);
            break;
    }
}
