// Clicking an element with data-toggle="<selector>" shows or hides the element matching the selector.
document.addEventListener("click", (event) => {
    const trigger = event.target.closest("[data-toggle]");
    if (trigger) {
        document.querySelector(trigger.dataset.toggle)?.classList.toggle("hidden");
    }
});
