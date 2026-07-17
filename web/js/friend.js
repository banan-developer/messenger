const App2 = {
            data() {
                return {
                    avatarURL: "",
                    name: "",
                    sex: "",
                    group: "",
                    aboutPerson: "",
                    wall: [],
                    friends: [],
                    userID: null,
                    ChatID: null,
                    lightboxImg: null,
                    SessionsID: null,
                    friendshipStatus: "none",
                    friendshipLoading: true

                }
            },
            computed: {
                isOwnProfile() {
                    return this.SessionsID !== null && Number(this.SessionsID) === Number(this.userID)
                },
                friendshipLabel() {
                    if (this.friendshipStatus === "friend") return "Удалить из друзей"
                    if (this.friendshipStatus === "outgoing") return "Отменить заявку"
                    return "Добавить в друзья"
                },
                friendshipIcon() {
                    if (this.friendshipStatus === "friend") return "fas fa-user-minus"
                    if (this.friendshipStatus === "outgoing") return "fas fa-xmark"
                    return "fas fa-user-plus"
                }
            },
            mounted(){
                const url = new URLSearchParams(window.location.search);
                this.userID = url.get('id');
                if (!this.userID) {
                this.error = "ID пользователя не указан";
                this.isLoading = false;
                return;
            }
                this.getName()
                this.getWall()
                this.loadFriend()
                this.Profile()
                this.loadFriendshipStatus()
            },
            methods: {
               async getName(){
                    try{
                        const res = await fetch( `/api/profile?id=${this.userID}`)
                        if (!res.ok) throw new Error("ошибка загрузки имени")
                        const data = await res.json()
                        this.name = data.name
                        this.aboutPerson = data.about
                        this.avatarURL = data.avatar
                        this.sex = data.sex
                        this.group = data.group
                    }catch(err){
                        console.log(err)
                    }
                },
                async Profile(){
                   try{
                        const res = await fetch("/api/profile")
                        if (!res.ok) throw new Error("ошибка загрузки айди пользователя")
                        const data = await res.json()
                        this.SessionsID = data.id
                   }catch(err){
                        console.log(err)
                   }
                },
                async getWall(){
                    try{
                        const res = await fetch(`/api/post?user_id=${this.userID}`)
                        if (!res.ok) throw new Error("ошибка загрузки стены")
                        const data = await res.json()
                        this.wall = data

                    }catch(err){
                        console.log(err)
                    }
                },

                async loadFriend(){
                    try{
                        const res = await fetch(`/api/friend?user_id=${this.userID}`)

                        if (!res.ok) throw new Error("Ошибка отправления id пользователя для последующего его добавления в друзья")
                        const data = await res.json()
                        this.friends = data
                        
                    }catch(err){
                        console.log(err)
                    }
                },

                async loadFriendshipStatus(){
                    this.friendshipLoading = true
                    try {
                        const [friendsRes, outgoingRes] = await Promise.all([
                            fetch('/api/friend'),
                            fetch('/api/outgoingrequest')
                        ])
                        if (!friendsRes.ok || !outgoingRes.ok) {
                            throw new Error('Не удалось определить статус дружбы')
                        }

                        const friends = await friendsRes.json()
                        const outgoing = await outgoingRes.json()
                        const profileID = Number(this.userID)

                        if (friends.some(friend => Number(friend.id) === profileID)) {
                            this.friendshipStatus = 'friend'
                        } else if (outgoing.some(request => Number(request.id) === profileID)) {
                            this.friendshipStatus = 'outgoing'
                        } else {
                            this.friendshipStatus = 'none'
                        }
                    } catch (err) {
                        console.log(err)
                    } finally {
                        this.friendshipLoading = false
                    }
                },

                async changeFriendship(){
                    if (this.friendshipLoading || this.isOwnProfile) return

                    this.friendshipLoading = true
                    try {
                        let endpoint = `/api/friend?id=${encodeURIComponent(this.userID)}`
                        let method = 'POST'

                        if (this.friendshipStatus === 'friend') {
                            method = 'DELETE'
                        } else if (this.friendshipStatus === 'outgoing') {
                            endpoint = `/api/outgoingrequest?id=${encodeURIComponent(this.userID)}`
                            method = 'DELETE'
                        }

                        const res = await fetch(endpoint, { method })
                        if (!res.ok) throw new Error('Не удалось изменить статус дружбы')

                        this.friendshipStatus = method === 'POST' ? 'outgoing' : 'none'
                    } catch (err) {
                        console.log(err)
                    } finally {
                        this.friendshipLoading = false
                    }
                },
                

                goToFriendProfile(friendId) {
                    if (friendId != this.SessionsID){
                        window.location.href = `/friend?id=${friendId}`
                    }
                    else{
                        console.log("На свой акканут нельзя перейти")
                    }
                },

                goBack(){
                    window.location.href = "/profile"
                },

                openLightbox(imgUrl) {
                    this.lightboxImg = imgUrl;
                },

                closeLightbox() {
                    this.lightboxImg = null;
                },

                async exitFromAccount(){
                    try{
                        const res = await fetch("/exit")
                        if (!res.ok) throw new Error("Ошибка выхода с аккаунта")
                        window.location.href = "/login"

                    }catch(err){
                        console.log(err)
                    }
                },

                async loadFriendMessage() {
                    const url = new URLSearchParams(window.location.search);
                    this.ChatID = url.get('id');
                    window.location.href = `/chat?id=${this.ChatID}`;
                },

            }
        }
    Vue.createApp(App2)
        .mount('#VUE2')
