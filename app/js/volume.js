// manages the UI on volume state changes
"use strict";

const icons =
    {
        on: "volume_up",
        off: "volume_off"
    };

let muteToggle = document.getElementById("volumeMuteToggle");
let iconDisplay = document.getElementById("volumeIcon");

muteToggle.addEventListener("change", function()
{
    if(muteToggle.checked)
        iconDisplay.textContent = icons.off;
    else
        iconDisplay.textContent = icons.on;
});
