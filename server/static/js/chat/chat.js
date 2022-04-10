class Chat {
    static CreateUserElement = (userObj) => {
        return createElementOfString(`<div class="chat-user">
            <img class="chat-user-image" src="${userObj.image}">
            <div class="chat-user-name">${userObj.nickname}</div>
            </div>
        `)
    }
    
    static CreateMessageElement = (body, image, time, me) => {
        return createElementOfString(`<div class="chat-message ${me ? 'chat-message-me' : ''}">
            <div class="chat-message-header">
                <img class="chat-message-image" src=${image}>
                <div class="chat-message-time">${time}</div>
            </div>
            <div class="chat-message-text">${body}</div>
        </div>`)
    }

    static GetHeader() {
        return document.querySelector('#open-chat-user-name').textContent
    }

    static SetHeader(nickname) {
        document.querySelector('#open-chat-user-name').textContent = nickname
    }

    static Clear() {
        document.querySelector('#chat-messages').innerHTML = ''
    }
    
    constructor(socket) {
        this.socket = socket

        // All the user's data goes into this list
        // This is essential for displaying and accessing them
        this.userManager = new UserManager()

        // The user of the client
        this.clientUser = undefined

        window.onkeydown = e => {
            if(e.code == 'KeyW') {
                this.socket.dispatchEvent(new Event('message', {
                    code: 100,
                    data: {
                        body: 'This is a message added from the "message" listener',
                        from: 3,
                        to: 1
                    }
                }))
            }
        }

        this.socket.addEventListener('open', _ => {
            // OPEN & CLOSE the chat
            $("#chat-close-button").click(function() {
                $('#chat').hide("slow", "swing")
            })

            $("#open-chat").click(function() {
                $('#chat').show("slow", "swing")
            })

            // Socket stuff
            this.socket.addEventListener('message', event => {
                event = {
                    code: 100,
                    data: {
                        body: 'This is a message added from the "message" listener',
                        time: '19:27',
                        from: 0,
                        to: 3
                    }
                }

                switch(event.code) {
                    case 100: // Receiving a message to the bottom
                        // If the message is sent in the active chat or not
                        switch(this.userManager.getUserByNickname(Chat.GetHeader())?.id) {
                            case event.data.from:
                            case event.data.to:
                                // If the message is sent by us sent by them
                                // Add the message element HTML
                                if(event.data.from == this.userManager.getClientUser().id) {
                                    this.addMessageBottom(event.data.body, event.data.time, true)
                                } else {
                                    this.addMessageBottom(event.data.body, event.data.time)
                                }
                                break
                            default:
                                // Make sure we are not receiving the message from ourselves
                                // (this could easily happen due to some server-side delays)
                                // (for example when the user switches tabs really quickly)
                                if(event.data.from != this.userManager.getClientUser().id) {
                                    // Send a notification
                                    alert(`THIS IS A NOTIFICATION! \nYou have received a message from a user named ${
                                        this.userManager.getUserById(event.data.from).nickname
                                    }! `)
                                }
                                break
                        }
                        break
                    case 101: // Receiving a message to the top. 
                        // Make sure that the message is being sent to the chat we are currently in
                        switch(this.userManager.getUserByNickname(Chat.GetHeader())?.id) {
                            case event.data.from: // These two cases say that the message has to be
                            case event.data.to:   // sent in the currently active chat window
                                // If the message is sent by us or by them. The `fromClient` flag changes
                                if(event.data.from == this.userManager.getClientUser().id) {
                                    this.addMessageTop(event.data.body, event.data.time, true)
                                } else {
                                    this.addMessageTop(event.data.body, event.data.time)
                                }
                        }
                        break
                    case 102: // Updating the user lists
                        this.updateUserList(event.data.recent, event.data.online)
                        break
                }
            })

            // Loading the past 10 messages when the client scrolls up
            let chatMessages = document.querySelector('#chat-messages')
            chatMessages.addEventListener('scroll', event => {
                // Once we have reached the top
                if(chatMessages.scrollTop == 0) {
                    // How many messages we currently have showing in the chat
                    let messageCount = document.querySelectorAll('.chat-message').length

                    // Send the request to the server to get the ten messages
                    this.socket.send({ code: 102, index: messageCount })

                    //! Debugging: Insert ten messages to the top automatically
                    let messages = [
                        { body: 'hi', time: '03:26' },
                        { body: 'Hey!', time: '03:26', fromClient: true },
                        { body: 'so whats up', time: '03:27' },
                        { body: 'Not much honestly. Yourself?', time: '03:27', fromClient: true },
                        { body: 'the same. bored.', time: '03:27' },
                        { body: 'Why?', time: '03:28', fromClient: true },
                        { body: 'well theres this program i need to finish but i do not want to so', time: '03:28' },
                        { body: 'Is it really that bad?', time: '03:28', fromClient: true },
                        { body: 'yes man', time: '03:28' },
                        { body: 'Okay, bye then. ', time: '03:29', fromClient: true }
                    ]
                    for(let i = messages.length - 1; i > -1; i--) {
                        this.addMessageTop(messages[i].body, messages[i].time, messages[i].fromClient)
                    }
                }
            })

            //! Debugging initialization
            this.updateUserList([
                { id: 0, nickname: 'Laura-Eliise', image: '/profile-pictures/0-0.png' },
                { id: 1, nickname: 'Olari', image: '/profile-pictures/0-0.png' }
            ], [
                { id: 2, nickname: 'Kris', image: '/profile-pictures/0-0.png' },
                { id: 3, nickname: 'Kanguste', image: '/profile-pictures/0-0.png' }
            ])
        })
    }

    updateUserList(recent, online) {
        // Make sur the list of users is empty
        this.userManager.clear()

        // Empty recent and add recent users
        let recentList = document.querySelector('#chat-recent-users')
        recentList.innerHTML = ''
        recent.forEach(user => {
            // Add the user to the list
            this.userManager.addUser(user.id, user.nickname, user.image)

            // Create the user HTML and add it to the page
            let userElement = document.createElement('li')
            userElement.appendChild(Chat.CreateUserElement(user))
            recentList.appendChild(userElement)

            // And add an event listener to the user button
            // Open the chat with that user
            userElement.addEventListener('mousedown', event => {
                // Change the header username
                Chat.SetHeader(user.nickname)

                // Clear the messages
                Chat.Clear()

                //! Getting the last 10 messages from the server
                //! this.socket.send({code: 102, index: '0'})
                console.log('Requesting last ten messages...')
            })
        })

        // Empty online and add online users
        let onlineList = document.querySelector('#chat-online-users')
        onlineList.innerHTML = ''
        online.forEach(user => {
            // Add the user to the list
            this.userManager.addUser(user.id, user.nickname, user.image)

            // If it if our user, just break out of it
            // Adding the user, but not creating it in the HTML
            if(user.nickname == this.userManager.getClientUser()?.nickname) {
                return
            }

            // Create the user HTML and add it to the page
            let userElement = document.createElement('li')
            userElement.appendChild(Chat.CreateUserElement(user))
            onlineList.appendChild(userElement)

            // And add an event listener to the user button
            userElement.addEventListener('mousedown', event => {
                // Change the header username
                Chat.SetHeader(user.nickname)

                // Clear the messages
                Chat.Clear()
            })
        })
    }

    // To the bottom, with scrolling down
    addMessageBottom(body, time, fromClient = false) {
        // Get the active chat user
        let user = this.userManager.getClientUser()

        let messageElement = Chat.CreateMessageElement(body, user.image, time, fromClient)
        let chatMessages = document.querySelector('#chat-messages')
        chatMessages.appendChild(messageElement)

        // And scroll all the way down
        chatMessages.scrollTop = chatMessages.scrollHeight
    }

    // To the top, without scrolling to the top automatically, just adding the messages
    addMessageTop(body, time, fromClient = false) {
        let user = this.userManager.getClientUser()

        let messageElement = Chat.CreateMessageElement(body, user.image, time, fromClient)
        let chatMessages = document.querySelector('#chat-messages')
        let previousScrollHeight = chatMessages.scrollHeight
        chatMessages.insertBefore(messageElement, chatMessages.firstChild)

        // Make sure the scrolling comes with the addition of the new message
        chatMessages.scrollTop += chatMessages.scrollHeight - previousScrollHeight
    }

    sendMessage(body) {
        // The ID of the user that we are sending the message to
        let toID = this.userManager.getUserByNickname(Chat.GetHeader()).id

        // Send the message to the user through the socket connection
        this.socket.send({ code: 103, body: body, to: toID })
     }
}