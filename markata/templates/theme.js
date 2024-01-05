function setTheme(theme) {
    document.documentElement.setAttribute("data-theme", theme);
}

function detectColorSchemeOnLoad() {
    //local storage is used to override OS theme settings
    if (localStorage.getItem("theme")) {
        if (localStorage.getItem("theme") == "dark") {
            setTheme("dark");
        } else if (localStorage.getItem("theme") == "light") {
            setTheme("light");
        }
    } else if (!window.matchMedia) {
        //matchMedia method not supported
        setTheme("light");
        return false;
    } else if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
        //OS theme setting detected as dark
        setTheme("dark");
    } else {
        setTheme("light");
    }
}
detectColorSchemeOnLoad();
    document.addEventListener(
    "DOMContentLoaded",
    function () {
        //identify the toggle switch HTML element
        const toggleSwitch = document.querySelector(
            '#theme-switch input[type="checkbox"]',
        );

        //function that changes the theme, and sets a localStorage variable to track the theme between page loads
        function switchTheme(e) {
            if (e.target.checked) {
                localStorage.setItem("theme", "dark");
                document.documentElement.setAttribute("data-theme", "dark");
                toggleSwitch.checked = true;
            } else {
                localStorage.setItem("theme", "light");
                document.documentElement.setAttribute("data-theme", "light");
                toggleSwitch.checked = false;
            }
        }

        //listener for changing themes
        toggleSwitch.addEventListener("change", switchTheme, false);

        //pre-check the dark-theme checkbox if dark-theme is set
        if (document.documentElement.getAttribute("data-theme") == "dark") {
            toggleSwitch.checked = true;
        }
    },
    false,
);
