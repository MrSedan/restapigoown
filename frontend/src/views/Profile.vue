<template>
    <div class="flex-container">
        <LeftMenu></LeftMenu>
        <div class="content">
                <div class="profile">
                    <div class="profile-flex">
                        <img src="../static/photo.jpg" alt="">
                        <div class="profile-info">
                            <h3 class="name">{{name}}</h3>
                            <h6 class="online">Дофига онлайн</h6>
                            <p v-if="about">{{about}}</p>
                        </div>
                    </div>
                </div>
                <div class="posts">
                </div>
        </div>
    </div>
</template>

<script>
import LeftMenu from "@/components/LeftMenu.vue"
export default {
    name: 'Profie',
    data(){
        return {
            name: "",
            about: ""
        }
    },
    components: {
        LeftMenu
    },
    mounted() {
        if (!localStorage.getItem('account')){
            this.$router.push('/login')
            return
        }
        let email = JSON.parse(localStorage.getItem('account')).email
        this.$http.get(`api/user/${email}/profile`)
        .then(r => {
            this.name = r.data.first_name+' '+r.data.last_name
            this.about = r.data.about
        })
        .catch(e => {
            alert(e)
        })
    },
    created() {
       if (!localStorage.getItem('account')){
            this.$router.push('/login')
            return
        } 
    }
}
</script>