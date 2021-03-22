class ImageSelect extends HTMLElement {
    connectedCallback() {
        this.appendChild(template.content.cloneNode(true));

        this.placeholder = this.querySelector('.placeholder');
        this.image = this.querySelector('img');
        const src = this.getAttribute('src');
        if (src) {
            this.showImage(src);
        }

        this.fileSelect = this.querySelector('.file-select');
        this.fileSelect.setAttribute("name", this.getAttribute('name'));
        this.fileSelect.addEventListener("change", this.changeImage);
        this.querySelector(".add-icon").addEventListener("click", this.openFileSelect);
    }

    showImage(src) {
        this.image.src = src;
        this.image.classList.remove("hidden");
        this.placeholder.classList.add("hidden");
    }

    changeImage = (e) => {
        if (this.fileSelect.files.length > 0) {
            this.value = this.fileSelect.files[0];
            this.showImage(URL.createObjectURL(this.value));
        }
    }

    openFileSelect = (e) => {
        this.fileSelect.click();
    }
}

const template = document.createElement('template');
template.innerHTML = `
  <style>
    .image-select {
      width: 100%;
      display: flex;
      justify-content: center;
    }
    .hidden {display: none !important;}
    
    .add-icon {
        cursor: pointer;
        fill: #949494;
        height: 24px;
        width: 24px;
        position: absolute;
        bottom: 10px;
        right: 10px;
        padding: 2px;
        background-color: #e0e0e0;
        border-radius: 13px;
    }
    .image-select .image-select-image {
        position: relative;
        --image-size: 12rem;
        
        width: var(--image-size);
        height: var(--image-size);
    }
    .image-select .image-select-image img {
      width: var(--image-size);
      height: var(--image-size);
      object-fit: cover;
    }
    
    .image-select .image-select-image .placeholder {
        background-color: var(--grey-lighter);
        width: var(--image-size);
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .image-select .image-select-image .placeholder svg {
        fill: var(--grey-lightest);
        width: 4rem;
        height: 4rem;
    }
  </style>
  
  <div class="image-select">
    <div class="image-select-image">
        <img src="" alt="meal" class="hidden">
        <div class="placeholder">
            <svg>
                <use xlink:href="assets/img/solid.svg#hamburger"></use>
            </svg>
        </div>
        <svg class="add-icon">
          <use xlink:href="assets/img/solid.svg#plus"></use>
        </svg>
    </div>

    <input class="file-select hidden" type="file" accept="image/png, image/jpeg" hidden>
  </div>`;

customElements.define("image-select", ImageSelect);