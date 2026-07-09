const MessagesApp = {
    data() {
        return {
            userName: "asdas",
            userAvatar: "",
            globalSearch: "",
            chatSearch: "",
            chats: []
        }
    },
    mounted(){
                this.GetListChatsWithLastMessage()
            },
    computed: {
        filteredChats() {
            const query = this.chatSearch.trim().toLowerCase()
        
            let result = this.chats.filter(c => c.last_message && c.last_message.trim() !== '')
        
            if (query) {
            result = result.filter(c => c.name.toLowerCase().includes(query))
            }
            return result
            }
    },
    methods: {
        openChat(chat) {
            window.location.href = `/chat.html?id=${chat.id}`
        },
        exitFromAccount() {
            window.location.href = "/login"
        },
        async GetListChatsWithLastMessage(){
            try{
                const res = await fetch("/api/messages",{
                method: "GET",
                credentials: "same-origin",
                })
                if (!res.ok) throw new Error("Ошибка при получении списка последних сообщений")
                const data = await res.json()
                this.chats = data
            }catch(err){
                console.log(err)
            }
        },
        goToChatByFrinedId(ChatID){
            window.location.href = `/chat?id=${ChatID}`
        }
    }
}

Vue.createApp(MessagesApp)
    .component('app-sidebar', AppSidebar)
    .mount('#MESSAGES')
