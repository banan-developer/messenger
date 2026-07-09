const AppSidebar = {
    props: {
        active: { type: String, default: '' }
    },
    data() {
        return {
            currentUser: {
                name: '',
                avatar: ''
            }
        }
    },
    async mounted() {
        await this.loadCurrentUser()
    },
    methods: {
        async loadCurrentUser() {
            try {
                const res = await fetch("/api/profile", {
                    credentials: "same-origin"
                })
                if (!res.ok) throw new Error("Ошибка загрузки профиля")
                const data = await res.json()
                this.currentUser.name = data.name
                this.currentUser.avatar = data.avatar || '/static/avatars/default.jpg'
            } catch (err) {
                console.log(err)
            }
        },
        handleLogout() {
            this.$emit('logout')
        },
        goProfile(){
            window.location.href = "/profile"
        },
    },
    emits: ['logout'],
    template: `
        <aside class="app-sidebar">
            <div class="sidebar-logo">
                <span>🧭</span>
                <span>ГУАП</span>
            </div>

            <ul class="sidebar-menu">
                <li>
                    <a href="/profile" class="sidebar-link" :class="{ 'is-active': active === 'profile' }">
                        <span class="sidebar-icon">👤</span>
                        <span>Профиль</span>
                    </a>
                </li>
                <li>
                    <a href="/messages" class="sidebar-link" :class="{ 'is-active': active === 'messages' }">
                        <span class="sidebar-icon">💬</span>
                        <span>Сообщения</span>
                    </a>
                </li>
                <li>
                    <a href="/friends" class="sidebar-link" :class="{ 'is-active': active === 'friends' }">
                        <span class="sidebar-icon">👥</span>
                        <span>Друзья</span>
                    </a>
                </li>
            </ul>

            <div class="sidebar-bottom" >
                <div class="mini-profile">
                    <img :src="currentUser.avatar" alt="avatar" class="mini-avatar" @click="goProfile">
                    <span  class="mini-name">{{ currentUser.name }}</span>
                </div>
                <button class="sidebar-logout" @click="handleLogout">
                    <span class="sidebar-icon">🚪</span>
                    <span>Выход</span>
                </button>
            </div>
        </aside>
    `
}