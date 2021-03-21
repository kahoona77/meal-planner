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

        this.image = this.shadowRoot.querySelector('img');
        const src = this.getAttribute('src');
        if (src) {
            this.image.src = src;
        }

        this.fileSelect = this.shadowRoot.querySelector('.file-select');
        this.fileSelect.addEventListener("change", this.changeImage);
        this.shadowRoot.querySelector(".add-icon").addEventListener("click", this.openFileSelect);
    }

    changeImage = (e) => {
        if (this.fileSelect.files.length > 0) {
            this.value = this.fileSelect.files[0];
            this.image.src = URL.createObjectURL(this.value);
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
      position: relative;
    }
    img {
      max-width: 100%;
    }
    .add-icon {
        cursor: pointer;
        fill: #4a4a4a;
        height: 24px;
        width: 24px;
        position: absolute;
        bottom: 10px;
        right: 10px;
    }
  </style>
  
  <div class="image-select">
    <img src="https://via.placeholder.com/150" alt="meal">
    <svg class="add-icon">
      <use xlink:href="assets/img/solid.svg#plus"></use>
    </svg>

    <input class="file-select" type="file" accept="image/png, image/jpeg" hidden>
  </div>`;

customElements.define("image-select", ImageSelect);