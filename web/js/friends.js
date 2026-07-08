const FriendsApp = {
    data() {
        return {
            userName: "asdas",
            userAvatar: "",
            activeTab: "all",
            searchQuery: "",
            openMenuId: null,
            friends: [
                { id: 1, name: "Светлана Анисимова", initials: "СА", status: "В сети, 1241" },
                { id: 2, name: "Андрей Абрамов", initials: "АА", status: "Не в сети. Факультет 4" },
                { id: 3, name: "Василий Васильев", initials: "ВВ", status: "В сети, 1441" }
            ],
            requests: [
                { id: 4, name: "Мария Иванова", initials: "МИ", status: "3 общих друга" },
                { id: 5, name: "Дмитрий Орлов", initials: "ДО", status: "1 общий друг" },
                { id: 6, name: "Ольга Петрова", initials: "ОП", status: "5 общих друзей" }
            ]
        }
    },
    computed: {
        filteredFriends() {
            const query = this.searchQuery.trim().toLowerCase()
            if (!query) return this.friends
            return this.friends.filter(f => f.name.toLowerCase().includes(query))
        }
    },
    methods: {
        toggleMenu(id) {
            this.openMenuId = this.openMenuId === id ? null : id
        },
        openChatWith(friend) {
            window.location.href = `/chat.html?id=${friend.id}`
        },
        exitFromAccount() {
            window.location.href = "/login"
        }
    }
}

Vue.createApp(FriendsApp)
    .component('app-sidebar', AppSidebar)
    .mount('#FRIENDS')
