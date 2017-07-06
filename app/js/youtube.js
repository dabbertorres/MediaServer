// provides a function to load a specified youtube video

"use strict";

let youtubePlayer = null;

// only function our stuff calls!
function playYouTube(url)
{
    if(youtubePlayer !== null)
        youtubePlayer.loadVideoByUrl(url);
    else
        setTimeout(playYouTube, 50, url);
}

// called by the YouTube API script
function onYouTubeIframeAPIReady()
{
    const minPlayerSize = 200;

    // values from docs: https://developers.google.com/youtube/player_parameters#Parameters
    youtubePlayer = new YT.Player("videoPlayer",
    {
        height: minPlayerSize,
        width:  minPlayerSize,

        cc_load_policy: 0,
        controls:       0,
        disablekb:      0,
        enablejsapi:    1,
        fs:             0,
        iv_load_policy: 3,
        modestbranding: 1,
        origin:         "", // TODO fill in with domain
        rel:            0,
        showinfo:       0,

        events:
        {
            "onReady":       onReady,
            "onStateChange": onStateChange,
            "onError":       onError,
        },
    });
}

function onReady(event)
{
    let toggle = document.getElementById("volumeMuteToggle");
    let volume = document.getElementById("volumeSlider");

    let checkMute = function()
    {
        if(toggle.checked)
            youtubePlayer.mute();
        else
            youtubePlayer.unMute();
    };

    let volumeChange = () => youtubePlayer.setVolume(volume.value);

    toggle.addEventListener("change", checkMute);
    volume.addEventListener("change", volumeChange);

    checkMute();
    volumeChange();

    event.target.playVideo();
}

function onStateChange(event)
{
    // YT does not have a constant for the unstarted state: https://developers.google.com/youtube/iframe_api_reference#onStateChange
    const UNSTARTED = -1;

    switch(event.data)
    {
        case UNSTARTED:
            break;

        case YT.PlayerState.ENDED:
            songEnd();
            break;

        case YT.PlayerState.PLAYING:
            break;

        case YT.PlayerState.PAUSED:
            break;

        case YT.PlayerState.BUFFERING:
            break;

        case YT.PlayerState.CUED:
            break;

        default:
            break;
    }
}

function onError(event)
{
    // defined in YouTube iframe api docs: https://developers.google.com/youtube/iframe_api_reference#onError
    const INVALID_PARAM     = 2;
    const CANNOT_PLAY_HTML5 = 5;
    const NOT_FOUND         = 100;
    const NO_EMBEDDED       = 101;
    const NO_EMBEDDED_2     = 150;

    switch(event.data)
    {
        case INVALID_PARAM:
            console.error("invalid error");
            break;

        case CANNOT_PLAY_HTML5:
            console.error("html5 player error");
            break;

        case NOT_FOUND:
            console.error("video not found");
            break;

        case NO_EMBEDDED:
        case NO_EMBEDDED_2:
            console.error("video not allowed to be played embedded");
            break;

        default:
            console.log("unknown value of: '" + event.data + "' received in onError");
    }
}
