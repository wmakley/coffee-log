import {Controller} from "@hotwired/stimulus"

export default class ConfirmationController extends Controller {
  static values = { message: String }

  confirm(event) {
    if (!confirm(this.messageValue)) {
      event.preventDefault()
      event.stopPropagation()
    }
  }
}
