const createElementOfString = (s) => {
    let element = document.createElement('div')
    element.innerHTML = s
    return element.firstChild
}