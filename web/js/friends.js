const FriendsApp = {
    data() {
        return {
            userName: "asdas",
            userAvatar: "",
            activeTab: "all",
            searchQuery: "",
            openMenuId: null,
            friends: [],
            requests: []
        }
    },
     mounted(){
            this.loadFriend()
            this.GetFrienedRequest()
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
		async deleteFriend(friendID) {
			if (!confirm('Удалить из друзей?')) return
			const res = await fetch(`/api/friend?id=${friendID}`, { method: 'DELETE' })
			if (res.ok) { this.toggleMenu(null); await this.loadFriend() }
		},
        exitFromAccount() {
            window.location.href = "/login"
        },
        async GetFrienedRequest(){
            try{
                const res =  await fetch("/api/incomingrequest",{
                    method: "GET",
                    credentials: "same-origin",
                })
                if (!res.ok) throw new Error("ошибка получения друзей")
                const data = await res.json()
                this.requests = data
            }catch(err){
                console.log(err)
            }
        },
        async loadFriend(){
            try{
                const res = await fetch("/api/friend", {
                    method: "GET"
                })

                if (!res.ok) throw new Error("Ошибка отправления id пользователя для последующего его добавления в друзья")
                const data = await res.json()
                this.friends = data
                        
            }catch(err){
                console.log(err)
            }
        },
        async SetAcceptStatus(friendID){
            try{
                const res = await fetch(`/api/friend?friendID=${friendID}`,{
                    method: "PUT",
                    credentials: "same-origin",
                })
                if (!res.ok) throw new Error("Принятия в друзья")
                await this.GetFrienedRequest()
                await this.loadFriend() 
            }catch(err){
                console.log(err)
            }
        },
    }
}

Vue.createApp(FriendsApp)
    .component('app-sidebar', AppSidebar)
    .mount('#FRIENDS')
