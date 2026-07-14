const MessagesApp = {
    data() {
        return {
            userName: "asdas",
            userAvatar: "",
            globalSearch: "",
            chatSearch: "",
            chats: [], friends: [], isGroupModalOpen: false, groupTitle: '', selectedFriendIds: []
        }
    },
    mounted(){
                this.GetListChatsWithLastMessage(); this.loadFriends()
            },
    computed: {
        filteredChats() {
            const query = this.chatSearch.trim().toLowerCase()
        
            let result = this.chats.filter(c => c.is_group || (c.last_message && c.last_message.trim() !== ''))
        
            if (query) {
            result = result.filter(c => c.name.toLowerCase().includes(query))
            }
            return result
            }
    },
    methods: {
        async loadFriends(){ const r=await fetch('/api/friend'); if(r.ok)this.friends=await r.json() },
        async createGroup(){ if(!this.groupTitle.trim()||!this.selectedFriendIds.length)return; const r=await fetch('/api/groups',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({title:this.groupTitle,user_ids:this.selectedFriendIds})}); if(r.ok){const d=await r.json();window.location.href='/chat?chat_id='+d.chat_id} },
        openChat(chat) {
            window.location.href = chat.is_group
                ? `/chat?chat_id=${chat.id}`
                : `/chat?id=${chat.user_id}`
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
            window.location.href = `/chat?chat_id=${ChatID}`
        }
    }
}

Vue.createApp(MessagesApp)
    .mount('#MESSAGES')
