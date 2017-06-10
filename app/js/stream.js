"use strict";

let setIcon = function(muted)
{
    const on  = "volume_up";
    const off = "volume_off";

    document.getElementById("volumeIcon")
        .textContent = muted ? on : off;
};

function stream(url, format, startTime)
{
    let progressBar = document.getElementById("progressBar");

    try
    {
        let audio         = new Audio(url);
        audio.currentTime = startTime;
        audio.type        = format;

        audio.addEventListener("loadstart", () => audio.currentTime = startTime);
        audio.addEventListener("canplay", audio.play);

        audio.addEventListener("timeupdate", function()
        {
            progressBar.style.transform = "scaleX(" + (this.currentTime / this.duration) + ")";
        });

        audio.addEventListener("ended", function()
        {
           // send message that this client finished
        });

        audio.addEventListener("error",
            () => window.alert("Error streaming song: " + audio.error.message + " (" + audio.error.code + ")"));

        document.getElementById("volumeMuteToggle")
            .addEventListener("change", function()
            {
                setIcon(this.checked);
                audio.muted = this.checked;
            });

        document.getElementById("volumeSlider")
            .addEventListener("change", function(){ audio.volume = this.value / this.max; });
    }
    catch(err)
    {
        window.alert("Your browser does not support HTML5 audio!\nGet a better browser (ie: Chrome or Firefox)!");
    }
}
