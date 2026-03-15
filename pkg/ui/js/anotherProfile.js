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
                    window.location.href = `/anotherProfile.html?id=${friendId}`;
                },

                goBack(){
                    window.location.href = "/profile"
                }
            }
        }
    Vue.createApp(App2).mount('#VUE2')