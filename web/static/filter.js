// Filters the items of a [data-filter] element while typing: an item ([data-filter-item]) stays
// visible if its data-name contains the text of the [data-filter-search] field and its data-tags
// (tag ids, separated by spaces) contain every checked checkbox in [data-filter-tags].
// [data-filter-empty] is shown if no item is left.
document.addEventListener("input", (event) => {
    const filter = event.target.closest("[data-filter]");
    if (filter) {
        applyFilter(filter);
    }
});

function applyFilter(filter) {
    const text = normalize(filter.querySelector("[data-filter-search]")?.value ?? "");
    const tags = [...filter.querySelectorAll("[data-filter-tags] input:checked")].map((input) => input.value);

    let visible = 0;
    for (const item of filter.querySelectorAll("[data-filter-item]")) {
        const itemTags = (item.dataset.tags ?? "").split(" ");
        item.hidden = !normalize(item.dataset.name ?? "").includes(text) || !tags.every((tag) => itemTags.includes(tag));
        if (!item.hidden) {
            visible++;
        }
    }

    const empty = filter.querySelector("[data-filter-empty]");
    if (empty) {
        empty.hidden = visible > 0;
    }
}

// lower case without accents, so "kurbis" finds "Kürbissuppe"
function normalize(text) {
    return text.toLocaleLowerCase().normalize("NFD").replace(/\p{Diacritic}/gu, "").trim();
}
