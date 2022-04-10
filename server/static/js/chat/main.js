$(document).ready(function() {
    var chat = new Chat(new WebSocket('ws://localhost:8080/chat'))

    document.querySelector('#chat-form').addEventListener('submit', event => {
        event.preventDefault()

        let chatInput = document.querySelector('#chat-message-input')

        //! Sending the message to the server
        //! this.chat.sendMessage(chatInput)

        //* Testing adding messages
        chat.addMessageBottom(chatInput.value, '19:38', true)

        // Empty the input field
        chatInput.value = ''
    })
})

// Message from client to server
/*
    {
        to <int>, 
        body <string>
    }
*/

// Message from server to client
/* 
    {
        id: <int>, 
        from: <int>, 
        to: <int>, 
        body: <string>, 
        time: <string>(12:34)
    }
*/

// User from server to client
/* 
    {
        id: <int>,
        nickname: <string>,
        image: <string>(url of the image)
    }
*/