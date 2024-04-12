import { LitElement, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'


@customElement('toggle-visibility')
export class ToggleVisibility extends LitElement {

  @property({type: String})
  selector: string|undefined = undefined;

  render() {
    return html`
      <slot @click=${this._onClick}></slot>
    `
  }

  private _onClick() {
    if (!this.selector) {
      console.warn("no selector set!")
      return;
    }

    document.querySelector(this.selector)?.classList.toggle("hidden");
  }

}

declare global {
  interface HTMLElementTagNameMap {
    'toggle-visibility': ToggleVisibility
  }
}
