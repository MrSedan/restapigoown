<template>
    <div class="flex-container">
        <vue-title title="Profile"></vue-title>
        <div class="content">
                <div class="profile">
                    <div class="profile-flex">
                        <img :src="'/api/user/'+page_id+'/avatar'" alt="Error photo">
                        <div class="profile-info">
                            <div v-if="name != 'undefined undefined'">
                            <h3 class="name">{{name}}</h3>
                            <h5>@{{user_name}}</h5>
                            </div>
                            <div v-else>
                                <h3 class="name">@{{user_name}}</h3>
                            </div>
                            <h6 class="online">Дофига онлайн</h6>
                            <p v-if="about" id="about">{{about}}</p>
                        </div>
                        <div class="after-profile">
                            <router-link :to="'/chat/'+page_id" id="messag"><i class="fas fa-paper-plane"></i></router-link>
                            <router-link to="/editprofile" id="editProfile" v-if="page_id == id">Edit</router-link>
                        </div>
                    </div>
                </div>
                <div class="posts">
                </div>
        </div>
    </div>
</template>

<script>
export default {
    name: 'Profie',
    data(){
        return {
            user_name: "",
            about: "",
            id: 0,
            page_id: this.$route.params.id,
            name: ""
        }
    },
    mounted() {
        let id = this.$route.params.id
        this.$http.get(`/api/user/${id}/profile`)
        .then(r => {
            this.user_name = r.data.user_name
            this.about = r.data.about
            this.name = r.data.first_name + ' ' + r.data.last_name
        }).catch(()=>{
            let myId = JSON.parse(localStorage.getItem('account')).id
            if (myId==id){
                localStorage.removeItem('account')
            }
            this.$router.push("/404")
        })
        if(localStorage.getItem('account')){
            var u = JSON.parse(localStorage.getItem('account'))
            if (u.id == 0 || u.token == 0){
                localStorage.removeItem('account')
                this.$router.push("/login")
                return
            }
            this.$http.post("/api/checkauth", this.$qs.stringify({id: u.id, token: u.token}), {
                'Content-Type': 'application/x-www-form-urlencoded'
            }).catch(()=>{
                localStorage.removeItem('account')
                this.$router.push("/login")
                return
            })
            this.id = u.id
        } else {
            this.$router.push("/login")
            return
        }
        this.id = u.id
    }
}
</script>

<style lang="scss" src="@/static/profile.scss" scoped></style>