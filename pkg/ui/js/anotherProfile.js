const App2 = {
            data() {
                return {
                    avatarURL: "",
                    name: "",
                    sex: "",
                    aboutPerson: "",
                    imgURL: "",
                    valueEditName: "",
                    valueEditAboutPers: "",

                    // стена
                    wall: [
                        {id: 1 ,title: "Первая запись на стене!", text: "Сегодня я впервые написал на своей стене, классный момент"},
                        {id: 2, title: "абоба", text: "стенаааааааааааа"} 
                    ],

                    friends: [],
                    searchname: "",
                    foundUsers: [],
                    inputValueTitle: "",
                    inputValueText: "",

                    valueEditTitle: "",
                    valueEditText: "",
                    editingID: null,

                }
            },
            mounted(){
                this.getName()
                this.getWall()
                this.loadFriend()
            },
            methods: {
               async getName(){
                    try{
                        const res = await fetch("/test", {
                            credentials: "same-origin"
                        })
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

// ........................................................ СТЕНА .............................................................
               // получение данных стены с сервера
                async getWall(){
                    try{
                        const res = await fetch("/pushWall",{
                            credentials: "same-origin"
                        })
                        if (!res.ok) throw new Error("ошибка загрузки стены")
                        const data = await res.json()
                        this.wall = data

                    }catch(err){
                        console.log(err)
                    }
                },

                // загрузка фотки пользователя
                async uploadAvatar(event){
                    const file = event.target.files[0]
                    if (!file) return
                    const formData = new FormData()
                    formData.append("avatar", file)
                    try{
                        const res = await fetch("/uploadAvatar",{
                            credentials: "same-origin",
                            method: "POST",
                            body: formData
                        })
                        if (!res.ok) throw new Error("ошибка редактирования поста")
                        const data = await res.json()
                        this.avatarURL = data.avatar

                    }catch(err){
                        console.log(err)
                    }
                },

                async OpenSearchFriend(){
                    document.getElementById("friends").showModal()
                },
                async SearchFriends(){
                    if (this.searchname == ''){
                            alert("Введи имя пользователя")
                            return
                        }
                    try{
                        const res = await fetch(`/searchFriends?name=${this.searchname}`)
                        if (!res.ok) throw new Error("Ошибка получения данных пользователя, при его поиске")
                        const data = await res.json()
                        this.foundUsers = data

                    }catch(err){
                        console.log(err)
                    }
                },

                async loadFriend(){
                    try{
                        const res = await fetch("/loadFriend", {
                            method: "GET"
                        })

                        if (!res.ok) throw new Error("Ошибка отправления id пользователя для последующего его добавления в друзья")
                        const data = await res.json()
                        this.friends = data
                        
                    }catch(err){
                        console.log(err)
                    }
                },

                async loadFriendProfile(friendID){
                    try{
                        const res = await fetch(`/loadFriendProfile?id=${friendID}`,{
                            method: "POST"
                        }
                        )

                    }catch(err){
                        console.log(err)
                    }
                }

     
                
            }
        }
    Vue.createApp(App2).mount('#VUE2')