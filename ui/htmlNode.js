function Node(v) {
    this.el = document.createElement(v)
    this.addText = v => {
        this.el.appendChild(document.createTextNode(v))
        return this
    }
    this.append = v => {
        this.el.appendChild(v.el)
        return this
    }
    this.addClass = v => {
        this.el.classList.add(v)
        return this
    }
    this.removeClass = v => {
        this.el.classList.remove(v)
        return this
    }
    this.containsClass = v => {
        return this.el.classList.contains(v)
    }
    this.appendTo = v => {
        v.appendChild(this.el)
        return this
    }
    this.prepend = v => {
        v.prepend(this.el)
        return this
    }
}