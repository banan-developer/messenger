const MessagesApp = {
    data() {
        return {
            userName: "asdas",
            userAvatar: "",
            globalSearch: "",
            chatSearch: "",
            chats: [
                { id: 1, name: "Алексей Петров", initials: "АП", lastMessage: "Ок, спасибо!", time: "Вчера", unread: 0, isGroup: false },
                { id: 2, name: "Группа 1442", initials: "👥", lastMessage: "Пётр: Кто завтра в библиотеку?", time: "10:15", unread: 0, isGroup: true },
                { id: 3, name: "Светлана Анисимова", initials: "СА", lastMessage: "Привет! Давай обсудим проект?", time: "10:30", unread: 2, isGroup: false },
                { id: 4, name: "Сергей Михайлов", initials: "СМ", lastMessage: "Не забудь про нашу встречу завтра в 11.", time: "15:45", unread: 0, isGroup: false }
            ]
        }
    },
    computed: {
        filteredChats() {
            const query = this.chatSearch.trim().toLowerCase()
            if (!query) return this.chats
            return this.chats.filter(c => c.name.toLowerCase().includes(query))
        }
    },
    methods: {
        openChat(chat) {
            window.location.href = `/chat.html?id=${chat.id}`
        },
        exitFromAccount() {
            window.location.href = "/login"
        }
    }
}

Vue.createApp(MessagesApp)
    .component('app-sidebar', AppSidebar)
    .mount('#MESSAGES')
