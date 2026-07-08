const AppSidebar = {
    props: {
        active: { type: String, default: '' },
        userName: { type: String, default: '' },
        userAvatar: { type: String, default: '' }
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

            <div class="sidebar-bottom">
                <div class="mini-profile">
                    <img :src="userAvatar" alt="avatar" class="mini-avatar">
                    <span class="mini-name">{{ userName }}</span>
                </div>
                <button class="sidebar-logout" @click="$emit('logout')">
                    <span class="sidebar-icon">🚪</span>
                    <span>Выход</span>
                </button>
            </div>
        </aside>
    `
}
