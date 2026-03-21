const App = {
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
                    wall: [],

                    friends: [],
                    searchname: "",
                    foundUsers: [],
                    inputValueTitle: "",
                    inputValueText: "",

                    valueEditTitle: "",
                    valueEditText: "",
                    editingID: null,
                    test: 2,

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
                        const res = await fetch("/api/profile", {
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
                openEditPerson(){
                    this.valueEditName = this.name
                    this.valueEditAboutPers = this.aboutPerson
                    document.getElementById('myDialog').showModal()
                },

                closeEditPerson(){
                    document.getElementById("myDialog").close()
                },
            
                async PushEditedProfile(){
                    if(this.valueEditName == ""){
                        alert("Имя не может быть пустым!")
                        return
                    }
                        
                    try{
                        const res = await fetch("/api/profile",{
                            method: "PUT",
                            credentials: "same-origin",
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify({
                                name: this.valueEditName,
                                about: this.valueEditAboutPers,
                            })
                        })
                         if (!res.ok) {
                            throw new Error(`HTTP error! status: ${res.status}`);
                        }
                        const data = await res.json();
                        console.log('Успешно обновлено:', data);

                    }catch(err){
                        console.log("Ошибка при отправки, измененных данных на сервер")
                    }
                    document.getElementById('myDialog').close()
                    await this.getName()
                },
// ........................................................ СТЕНА .............................................................
                // отправка данных стены на сервер
                async pushWalltoServer(){
                    if (this.inputValueText == '' || this.inputValueTitle == ''){
                        alert("Нельзя опубликовать пустой пост")
                        return
                    }
                    try{
                        const res = await fetch("/api/post", {
                            method: "POST",
                            headers: {'Content-Type': 'application/json'},
                            credentials: "same-origin",
                            body: JSON.stringify({
                                title: this.inputValueTitle,
                                text: this.inputValueText
                            })
                        })
                        await this.getWall();
                        this.inputValueText = ""
                        this.inputValueTitle = ""
                       

                    }catch(err){
                        console.log(err)
                    }
                },
               // получение данных стены с сервера
                async getWall(){
                    try{
                        const res = await fetch("/api/post",{
                            method: "GET",
                            credentials: "same-origin"
                        })
                        if (!res.ok) throw new Error("ошибка загрузки стены")
                        const data = await res.json()
                        this.wall = data

                    }catch(err){
                        console.log(err)
                    }
                },
                // Удаление записей из стены
                async deleteWall(wall_ID){
                    try{
                        const res = await fetch(`/api/post?id=${wall_ID}`, {
                            method: "DELETE",
                            credentials: "same-origin"
                        })
                        if (!res.ok) throw new Error('Ошибка удаления')
                        await this.getWall()

                    }catch(err){
                        console.log(err)
                    }
                },


                closeEditWall(){
                    document.getElementById("myDialog2").close()
                    document.getElementById("friends").close()
                },
                // редактирования записей на стене
                async editWall(wall_ID){
                    document.getElementById("myDialog2").showModal()
                    try{
                        const res = await fetch(`/api/post?id=${wall_ID}`, {
                            credentials: "same-origin",
                            method: "GET"
                        })
                        if (!res.ok) throw new Error("ошибка редактирования поста")
                        const data = await res.json()
                        this.valueEditTitle = data.title
                        this.valueEditText = data.text
                        this.editingID = wall_ID

                    }catch(err){
                        console.log(err)
                    }
                },

                async updateEditingWall(wall_ID){
                    document.getElementById("myDialog2").close()
                    try{
                        const res = await fetch(`/api/post?id=${wall_ID}`, {
                            credentials: "same-origin",
                            method: "PUT",
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify({title: this.valueEditTitle, text: this.valueEditText})
                        })
                        if (!res.ok) throw new Error("ошибка редактирования поста")

                    }catch (err){
                        console.log(err)
                    }
                    await this.getWall()
                },
                // загрузка фотки пользователя
                async uploadAvatar(event){
                    const file = event.target.files[0]
                    if (!file) return
                    const formData = new FormData()
                    formData.append("avatar", file)
                    try{
                        const res = await fetch("/api/profile/avatar",{
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

                async exitFromAccount(){
                    try{
                        const res = await fetch("/exit")
                        if (!res.ok) throw new Error("Ошибка выхода с аккаунта")
                        window.location.href = "/login"

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
                        const res = await fetch(`/api/friend?name=${this.searchname}`)
                        if (!res.ok) throw new Error("Ошибка получения данных пользователя, при его поиске")
                        const data = await res.json()
                        this.foundUsers = data

                    }catch(err){
                        console.log(err)
                    }
                },
                // Добавление в друзья
                async addToFriend(idUser){
                    try{
                        const res = await fetch(`/api/friend?id=${idUser}`,{
                            method: "POST"
                        })
                        if (!res.ok) throw new Error("Ошибка отправления id пользователя для последующего его добавления в друзья")
                
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
                async loadFriendProfile(friendID) {
                    window.location.href = `/anotherProfile.html?id=${friendID}`;
                },
                openDialogInfo(){
                    document.getElementById("myDialogInfo").showModal()
                },
                exitNtfct(){
                    document.getElementById("myDialogInfo").close()
                }
            }
        }
    Vue.createApp(App).mount('#VUE')