class UserManager {
    constructor() {
        this.userList = []
    }

    // Getters
    getUserById(id) {
        for(let user of this.userList) {
            if(id == user.id) {
                return user
            }
        }
    }

    getUserByNickname(nickname) {
        for(let user of this.userList) {
            if(nickname == user.nickname) {
                return user
            }
        }
    }

    getClientUser() {
        let nickname = document.querySelector('#header-right > span > a').textContent
        return this.getUserByNickname(nickname)
    }

    // Other functions
    addUser(id, nickname, image) {
        this.userList.push(new User(
            id,
            nickname,
            image
        ))
    }

    clear() {
        this.userList = []
    }
}
