class Chat {
    static CreateUserElement = (userObj) => {
        return createElementOfString(`<div class="chat-user">
            <img class="chat-user-image" src="/profile-pictures/${userObj.image}">
            <div class="chat-user-name">${userObj.nickname}</div>
            </div>
        `)
    }
    
    static CreateMessageElement = (body, image, time, me) => {
        return createElementOfString(`<div class="chat-message ${me ? 'chat-message-me' : ''}">
            <div class="chat-message-header">
                <img class="chat-message-image" src=/profile-pictures/${image}>
                <div class="chat-message-time">${time}</div>
            </div>
            <div class="chat-message-text">${body}</div>
        </div>`)
    }

    static GetHeader() {
        return document.querySelector('#open-chat-user-name').textContent
    }

    static SetHeader(nickname, image) {
        document.querySelector('#open-chat-user-name').textContent = nickname
        document.querySelector('#open-chat-user-image').src = `/profile-pictures/${image}`
    }

    static Clear() {
        document.querySelector('#chat-messages').innerHTML = ''
    }

    static formatDate(date) {
        return moment(date).format('DD. MMM, YYYY<br/>HH:mm:ss')
    }

    static GetMessageCount() {
        return document.querySelectorAll('.chat-message').length
    }

    /* static SortOnlineList() {
        let users = document.querySelectorAll('#chat-recent-users li')

        // Loop through everything
        for(let i = 0; i < users.length; i++) {
            // Get the current text
            let firstUserText = users[i].textContent.toLowerCase().trim()

            let smallestIndex = i
            
            // Loop through everything else from the point of i
            for(let j = i; j < users.length; j++) {
                // Second text
                let secondUserText = users[j].textContent.toLowerCase().trim()

                // If second one is ahead in the alphabet
                if(secondUserText < firstUserText) {
                    // Switch the users
                    //users[j].parentNode.insertBefore(users[j], users[i])
                    smallestIndex = j
                }
            }

            users[i].parentNode.insertBefore(users[smallestIndex], users[i])
        }
    } */
    
    constructor(socket) {
        this.socket = socket

        // All the user's data goes into this list
        // This is essential for displaying and accessing them
        this.userManager = new UserManager()

        // The user of the client, will be defined later
        this.clientUser = undefined

        this.socket.addEventListener('open', _ => {
            // OPEN & CLOSE the chat
            $("#chat-close-button").click(function() {
                $('#chat').hide("slow", "swing")
            })

            $("#open-chat").click(function() {
                $('#chat').show("slow", "swing")
            })

            // Logging out
            this.loggedOut = false
            let logoutButton = document.querySelector('#logout-button')
            logoutButton.addEventListener('click', event => {
                if(!this.loggedOut) {
                    event.preventDefault()
                }

                this.socket.send(JSON.stringify({ messageType: 'logout' }))
                this.loggedOut = true
                logoutButton.click()
            })

            // Socket stuff
            this.socket.addEventListener('message', event => {
                
                // Convert the string data into a javascript object
                let data = JSON.parse(event.data)
                
                switch(data.messageType) {
                    case 'message-server-client': // Receiving a message to the bottom
                        // If the message is sent in the active chat or not
                        switch(this.userManager.getUserByNickname(Chat.GetHeader())?.id) {
                            case data.content.from:
                            case data.content.to:
                                this.addMessageBottom(data.content.body, data.content.date, data.content.from)
                                break
                            default:
                                // Make sure we are not receiving the message from ourselves
                                // (this could easily happen due to some server-side delays)
                                // (for example when the user switches tabs really quickly)
                                if(data.content.from != this.userManager.getClientUser().id) {
                                    // Send a notification
                                    alert(`THIS IS A NOTIFICATION! \nYou have received a message from a user named ${
                                        this.userManager.getUserById(data.content.from).nickname
                                    }! `)
                                }
                                break
                        }
                        break
                    case 'history': // Receiving a message to the top. 
                        // Make sure that there are messages, that the array is not empty
                        if(data.content) {
                            this.addMessagesTop(data.content)
                        }
                        break
                    case 'update-user-list': // Updating the user lists
                        data.content.online.sort((a, b) => {
                            return a.nickname.toLowerCase() < b.nickname.toLowerCase() ? -1 : 1
                        })
                        this.updateUserList(data.content.recent, data.content.online)
                        break
                }
            })

            // Loading the past 10 messages when the client scrolls up
            let chatMessages = document.querySelector('#chat-messages')
            chatMessages.addEventListener('scroll', event => {
                // Once we have reached the top
                if(chatMessages.scrollTop == 0) {
                    // Get the id of the currently open chat
                    let user = this.userManager.getUserByNickname(Chat.GetHeader())
                    if(!user) {
                        return
                    }

                    // Request message history
                    this.socket.send(JSON.stringify({ messageType: 'request-history', content: { id: user.id, index: Chat.GetMessageCount() } }))
                }
            })
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
                Chat.SetHeader(user.nickname, user.image)

                // Clear the messages
                Chat.Clear()

                //! Request last 10 messages (THIS IS RECENT)
                this.socket.send(JSON.stringify({ messageType: 'request-history', content: { id: user.id, index: Chat.GetMessageCount() } }))
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
                Chat.SetHeader(user.nickname, user.image)

                // Clear the messages
                Chat.Clear()

                //! Request 10 last messages, even though we shouldn't
                this.socket.send(JSON.stringify({ messageType: 'request-history', content: { id: user.id, index: Chat.GetMessageCount() } }))
            })
        })
    }

    // To the bottom, with scrolling down
    addMessageBottom(body, date, from) {
        // Get the active chat user
        let user = this.userManager.getUserById(from)

        // Determine whether the message was sent by us or not
        let fromClient = from == this.userManager.getClientUser().id ? true : false

        let messageElement = Chat.CreateMessageElement(body, user.image, Chat.formatDate(date), fromClient)
        let chatMessages = document.querySelector('#chat-messages')
        chatMessages.appendChild(messageElement)

        // And scroll all the way down
        chatMessages.scrollTop = chatMessages.scrollHeight
    }

    // To the top, without scrolling to the top automatically, just adding the messages
    addMessagesTop(messages) {
        let chatMessages = document.querySelector('#chat-messages')
        let lastHeight = chatMessages.scrollHeight

        messages.forEach(message => {
            // Make sure the message is sent either by us or by them, otherwise don't add anything
            if(message.from != this.userManager.getUserByNickname(Chat.GetHeader()).id &&
               message.to != this.userManager.getUserByNickname(Chat.GetHeader()).id) {
                return
            }

            let user = this.userManager.getUserById(message.from)

            // Determine whether the message was sent by us or not
            let fromClient = message.from == this.userManager.getClientUser().id ? true : false

            // Create the message element
            let messageElement = Chat.CreateMessageElement(message.body, user.image, Chat.formatDate(message.date), fromClient)
            
            // Append the message element to the chat
            chatMessages.insertBefore(messageElement, chatMessages.firstChild)
        })

        // Make sure the scrolling comes with the addition of the new message
        //chatMessages.pageYOffset += chatMessages.scrollHeight - previousScrollHeight
        chatMessages.scrollTop += chatMessages.scrollHeight - lastHeight
    }

    sendMessage(body) {
        // The ID of the user that we are sending the message to
        let user = this.userManager.getUserByNickname(Chat.GetHeader())
        if(!user) {
            return
        }

        // Send the message to the user through the socket connection
        let message = JSON.stringify({ messageType: "message-client-server", content: { body: body, to: user.id } })
        this.socket.send(message)
     }
}