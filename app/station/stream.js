"use strict";

function stream(url, format, startTime)
{
    let progressBar        = document.getElementById("progressBar");
    let progressBarTooltip = document.getElementById("progressBarTooltip");
    let volumeTooltip      = document.getElementById("volumeSliderTooltip");

    try
    {
        let audio         = new Audio(url);
        audio.currentTime = startTime;
        audio.type        = format;

        audio.addEventListener("loadstart", () => audio.currentTime = startTime);
        audio.addEventListener("canplay", audio.play);

        audio.addEventListener("timeupdate", () =>
        {
            progressBar.style.transform    = "scaleX(" + (audio.currentTime / audio.duration) + ")";
            progressBarTooltip.textContent = audio.currentTime.toString();
        });

        audio.addEventListener("ended", songEnd);

        audio.addEventListener("error", () =>
            window.alert("Error streaming song (" + audio.error.code + "): " + audio.error.message));

        document.getElementById("volumeMuteToggle")
                .addEventListener("change", () => audio.muted = this.checked);

        document.getElementById("volumeSlider")
                .addEventListener("change", () =>
                {
                    let val                   = this.value / this.max;
                    audio.volume              = val;
                    volumeTooltip.textContent = Math.floor(val) + "%";
                });
    }
    catch(err)
    {
        window.alert("Your browser does not support HTML5 audio!\nGet a better browser (eg: Chrome or Firefox)!");
    }
}

function songEnd()
{
    // TODO tell server we finished the playing song
    // When server gets the first "finished" from any listener to a station,
    // start a timeout to receive "finished" messages from the rest of the station's listeners.
    // This way we can keep listeners in (relative) sync while keeping the listening experience
    // not terrible (ie: syncing mid song).
    // When the timeout is up, send the next song to all listeners.
    // (server can then check here for lost connections)
}
