<template>
    <div class="flex-container">
        <LeftMenu></LeftMenu>
        <div class="content">
                <div class="profile">
                    <div class="profile-flex">
                        <img :src="'https://www.gravatar.com/avatar/' + id" alt="">
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
            about: "",
            id: 0
        }
    },
    methods: {
        getHash(s){
            var hash = 0, i, chr;
            if (s.length === 0) return hash;
            for (i = 0; i < s.length; i++) {
                chr   = s.charCodeAt(i);
                hash  = ((hash << 5) - hash) + chr;
                hash |= 0; // Convert to 32bit integer
            }
            return hash;
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
        try{
            let id = JSON.parse(localStorage.getItem('account')).id
            this.id = id
            this.$http.get(`api/user/${id}/profile`)
            .then(r => {
                this.name = r.data.first_name+' '+r.data.last_name
                this.about = r.data.about
            })
        }
        catch(e){
            localStorage.removeItem('account')
            this.$router.push('/login')
        }
    }
}
</script>