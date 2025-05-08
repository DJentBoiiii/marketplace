function updateFormFields() {
    const typeSelect = document.getElementById("typeSelect");
    const genreContainer = document.getElementById("genreContainer");
    const typeContainer = document.getElementById("typeContainer");

    if (typeSelect.value === "audio") {
        genreContainer.style.display = "block";
        typeContainer.style.display = "none";
    } else if (typeSelect.value === "midi" || typeSelect.value === "samples") {
        genreContainer.style.display = "none";
        typeContainer.style.display = "block";
    } else {
        genreContainer.style.display = "none";
        typeContainer.style.display = "none";
    }
}

document.getElementById("typeSelect").addEventListener("change", updateFormFields);

// Initial call to set correct visibility on page load
updateFormFields();

document.getElementById("fileInput").addEventListener("change", function () {
    const file = this.files[0];
    const typeSelect = document.getElementById("typeSelect");

    if (!file || !typeSelect.value) {
        alert("Оберіть файл і тип.");
        return;
    }

    const audioFormats = [".mp3", ".wav", ".flac", ".ogg", ".m4a"];
    const archiveFormats = [".zip"];
    const fileExt = file.name.split(".").pop().toLowerCase();

    if (typeSelect.value === "audio" && !audioFormats.includes("." + fileExt)) {
        alert("Файл має бути аудіоформатом.");
        this.value = "";
    } else if ((typeSelect.value === "midi" || typeSelect.value === "samples") && !archiveFormats.includes("." + fileExt)) {
        alert("Файл має бути ZIP-архівом.");
        this.value = "";
    }
});

document.getElementById("imageInput").addEventListener("change", function () {
    const file = this.files[0];
    if (!file) return;

    const imgFormats = [".jpg", ".jpeg", ".png", ".webp"];
    const fileExt = file.name.split(".").pop().toLowerCase();

    if (!imgFormats.includes("." + fileExt)) {
        alert("Зображення має бути у форматі JPG, JPEG, PNG або WebP.");
        this.value = "";
    }
});
