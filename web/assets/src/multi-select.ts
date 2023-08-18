import {LitElement, css, html, nothing} from 'lit'
import {customElement, property} from 'lit/decorators.js'
import {classMap} from 'lit/directives/class-map.js';
import {ComplexAttributeConverter} from "@lit/reactive-element";

class SelectOption {
    constructor(public id: string, public name: string) {
    }
}

class SelectOptionsConverter implements ComplexAttributeConverter<Array<SelectOption>> {
    fromAttribute(value: string) {
        let selected: Array<any> = [];
        if (value) {
            selected = JSON.parse(value)
        }

        if (selected) {
            selected = selected.map((s: any) => new SelectOption(s.id, s.name));
        } else {
            selected = [];
        }

        return selected;
    }

    toAttribute(selected: Array<SelectOption>) {
        return JSON.stringify(selected);
    }
}

@customElement('multi-select')
export class MultiSelect extends LitElement {
    static formAssociated = true;
    private internals: ElementInternals;

    constructor() {
        super();
        this.internals = this.attachInternals();
    }

    @property({
        attribute: 'value',
        type: Array<SelectOption>,
        converter: new SelectOptionsConverter(),
    })
    selected: Array<SelectOption> = [];

    @property({
        attribute: 'options',
        type: Array<SelectOption>,
        converter: new SelectOptionsConverter(),
    })
    items: Array<SelectOption> = [];

    @property({attribute: false})
    showItems: boolean = false;

    get form() {
        return this.internals.form;
    }

    get name() {
        return this.getAttribute('name')
    };

    private updateFormValue(): void {
        const value = new SelectOptionsConverter().toAttribute(this.selected);
        this.internals.setFormValue(value);
    }

    private removeSelected(event: Event, selectedItem: SelectOption) {
        event.stopPropagation();
        const indexToRemove = this.selected.indexOf(selectedItem);
        this.selected = [...this.selected.slice(0, indexToRemove), ...this.selected.slice(indexToRemove + 1)];
        this.updateFormValue();
        this.showItems = false;
    }

    private addItem(item: SelectOption) {
        if (this.selected.indexOf(item) < 0) {
            this.selected = [...this.selected, item];
            this.updateFormValue();
            this.showItems = false;
        }
    }

    render() {
        return html`
            <!-- component -->
            <div class="container">
                <div class="inner-container">
                    <div class="selection">
                        <div class="selected" @click=${() => this.showItems = true}>
                            ${this.selected.map((select) =>
                                    html`
                                        <div class="item">
                                            <div class="text">${select.name}</div>
                                            <svg @click=${(event: Event) => this.removeSelected(event, select)}
                                                 xmlns="http://www.w3.org/2000/svg" width="100%" height="100%"
                                                 fill="none" viewBox="0 0 24 24" stroke-width="2"
                                                 stroke-linecap="round" stroke-linejoin="round">
                                                <line x1="18" y1="6" x2="6" y2="18"></line>
                                                <line x1="6" y1="6" x2="18" y2="18"></line>
                                            </svg>
                                        </div>
                                    `
                            )}
                        </div>
                        <div class="toggle">
                            <button>
                                <svg xmlns="http://www.w3.org/2000/svg" width="100%" height="100%" fill="none"
                                     viewBox="0 0 24 24" stroke-width="2"
                                     stroke-linecap="round" stroke-linejoin="round">
                                    <polyline points="18 15 12 9 6 15"></polyline>
                                </svg>
                            </button>
                        </div>
                    </div>

                    ${
                            this.showItems ?
                                    html`
                                        <div class="items-container">
                                            ${this.items.map((item) => html`
                                                <div class="${classMap({
                                                    item: true,
                                                    disabled: this.selected.indexOf(item) >= 0
                                                })}" @click=${() => this.addItem(item)}>${item.name}
                                                </div>`)}
                                        </div>` : nothing
                    }

                </div>

            </div>
        `
    }


    static styles = css`
      :host {
        --ms-text-color: var(--text-color, rgb(15 118 110));
        --ms-accent-color: var(--accent-color, rgb(15 118 110));
        --ms-bg-color: var(--bg-color, rgb(15 118 110));
      }

      * {
        box-sizing: border-box;
      }

      .container {
        align-items: center;
        flex-direction: column;
        width: 100%;
        display: flex;
      }

      .inner-container {
        position: relative;
        align-items: center;
        flex-direction: column;
        display: flex;
        width: 100%;
      }

      .selection {
        width: 100%;
        padding: 0.25rem;
        background-color: white;
        border-color: rgb(229 231 235);
        border-width: 1px;
        border-style: solid;
        border-radius: 0.25rem;
        display: flex;
        margin-top: 0.5rem;
        margin-bottom: 0.5rem;
      }

      .selected {
        display: flex;
        flex-wrap: wrap;
        flex: 1 1 auto;
      }

      .selected .item {
        color: var(--ms-text-color);
        font-weight: 500;
        padding: 0.25rem 0.5rem;
        background-color: var(--ms-bg-color);
        border-color: var(--ms-accent-color);
        border-width: 1px;
        border-style: solid;
        border-radius: 9999px;
        justify-content: center;
        align-items: center;
        display: flex;
        margin: 0.25rem;
      }

      .selected .item .text {
        font-weight: 400;
        font-size: 0.75rem;
        line-height: 1rem;
        flex: 0 1 auto;
        max-width: 100%;
      }

      .selected .item svg {
        border-radius: 9999px;
        cursor: pointer;
        width: 1rem;
        height: 1rem;
        margin-left: 0.5rem;
        stroke: var(--ms-accent-color);
      }

      .selected .item svg:hover {
        filter: brightness(200%);
      }

      .toggle {
        color: rgb(209 213 219);
        padding: 0.25rem 0.25rem 0.25rem 0.5rem;
        border-color: rgb(229 231 235);
        border-width: 0 0 0 1px;
        border-style: solid;
        align-items: center;
        width: 2rem;
        display: flex;
      }

      .toggle button {
        outline: 2px solid transparent;
        background-color: transparent;
        border: none;
        color: var(--ms-text-color);
        cursor: pointer;
        width: 1.5rem;
        height: 1.5rem;
        padding: 0;
      }

      .toggle svg {
        stroke: var(--ms-text-color);
        width: 1rem;
        height: 1rem;
      }

      .items-container {
        max-height: 300px;
        top: 100%;
        box-shadow: 0 0 0 2px white, 0.3em 0.3em 1em rgba(0, 0, 0, 0.3);
        border-radius: 0.25rem;
        background-color: white;
        overflow-y: auto;
        width: 100%;
        z-index: 40;
        position: absolute;
        flex-direction: column;
        display: flex;
      }

      .items-container .item {
        border-color: rgb(243 244 246);
        border-bottom-width: 1px;
        border-top-left-radius: 0.25rem;
        border-top-right-radius: 0.25rem;
        cursor: pointer;
        width: 100%;
        padding: 0.5rem;
      }

      .items-container .item:hover {
        background-color: var(--ms-bg-color);
      }

      .items-container .item.disabled {
        cursor: not-allowed;
        filter: brightness(200%);
      }

      .items-container .item.disabled:hover {
        background-color: transparent;
      }
    `
}

declare global {
    interface HTMLElementTagNameMap {
        'multi-select': MultiSelect
    }
}
