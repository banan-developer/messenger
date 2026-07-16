const FriendsApp = {
    data() {
        return {
            userName: "asdas",
            userAvatar: "",
            currentUserID: null,
            activeTab: "all",
            searchQuery: "",
            openMenuId: null,
            friends: [],
            requests: [],
            showAddFriend: false,
            friendSearchQuery: '',
            friendSearchResults: [],
            friendSearchLoading: false,
            friendSearchTimer: null,
            sentFriendRequestIDs: []
        }
    },
     async mounted(){
            await Promise.all([
                this.loadCurrentUser(),
                this.loadFriend(),
                this.GetFrienedRequest()
            ])
        },
    computed: {
        filteredFriends() {
            const query = this.searchQuery.trim().toLowerCase()
            if (!query) return this.friends
            return this.friends.filter(f => f.name.toLowerCase().includes(query))
        }
    },
    methods: {
        loadSentFriendRequestIDs() {
            try {
                const saved = JSON.parse(localStorage.getItem(`sentFriendRequestIDs:${this.currentUserID}`) || '[]')
                this.sentFriendRequestIDs = Array.isArray(saved)
                    ? saved.map(Number).filter(Number.isInteger)
                    : []
            } catch {
                this.sentFriendRequestIDs = []
            }
        },
        saveSentFriendRequestIDs() {
            localStorage.setItem(
                `sentFriendRequestIDs:${this.currentUserID}`,
                JSON.stringify(this.sentFriendRequestIDs)
            )
        },
        async loadCurrentUser() {
            try {
                const res = await fetch('/api/profile', { credentials: 'same-origin' })
                if (!res.ok) throw new Error('Не удалось загрузить профиль')
                const user = await res.json()
                this.currentUserID = Number(user.id)
                this.loadSentFriendRequestIDs()
            } catch (err) {
                console.log(err)
            }
        },
        toggleMenu(id) {
            this.openMenuId = this.openMenuId === id ? null : id
        },
        openChatWith(friend) {
            window.location.href = `/chat?id=${friend.id}`
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

        friendInitials(name) {
            if (!name) return '?'
            return name
                .trim()
                .split(/\s+/)
                .slice(0, 2)
                .map(part => part[0])
                .join('')
                .toUpperCase()
        },
        openAddFriend() {
            this.showAddFriend = true
            this.friendSearchQuery = ''
            this.friendSearchResults = []
        },
        closeAddFriend() {
            this.showAddFriend = false
        },
        onSearchFriends() {
            clearTimeout(this.friendSearchTimer)
            const query = this.friendSearchQuery
            if (!query) {
                this.friendSearchResults = []
                return
            }
            this.friendSearchTimer = setTimeout(() => this.searchFriends(query), 300)
        },
        async searchFriends(query) {
            this.friendSearchLoading = true;
            try {
                const res = await fetch(`/api/friend?name=${encodeURIComponent(query)}`, {
                credentials: 'same-origin'
                });
                if (!res.ok) throw new Error('Не удалось выполнить поиск');
                const results = await res.json();
                const friendIDs = new Set(this.friends.map(friend => Number(friend.id)))
                const sentRequestIDs = new Set(this.sentFriendRequestIDs)
                this.friendSearchResults = results
                    .filter(person => Number(person.id) !== this.currentUserID)
                    .filter(person => !friendIDs.has(Number(person.id)))
                    .map(person => ({
                        ...person,
                        _added: sentRequestIDs.has(Number(person.id)),
                        _adding: false
                    }));
            } catch (err) {
                console.log(err);
                this.friendSearchResults = [];
            } finally {
                this.friendSearchLoading = false;
            }
        },

        async addFriend(person) {
        if (person._added || person._adding) return
        person._adding = true
        try {
            const res = await fetch(`/api/friend?id=${encodeURIComponent(person.id)}`, {
            method: 'POST',
            credentials: 'same-origin'
            });
            if (!res.ok) throw new Error('Не удалось добавить друга');
            person._added = true;
            const personID = Number(person.id)
            if (!this.sentFriendRequestIDs.includes(personID)) {
                this.sentFriendRequestIDs.push(personID)
                this.saveSentFriendRequestIDs()
            }
            await this.loadFriend();
        } catch (err) {
            console.log(err);
        } finally {
            person._adding = false
        }
        },
        goToProfile(friend){
            window.location.href = `/friend?id=${friend}`
        }
    }
}

Vue.createApp(FriendsApp)
    .mount('#FRIENDS')
