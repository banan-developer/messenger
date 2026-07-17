const App3 = {
    data() {
        return {
            name: "",
            avatar_url: "",
            userGroup: "",
            friendID: null,
            currentUserID: null,
            messanges: [],
            editingMes: "",
            socket: null,
            inputText: "",
			selectedFile: null,
			selectedFileName: "",
			isUploading: false,
			groupMembers: [], groupFriends: [], groupMemberID: '',
			groupTitle: '', isGroupSettingsOpen: false,
            editingMessage: null,
            inputEditingText: "",
            isediting: 0,
            isLoading: true,
            error: null,
            lightboxImg: null
        }
    },
    async mounted() {
        const url = new URLSearchParams(window.location.search);
        this.friendID = Number(url.get('id'));
        this.chatID = Number(url.get('chat_id'));
        if (this.chatID) this.friendID = 0;
        if (!this.friendID && !this.chatID) {
            this.error = "ID пользователя не указан";
            this.isLoading = false;
            return;
        }
        
        try {
            await this.getCurrentId()
            await this.GetName()
            this.connectSocket()
            await this.GetAllMessage()
			if (this.chatID) await this.loadGroupMembers()
        } catch (err) {
            console.error('Ошибка загрузки чата:', err)
            this.error = 'Ошибка загрузки чата'
        } finally {
            this.isLoading = false
        }
    },
    methods: {
        async GetName() {
            try {
                if (this.chatID) {
                    const res = await fetch('/api/groups', { credentials: 'same-origin' })
                    if (!res.ok) throw new Error("ошибка загрузки группы")
                    const groups = await res.json()
                    const group = groups.find(item => Number(item.id) === Number(this.chatID))
                    if (!group) throw new Error("группа не найдена")

                    this.name = group.title || 'Групповая беседа'
                    this.groupTitle = group.title || ''
                    this.avatar_url = group.avatar_url || group.avatar || ''
                    return
                }

                const res = await fetch(`/api/profile?id=${this.friendID}`)
                if (!res.ok) throw new Error("ошибка загрузки имени")
                const data = await res.json()
                this.name = data.name
                this.avatar_url = data.avatar || '/static/icons/def_picture.png'
                this.userGroup = data.group || ''
            } catch (err) {
                console.log(err)
            }
        },
        getBack() {
            window.location.href = "/messages"
        },
        openLightbox(imgUrl) {
            this.lightboxImg = imgUrl;
        },
        closeLightbox() {
            this.lightboxImg = null;
        },
        async loadFriendProfile() {
            if (this.chatID) return
            window.location.href = `/friend?id=${this.friendID}`;
        },
        connectSocket() {
            const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
            const wsHost = window.location.host
            const target = this.chatID ? `chat_id=${this.chatID}` : `id=${this.friendID}`
            this.socket = new WebSocket(`${wsProtocol}//${wsHost}/ws?${target}&user_id=${this.currentUserID}`)

            this.socket.onopen = () => {
                console.log("ws подключен")
            }

            this.socket.onmessage = (event) => {
                console.log("RAW:", event.data)
                try {
                    const data = JSON.parse(event.data)
                    console.log("PARSED:", data)

                    if (data.event === 'new_message') {
                        this.upsertMessage(data.message)
                    }

                    if (data.event === 'message_updated') {
                        const index = this.messanges.findIndex(m => m.id === data.data.id)
                        if (index !== -1) {
                            this.messanges[index].text = data.data.text
                        }
                    }

                    if (data.event === 'message_deleted') {
                        this.messanges = this.messanges.filter(m => m.id !== data.data.id)
                    }
                } catch (err) {
                    console.error('Ошибка парсинга:', err)
                }
            }

            this.socket.onclose = () => {
                console.log("ws закрыт")
            }

            this.socket.onerror = (error) => {
                console.error("WebSocket ошибка:", error)
            }
        },
        sendMessage() {
            if (this.selectedFile) {
                this.sendAttachment()
                return
            }
            if (!this.inputText.trim()) return
            
            this.socket.send(JSON.stringify({
                text: this.inputText.trim()
            }))
            this.inputText = ""
        },
		selectFile(event) {
			const [file] = event.target.files
			if (!file) {
				this.clearSelectedFile()
				return
			}
			const extension = `.${file.name.split('.').pop().toLowerCase()}`
			const allowed = ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.pdf', '.docx', '.xml']
			if (!allowed.includes(extension)) {
				this.error = 'Разрешены изображения, PDF, DOCX и XML'
				this.clearSelectedFile()
				return
			}
			if (file.size > 20 * 1024 * 1024) {
				this.error = 'Размер файла не должен превышать 20 МБ'
				this.clearSelectedFile()
				return
			}
			this.error = null
			this.selectedFile = file || null
			this.selectedFileName = file ? file.name : ""
		},
		clearSelectedFile() {
			this.selectedFile = null
			this.selectedFileName = ''
			if (this.$refs.fileInput) this.$refs.fileInput.value = ''
		},
		async sendAttachment() {
			if (!this.selectedFile || this.isUploading) return

			this.isUploading = true
			this.error = null
			try {
				const formData = new FormData()
				if (this.chatID) formData.append("chat_id", this.chatID)
				else formData.append("friend_id", this.friendID)
				formData.append("text", this.inputText.trim())
				formData.append("file", this.selectedFile)

				const res = await fetch("/api/messages/file", {
					method: "POST",
					body: formData
				})
				if (!res.ok) throw new Error((await res.text()) || "Не удалось отправить файл")
				const message = await res.json()
				this.upsertMessage({
					...message,
					from_id: this.currentUserID
				})

				this.inputText = ""
				this.clearSelectedFile()
			} catch (err) {
				console.error(err)
				this.error = err.message || "Не удалось отправить файл"
			} finally {
				this.isUploading = false
			}
		},
		upsertMessage(message) {
			const normalizedMessage = {
				id: message.id,
				text: message.text || '',
				from_id: message.from_id,
				to_id: message.to_id || 0,
				chat_id: message.chat_id || 0,
				created_at: message.created_at,
				attachment_url: message.attachment_url || '',
				attachment_name: message.attachment_name || '',
				attachment_type: message.attachment_type || '',
				attachment_size: Number(message.attachment_size) || 0
			}
			const index = this.messanges.findIndex(item => item.id === normalizedMessage.id)
			if (index === -1) {
				this.messanges.push(normalizedMessage)
			} else {
				this.messanges[index] = normalizedMessage
			}
			this.scrollToBottom()
		},
		isImageAttachment(message) {
			if (message.attachment_type) return message.attachment_type.startsWith('image/')
			return /\.(jpe?g|png|gif|webp)$/i.test(message.attachment_url || '')
		},
		fileNameFromURL(url) {
			const name = (url || '').split('/').pop() || 'Файл'
			try { return decodeURIComponent(name) } catch { return name }
		},
		fileIcon(message) {
			const name = (message.attachment_name || message.attachment_url || '').toLowerCase()
			if (name.endsWith('.pdf')) return 'fas fa-file-pdf'
			if (name.endsWith('.docx')) return 'fas fa-file-word'
			if (name.endsWith('.xml')) return 'fas fa-file-code'
			return 'fas fa-file'
		},
		fileTypeLabel(message) {
			const name = (message.attachment_name || message.attachment_url || '').toLowerCase()
			if (name.endsWith('.pdf')) return 'PDF'
			if (name.endsWith('.docx')) return 'DOCX'
			if (name.endsWith('.xml')) return 'XML'
			return 'Файл'
		},
		formatFileSize(bytes) {
			const size = Number(bytes) || 0
			if (size < 1024) return `${size} Б`
			if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} КБ`
			return `${(size / 1024 / 1024).toFixed(1)} МБ`
		},
		async loadGroupMembers(){ const r=await fetch(`/api/groups?chat_id=${this.chatID}`); if(r.ok)this.groupMembers=await r.json(); const f=await fetch('/api/friend'); if(f.ok)this.groupFriends=await f.json() },
		async addGroupMember(){ const r=await fetch(`/api/groups?chat_id=${this.chatID}&user_id=${this.groupMemberID}`,{method:'PUT'}); if(r.ok){this.groupMemberID='';this.loadGroupMembers()} },
		async removeGroupMember(id){ const r=await fetch(`/api/groups?chat_id=${this.chatID}&user_id=${id}`,{method:'DELETE'}); if(r.ok)this.loadGroupMembers() },
		async renameGroup(){ await fetch(`/api/groups?chat_id=${this.chatID}`,{method:'PATCH',headers:{'Content-Type':'application/json'},body:JSON.stringify({title:this.groupTitle})}); this.isGroupSettingsOpen=false },
		async deleteGroup(){ if(!confirm('Удалить беседу?'))return; const r=await fetch(`/api/groups?chat_id=${this.chatID}`,{method:'DELETE'}); if(r.ok)window.location.href='/messages' },
		messageAuthor(id) { return this.groupMembers.find(member => Number(member.id) === Number(id)) || null },
        async getCurrentId() {
            try {
                const res = await fetch("/api/profile")
                if (!res.ok) throw new Error("ошибка загрузки профиля")
                const data = await res.json()
                this.currentUserID = data.id
                console.log("CurrentID:", this.currentUserID)
            } catch (err) {
                console.log("Current id error", err)
            }
        },
        async GetAllMessage() {
            try {
                const target = this.chatID ? `chat_id=${this.chatID}` : `id=${this.friendID}`
                const res = await fetch(`/api/messages?${target}&user_id=${this.currentUserID}`)
                if (!res.ok) throw new Error("Ошибка получения сообщений")
                const data = await res.json()
                this.messanges = data
                this.scrollToBottom()
            } catch (err) {
                console.log(err)
            }
        },
        formatDate(dateString) {
            if (!dateString) return '';

            // Если дата уже пришла обрезанной из БД в формате "MM-DD HH:MM"
            if (/^\d{2}-\d{2}\s\d{2}:\d{2}$/.test(dateString)) {
                return dateString;
            }

            // Для новых сообщений заменяем пробел на 'T', 
            // чтобы формат "YYYY-MM-DD HH:MM:SS" без ошибок работал в Safari/iOS
            const safeDateString = dateString.replace(' ', 'T');
            const date = new Date(safeDateString);

            // Если парсинг всё равно не удался, возвращаем исходную строку, чтобы не было NaN
            if (isNaN(date.getTime())) {
                return dateString;
            }

            const month = String(date.getMonth() + 1).padStart(2, '0');
            const day = String(date.getDate()).padStart(2, '0');
            const hours = String(date.getHours()).padStart(2, '0');
            const minutes = String(date.getMinutes()).padStart(2, '0');
            
            return `${month}-${day} ${hours}:${minutes}`;
        },
        async editMessHttpFallback() {
            try {
                this.isediting = 1
                document.getElementById("dialog-msg-tool").close()
                const res = await fetch(`/api/message?msg_id=${this.editingMessage}`)
                if (!res.ok) throw new Error("Ошибка получения сообщения")
                const data = await res.json()
                this.editingMes = data.text
                this.inputEditingText = this.editingMes
            } catch (err) {
                console.log(err)
            }
        },
        async SaveEditingMessHttpFallback() {
            this.isediting = 0
            try {
                const res = await fetch(`/api/message`, {
                    method: "PUT",
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ id: this.editingMessage, text: this.inputEditingText })
                })
                if (!res.ok) throw new Error("Ошибка редактирования")
                await this.GetAllMessage()
            } catch (err) {
                console.log(err)
            }
        },
        async deleteMessHttpFallback() {
            try {
                document.getElementById("dialog-msg-tool").close()
                const res = await fetch(`/api/message?id=${this.editingMessage}`, {
                    method: "DELETE"
                })
                if (!res.ok) throw new Error("Ошибка удаления")
                await this.GetAllMessage()
            } catch (err) {
                console.log(err)
            }
        },
        editMess() {
            const message = this.messanges.find(item => item.id === this.editingMessage)
            if (!message) return

            this.isediting = 1
            this.inputEditingText = message.text
            document.getElementById("dialog-msg-tool").close()
        },
        SaveEditingMess() {
            const text = this.inputEditingText.trim()
            if (!text || !this.isSocketOpen()) return

            this.socket.send(JSON.stringify({
                event: "edit_message",
                id: this.editingMessage,
                text
            }))
            this.isediting = 0
        },
        deleteMess() {
            document.getElementById("dialog-msg-tool").close()
            if (!this.isSocketOpen()) return

            this.socket.send(JSON.stringify({
                event: "delete_message",
                id: this.editingMessage
            }))
        },
        isSocketOpen() {
            if (this.socket?.readyState === WebSocket.OPEN) return true
            this.error = "Соединение с чатом ещё не установлено"
            return false
        },
        openMsgTools(msg) {
            this.editingMessage = msg
            const dialog = document.getElementById("dialog-msg-tool")
            dialog.showModal()
        },
        closeMsgTools() {
            document.getElementById("dialog-msg-tool").close()
        },
        scrollToBottom() {
            const container = document.querySelector('.messanges')
            if (container) {
                setTimeout(() => {
                    container.scrollTop = container.scrollHeight
                }, 100)
            }
        }
    }
}

Vue.createApp(App3).mount('#VUECHAT')
