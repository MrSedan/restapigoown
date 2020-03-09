<template>
    <div class="flex-container">
        <div class="content">
                <div class="profile">
                    <div class="profile-flex">
                        <img :src="'/api/user/'+id+'/avatar'" alt="Error photo">
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
export default {
    name: 'Profie',
    data(){
        return {
            name: "",
            about: "",
            id: 0,
        }
    },
    mounted() {
        let id = this.$route.params.id
        this.$http.get(`/api/user/${id}/profile`)
        .then(r => {
            this.name = r.data.first_name+' '+r.data.last_name
            this.about = r.data.about
        }).catch(()=>{
            this.$router.push("/404")
        })
        this.id = id
    }
}
</script>