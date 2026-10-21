// Selecting a file in an <input type="file"> shows it in the <img data-image-preview> of the
// same <label> (see imageInput in web/views/meals.templ).
document.addEventListener("change", (event) => {
    const input = event.target;
    const preview = input.closest?.("label")?.querySelector("img[data-image-preview]");
    if (input.type === "file" && preview && input.files.length > 0) {
        preview.src = URL.createObjectURL(input.files[0]);
        preview.classList.remove("hidden");
    }
});
