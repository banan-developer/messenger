const App5 = {
    data() {
        return {
            test: 1,
            friends: []
        }
    },
    async mounted() {
        this.GetFrienedRequest()
    },
    methods: {
        async GetFrienedRequest(){
            try{
                const res =  await fetch("/api/incomingrequest",{
                    method: "GET",
                    credentials: "same-origin",
                })
                if (!res.ok) throw new Error("ошибка получения друзей")
                const data = await res.json()
                this.friends = data
            }catch(err){
                console.log(err)
            }
        },
        async SetAcceptStatus(friendID){
            try{
                const res = await fetch(`/api/friend?friendID=${friendID}`,{
                    method: "PUT",
                    credentials: "same-origin",
                })
                if (!res.ok) throw new Error("Принятия в друзья")
                await this.GetFriendRequest()
            }catch(err){
                console.log(err)
            }
        },
        RedirectOnProfile(){
            window.location.href = "/profile"
        }
    }
};

Vue.createApp(App5).mount('#app5');