const App2 = {
            data() {
                return {
                    avatarURL: "",
                    name: "",
                    sex: "",
                    aboutPerson: "",
                    wall: [],
                    friends: [],
                    userID: null,
                    ChatID: null,
                    lightboxImg: null,
                    SessionsID: null

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