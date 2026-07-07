const App3 = {
            data() {
                return {
                    name: "",
                    avatar_url: "",
                    friendID: null,
                    currentUserID: null,
                    messanges: [],
                    editingMes: "",
                    socket: null,
                    inputText: "",
                    editingMessage: null,
                    inputEditingText: "",
                    isediting: 0,

                }
            },
            async mounted(){
                const url = new URLSearchParams(window.location.search);
                this.friendID = Number(url.get('id'));
                if (!this.friendID) {
                    this.error = "ID пользователя не указан";
                    this.isLoading = false;
                    return;
                }
                
                await this.GetName()
                await this.getCurrentId()
                this.connectSocket()
                await this.GetAllMessage()
            },
            methods: {
                async GetName(){
                    try{
                        const res = await fetch(`/api/profile?id=${this.friendID}`)
                        if (!res.ok) throw new Error("ошибка загрузки имени")
                        const data = await res.json()
                        this.name = data.name
                        this.avatar_url = data.avatar

                    }catch(err){
                        console.log(err)
                    }
                },
                getBack(){
                    window.location.href = "/profile"
                },
                async loadFriendProfile() {
                    window.location.href = `/friend?id=${this.friendID}`;
                },

                connectSocket(){
                    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
                    const wsHost = window.location.host
                    this.socket = new WebSocket(`${wsProtocol}//${wsHost}/ws?id=${this.friendID}&user_id=${this.currentUserID}`)
                    this.socket.onopen = () => {
                        console.log("ws подлючен")
                    }
                    this.socket.onmessage = (event) => {
                        console.log("RAW:", event.data)
                        const msg = JSON.parse(event.data)
                        console.log("PARSED:", msg)
                        this.messanges.push({
                            id: msg.id,
                            text: msg.text,
                            from_id: msg.from_id,
                            to_id: msg.chats_id,
                            created_at: msg.created_at
                        })
                        
                    }
                    this.socket.onclose = () => {
                        console.log("ws закрыт")
                    }
                },
                sendMessage(){
                    const tempId = Date.now()
                    if (!this.inputText) return
                    this.socket.send(JSON.stringify({
                        text: this.inputText
                    }))
                    this.inputText = ""
                },
                async getCurrentId(){
                     try{
                        const res = await fetch("/api/profile")
                        if (!res.ok) throw new Error("ошибка загрузки айди")
                        const data = await res.json()
                        this.currentUserID = data.id
                        console.log("CurrentID:", this.currentUserID)

                    }catch(err){
                        console.log("Current id error", err)
                    }
                },
                async GetAllMessage(){
                    try{
                        const res = await fetch(`/api/messages?id=${this.friendID}&user_id=${this.currentUserID}`)
                        if (!res.ok) throw new Error("Ошибка получения сообщения")
                        const data = await res.json()
                        this.messanges = data
                    }catch(err){
                        console.log(err)
                    }
                },
                formatDate(dateString) {
                    if (!dateString) return ''
                    const date = new Date(dateString)
                    const month = String(date.getMonth() + 1).padStart(2, '0')
                    const day = String(date.getDate()).padStart(2, '0')
                    const hours = String(date.getHours()).padStart(2, '0')
                    const minutes = String(date.getMinutes()).padStart(2, '0')
                    return `${month}-${day} ${hours}:${minutes}`
                },

                async editMess(){
                    try{
                        this.isediting = 1
                        document.getElementById("dialog-msg-tool").close()
                        const res = await fetch(`/api/message?msg_id=${this.editingMessage}`)
                        if (!res.ok) throw new Error("Ошибка получения сообщения")
                        const data = await res.json()
                        this.editingMes = data.text
                        this.inputEditingText = this.editingMes
                    }catch(err){
                        console.log(err)
                    }
                },
                async SaveEditingMess(){
                    this.isediting = 0
                    const res = await fetch(`/api/message`,{
                        method: "PUT",
                        headers: {'Content-Type': 'application/json'},
                        body: JSON.stringify({id: this.editingMessage, text: this.inputEditingText})
                    })
                    await this.GetAllMessage()
                },

                async deleteMess(){
                    try{
                        document.getElementById("dialog-msg-tool").close()
                        const res = await fetch(`/api/message?id=${this.editingMessage}`,{
                            method: "DELETE"
                        })
                        await this.GetAllMessage()
                        
                    }catch(err){
                        console.log(err)
                    }
                },
                openMsgTools(msg){
                    this.editingMessage = msg
                    dialog = document.getElementById("dialog-msg-tool")
                    dialog.showModal()
                },
                closeMsgTools(){
                    document.getElementById("dialog-msg-tool").close()
                }
            }
        }
    Vue.createApp(App3).mount('#VUECHAT')