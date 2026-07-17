const { createApp } = Vue;

createApp({
  data() {
    return {
      user: {
        id: null,
        name: '',
        about: '',
        sex: '',
        avatar: '',
        group: ''
      },

      posts: [],
      postsLoading: true,

      friends: [],
      friendsLoading: true,

      // wall composer
      newPost: { title: '', text: '', imageFile: null, imageName: '' },
      posting: false,
      postError: '',

      // post editing
      editingPostId: null,
      editPostForm: { title: '', text: '' },
      editPostError: '',

      // edit profile modal
      showEditProfile: false,
      editProfileForm: { name: '', about: '', group: '' },
      editProfileError: '',
      savingProfile: false,

      // mobile burger menu
      mobileMenuOpen: false,

      // lightbox
      lightboxImg: null
    };
  },

  computed: {
    initials() {
      return this.friendInitials(this.user.name);
    }
  },

  mounted() {
    this.fetchProfile();
    this.fetchPosts();
    this.fetchFriends();

    window.addEventListener('resize', () => {
      if (window.innerWidth > 720) this.mobileMenuOpen = false;
    });
  },

  methods: {
    friendInitials(name) {
      if (!name) return '?';
      return name
        .trim()
        .split(/\s+/)
        .slice(0, 2)
        .map(part => part[0])
        .join('')
        .toUpperCase();
    },

    goToProfile(userId) {
      // Укажите ваш роут: например, /profile?id=... или /profile.html?id=...
      window.location.href = `/friend?id=${userId}`;
    },

    // ===== PROFILE =====
    async fetchProfile() {
      try {
        const res = await fetch('/api/profile', { credentials: 'same-origin' });
        if (!res.ok) throw new Error('Не удалось загрузить профиль');
        this.user = await res.json();
      } catch (err) {
        console.log(err);
      }
    },

    openEditProfile() {
      this.editProfileForm.name = this.user.name || '';
      this.editProfileForm.about = this.user.about || '';
      this.editProfileForm.group = this.user.group || '';
      this.editProfileError = '';
      this.showEditProfile = true;
    },

    closeEditProfile() {
      this.showEditProfile = false;
    },

    openLightbox(imgUrl) {
      this.lightboxImg = imgUrl;
    },

    closeLightbox() {
      this.lightboxImg = null;
    },

    async saveProfile() {
      if (!this.editProfileForm.name) {
        this.editProfileError = 'Укажите имя — это поле не может быть пустым.';
        return;
      }
      this.editProfileError = '';
      this.savingProfile = true;
      try {
        const res = await fetch('/api/profile', {
          method: 'PUT',
          credentials: 'same-origin',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            name: this.editProfileForm.name,
            about: this.editProfileForm.about,
            group: this.editProfileForm.group
          })
        });
        if (!res.ok) throw new Error('Не удалось сохранить профиль');
        await this.fetchProfile();
        this.showEditProfile = false;
      } catch (err) {
        console.log(err);
        this.editProfileError = 'Не получилось сохранить изменения. Попробуйте ещё раз.';
      } finally {
        this.savingProfile = false;
      }
    },

    triggerAvatarUpload() {
      this.$refs.avatarInput.click();
    },

    async onAvatarSelected(e) {
      const file = e.target.files[0];
      if (!file) return;
      const formData = new FormData();
      formData.append('avatar', file);
      try {
        const res = await fetch('/api/profile/avatar', {
          method: 'POST',
          credentials: 'same-origin',
          body: formData
        });
        if (!res.ok) throw new Error('Не удалось обновить аватар');
        const data = await res.json();
        this.user.avatar = data.avatar;
      } catch (err) {
        console.log(err);
      } finally {
        e.target.value = '';
      }
    },

    // ===== WALL / POSTS =====
    async fetchPosts() {
      this.postsLoading = true;
      try {
        const res = await fetch('/api/post', { credentials: 'same-origin' });
        if (!res.ok) throw new Error('Не удалось загрузить стену');
        this.posts = await res.json();
      } catch (err) {
        console.log(err);
      } finally {
        this.postsLoading = false;
      }
    },

    triggerPostImage() {
      this.$refs.postImageInput.click();
    },

    onPostImageSelected(e) {
      const file = e.target.files[0];
      if (!file) return;
      this.newPost.imageFile = file;
      this.newPost.imageName = file.name;
    },

    clearPostImage() {
      this.newPost.imageFile = null;
      this.newPost.imageName = '';
      this.$refs.postImageInput.value = '';
    },

    async createPost() {
      if (!this.newPost.title || !this.newPost.text) {
        this.postError = 'Заполните тему и текст поста.';
        return;
      }
      this.postError = '';
      this.posting = true;
      try {
        const formData = new FormData();
        formData.append('title', this.newPost.title);
        formData.append('text', this.newPost.text);
        if (this.newPost.imageFile) formData.append('img', this.newPost.imageFile);

        const res = await fetch('/api/post', {
          method: 'POST',
          credentials: 'same-origin',
          body: formData
        });
        if (!res.ok) throw new Error('Не удалось опубликовать пост');

        this.newPost = { title: '', text: '', imageFile: null, imageName: '' };
        if (this.$refs.postImageInput) this.$refs.postImageInput.value = '';
        await this.fetchPosts();
      } catch (err) {
        console.log(err);
        this.postError = 'Не получилось опубликовать пост. Попробуйте ещё раз.';
      } finally {
        this.posting = false;
      }
    },

    startEditPost(post) {
      this.editingPostId = post.id;
      this.editPostForm.title = post.title;
      this.editPostForm.text = post.text;
      this.editPostError = '';
    },

    cancelEditPost() {
      this.editingPostId = null;
    },

    async saveEditPost(id) {
      if (!this.editPostForm.title || !this.editPostForm.text) {
        this.editPostError = 'Тема и текст не могут быть пустыми.';
        return;
      }
      try {
        const res = await fetch(`/api/post?id=${encodeURIComponent(id)}`, {
          method: 'PUT',
          credentials: 'same-origin',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            title: this.editPostForm.title,
            text: this.editPostForm.text
          })
        });
        if (!res.ok) throw new Error('Не удалось сохранить пост');
        this.editingPostId = null;
        await this.fetchPosts();
      } catch (err) {
        console.log(err);
        this.editPostError = 'Не получилось сохранить изменения.';
      }
    },

    async deletePost(id) {
      if (!confirm('Удалить этот пост?')) return;

      try {
        const res = await fetch(`/api/post?id=${encodeURIComponent(id)}`, {
          method: 'DELETE',
          credentials: 'same-origin'
        });
        if (!res.ok) throw new Error('Не удалось удалить пост');
        this.posts = this.posts.filter(p => p.id !== id);
      } catch (err) {
        console.log(err);
      }
    },

    // ===== FRIENDS =====
    async fetchFriends() {
      this.friendsLoading = true;
      try {
        const res = await fetch('/api/friend', { credentials: 'same-origin' });
        if (!res.ok) throw new Error('Не удалось загрузить друзей');
        this.friends = await res.json();
      } catch (err) {
        console.log(err);
      } finally {
        this.friendsLoading = false;
      }
    },

    async removeFriend(id) {
      try {
        const res = await fetch(`/api/friend?id=${encodeURIComponent(id)}`, {
          method: 'DELETE',
          credentials: 'same-origin'
        });
        if (!res.ok) throw new Error('Не удалось удалить из друзей');
        this.friends = this.friends.filter(f => f.id !== id);
      } catch (err) {
        console.log(err);
      }
    },

    // ===== MOBILE MENU =====
    toggleMobileMenu() {
      this.mobileMenuOpen = !this.mobileMenuOpen;
    },

    closeMobileMenu() {
      this.mobileMenuOpen = false;
    },

    // ===== SESSION =====
    async logout() {
      try {
        const res = await fetch('/exit', { credentials: 'same-origin' });
        if (res.ok) {
          window.location.href = '/login';
        }
      } catch (err) {
        console.log(err);
      }
    }
  }
}).mount('#app');
