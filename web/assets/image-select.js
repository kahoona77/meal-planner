class ImageSelect extends HTMLElement {
    static formAssociated = true;

    constructor() {
        super();
        // Get access to the internal form control APIs
        this._internals = this.attachInternals();
        // internal value for this control
        this._value = null;
    }

    // Form controls usually expose a "value" property
    get value() { return this._value; }
    set value(v) {
        this._value = v;
        this._internals.setFormValue(this._value);
    }

    // The following properties and methods aren't strictly required,
    // but browser-level form controls provide them. Providing them helps
    // ensure consistency with browser-provided controls.
    get form() { return this._internals.form; }
    get name() { return this.getAttribute('name'); }
    get type() { return this.localName; }
    get validity() {return this._internals.validity; }
    get validationMessage() {return this._internals.validationMessage; }
    get willValidate() {return this._internals.willValidate; }

    checkValidity() { return this._internals.checkValidity(); }
    reportValidity() {return this._internals.reportValidity(); }


    connectedCallback() {
        this.attachShadow({mode: 'open'});
        this.shadowRoot.appendChild(template.content.cloneNode(true));

        this.placeholder = this.shadowRoot.querySelector('.placeholder');
        this.image = this.shadowRoot.querySelector('img');
        const src = this.getAttribute('src');
        if (src) {
            this.showImage(src);
        }

        this.fileSelect = this.shadowRoot.querySelector('.file-select');
        this.fileSelect.addEventListener("change", this.changeImage);
        this.shadowRoot.querySelector(".add-icon").addEventListener("click", this.openFileSelect);
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
    .image {
        position: relative;
        --image-size: 12rem;
        
        width: var(--image-size);
        height: var(--image-size);
    }
    img {
      width: var(--image-size);
      height: var(--image-size);
      object-fit: cover;
    }
    
    .placeholder {
        background-color: var(--grey-lighter);
        width: var(--image-size);
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .placeholder svg {
        fill: var(--grey-lightest);
        width: 4rem;
        height: 4rem;
    }
  </style>
  
  <div class="image-select">
    <div class="image">
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

    <input class="file-select" type="file" accept="image/png, image/jpeg" hidden>
  </div>`;

customElements.define("image-select", ImageSelect);