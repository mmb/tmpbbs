const commit = document.querySelector('meta[name="commit"]').content,
  dateTime = new Intl.DateTimeFormat(document.documentElement.lang, {
    day: 'numeric',
    hour: 'numeric',
    minute: 'numeric',
    month: 'numeric',
    second: 'numeric',
    year: 'numeric'
  }),
  debounceInterval = 150,
  emojiSuggestions = document.querySelector('#emoji-suggestions'),
  minPrefixLength = 2,
  msPerSecond = 1000,
  qrCodeDialog = document.querySelector('#qr'),
  relativeTime = new Intl.RelativeTimeFormat(document.documentElement.lang, { numeric: 'auto', style: 'short' }),
  secondsPerDay = 86_400,
  secondsPerHour = 3600,
  secondsPerMinute = 60,
  timeAgo = (date) => {
    const secondsAgo = Math.max(0, Math.floor((Date.now() - date.getTime()) / msPerSecond))
    if (secondsAgo < secondsPerMinute) {
      return relativeTime.format(-secondsAgo, 'second')
    }
    if (secondsAgo < secondsPerHour) {
      return relativeTime.format(-Math.round(secondsAgo / secondsPerMinute), 'minute')
    }
    if (secondsAgo < secondsPerDay) {
      return relativeTime.format(-Math.round(secondsAgo / secondsPerHour), 'hour')
    }
    return relativeTime.format(-Math.round(secondsAgo / secondsPerDay), 'day')
  },
  timeAgoUpdateInterval = 60_000,
  updateTimeAgo = () => {
    for (const element of document.querySelectorAll('time')) {
      const date = new Date(element.dateTime)
      element.textContent = timeAgo(date)
      if (!element.title) {
        element.title = dateTime.format(date)
      }
    }
  },
  whitespaceRegex = /\s/u

updateTimeAgo()
setInterval(updateTimeAgo, timeAgoUpdateInterval)

class EmojiSuggester {
  static async #fetchSuggestions(query) {
    const params = new URLSearchParams([['q', query]]),
      response = await fetch(`/emoji-suggest/${commit}?${params}`, { cache: 'default' })
    return response.ok
      ? response.json()
      : []
  }

  #currentWord = ''
  #currentWordStartIndex = 0
  #currentWordEndIndex = 0
  #displayContainer
  #input

  constructor(input, displayContainer) {
    this.#input = input
    this.#displayContainer = displayContainer
  }

  #update() {
    this.#updateCurrentWord()
    this.#updateSuggestions()
  }

  // eslint-disable-next-line max-statements
  #updateCurrentWord() {
    let startIndex = this.#input.selectionStart
    if ((this.#indexIsWhitespace(startIndex) || startIndex === this.#input.value.length) && startIndex > 0 && !this.#indexIsWhitespace(startIndex - 1)) {
      startIndex -= 1
    }
    // eslint-disable-next-line one-var
    let endIndex = startIndex
    while (startIndex > 0 && !this.#indexIsWhitespace(startIndex - 1)) {
      startIndex -= 1
    }
    while (endIndex < this.#input.value.length - 1 && !this.#indexIsWhitespace(endIndex) && !this.#indexIsWhitespace(endIndex + 1)) {
      endIndex += 1
    }

    this.#currentWord = this.#input.value.slice(startIndex, endIndex + 1)
    this.#currentWordStartIndex = startIndex
    this.#currentWordEndIndex = endIndex
  }

  #indexIsWhitespace(index) {
    return whitespaceRegex.test(this.#input.value[index])
  }

  async #updateSuggestions() {
    const query = this.#currentWord.replaceAll(/^:|:$/gu, '').toLowerCase()

    if (!this.#currentWord.startsWith(':') || query.length < minPrefixLength) {
      this.#clearSuggestions()

      return
    }

    // eslint-disable-next-line one-var
    const suggestions = await EmojiSuggester.#fetchSuggestions(query)
    this.#clearSuggestions()
    for (const suggestion of suggestions) {
      this.#addSuggestionButton(suggestion)
    }
  }

  #clearSuggestions() {
    this.#displayContainer.replaceChildren()
  }

  #addSuggestionButton(suggestion) {
    const button = document.createElement('button'),
      emojiSpan = document.createElement('span')
    button.type = 'button'
    button.classList.add('emoji-suggestion')
    button.addEventListener('click', () => {
      this.#replaceCurrentWord(`:${suggestion.suggestion}:`)
      this.#input.focus()
    })
    emojiSpan.classList.add('emoji')
    emojiSpan.textContent = suggestion.pictogram
    button.append(emojiSpan, document.createTextNode(` ${suggestion.suggestion}`))
    this.#displayContainer.append(button)
  }

  #replaceCurrentWord(replacement) {
    this.#input.value = this.#input.value.substring(0, this.#currentWordStartIndex) + replacement + this.#input.value.substring(this.#currentWordEndIndex + 1)
    this.#input.selectionStart = this.#currentWordStartIndex + replacement.length
    this.#input.selectionEnd = this.#input.selectionStart
    this.#update()
  }

  listen() {
    let timeout
    const events = ['input', 'keyup', 'mouseup'],
      listener = () => {
        clearTimeout(timeout)
        timeout = setTimeout(() => this.#update(), debounceInterval)
      }

    for (const event of events) {
      this.#input.addEventListener(event, listener)
    }
  }
}

if (emojiSuggestions) {
  const emojiSuggester = new EmojiSuggester(document.querySelector('#body'), emojiSuggestions)
  emojiSuggester.listen()
}

if (qrCodeDialog) {
  const button = document.querySelector('#show-qr'),
    img = document.createElement('img')
  img.src = `/qr/${commit}?url=${encodeURIComponent(location.href)}`
  img.loading = 'lazy'
  qrCodeDialog.prepend(img)

  button.addEventListener('click', () => {
    qrCodeDialog.showModal()
  })
  button.hidden = false
}
