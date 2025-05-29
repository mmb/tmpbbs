const debounceInterval = 150,
      emojiSuggestions = document.getElementById("emoji-suggestions"),
          qrCodeDialog = document.getElementById("qr")

class EmojiSuggester {
  constructor(input, displayContainer) {
    this.input = input
    this.displayContainer = displayContainer
    this.minPrefixLength = 3
    this.whitespaceRegex = /\s/u
  }

  update() {
    this.#updateCurrentWord()
    this.#getSuggestions()
  }

  // eslint-disable-next-line max-statements
  #updateCurrentWord() {
    let startIndex = this.input.selectionStart
    if ((this.#indexIsWhitespace(startIndex) || startIndex === this.input.value.length) && startIndex > 0 && !this.#indexIsWhitespace(startIndex - 1)) {
      startIndex -= 1
    }
    // eslint-disable-next-line one-var
    let endIndex = startIndex
    while (startIndex > 0 && !this.#indexIsWhitespace(startIndex - 1)) {
      startIndex -= 1
    }
    while (endIndex < this.input.value.length - 1 && !this.#indexIsWhitespace(endIndex) && !this.#indexIsWhitespace(endIndex + 1)) {
      endIndex += 1
    }

    this.currentWord = this.input.value.slice(startIndex, endIndex + 1)
    this.currentWordStartIndex = startIndex
    this.currentWordEndIndex = endIndex
  }

  async #getSuggestions() {
    if (!this.currentWord.startsWith(":") || this.currentWord.length < this.minPrefixLength) {
      this.displayContainer.innerHTML = ""

      return
    }

    const params = new URLSearchParams([["q", this.currentWord]]),
     suggestions = await (await fetch(`/emoji-suggest?${params}`)).json()

    this.displayContainer.innerHTML = ""
    suggestions.forEach(suggestion => {
      const button = document.createElement("button"),
         emojiSpan = document.createElement("span")
      button.type = "button"
      button.className = "emoji-suggestion"
      button.addEventListener('mouseup', () => {
        this.#replaceCurrentWord(suggestion.suggestion)
        this.input.focus()
      })
      emojiSpan.className = "emoji"
      emojiSpan.innerHTML = suggestion.pictogram
      button.append(emojiSpan, document.createTextNode(` ${suggestion.suggestion}`))
      this.displayContainer.append(button)
    })
  }

  #replaceCurrentWord(replacement) {
    this.input.value = this.input.value.substring(0, this.currentWordStartIndex) + replacement + this.input.value.substring(this.currentWordEndIndex + 1)
    this.input.selectionStart = this.currentWordStartIndex + replacement.length
    this.input.selectionEnd =  this.input.selectionStart
    this.update()
  }

  #indexIsWhitespace(index) {
    return this.whitespaceRegex.test(this.input.value[index])
  }
}

if (emojiSuggestions) {
  const bodyTextarea = document.getElementById("body"),
      emojiSuggester = new EmojiSuggester(bodyTextarea, emojiSuggestions),
      events = ["input", "mouseup"]
    let timeout
    events.forEach(event => {
      bodyTextarea.addEventListener(event, () => {
        clearTimeout(timeout)
        timeout = setTimeout(() => emojiSuggester.update(), debounceInterval)
      })
    })
}

if (qrCodeDialog) {
  const button = document.getElementById("show-qr"),
           img = document.createElement("img")
  img.src = `/qr?url=${encodeURIComponent(window.location.href)}`
  img.loading = "lazy"
  qrCodeDialog.prepend(img)

  button.addEventListener("click", () => { qrCodeDialog.showModal() })
  button.style.display = ""
}
