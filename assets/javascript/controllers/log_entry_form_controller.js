import { Controller } from "@hotwired/stimulus"

export default class LogEntryFormController extends Controller {
  static targets = [
    "brewMethodSelect",
    "coffeeTextBox",
    "submitButton",
  ]

  connect() {
    this.enableOrDisableSubmit()
  }

  clear() {
    const emptyValue = (input) => {
      input.value = ""
    }

    this.element.querySelectorAll("input.form-control").forEach(emptyValue)
    this.element.querySelectorAll("textarea").forEach(emptyValue)
    this.brewMethodSelectTarget.selectedIndex = 0
  }

  enableOrDisableSubmit() {
    const val = this.coffeeTextBoxTarget.value.trim()
    this.submitButtonTarget.disabled = !val
  }
}
