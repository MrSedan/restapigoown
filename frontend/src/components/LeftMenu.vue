<template>
    <div class="left-menu">
        <ul>
            <li><router-link :to="'/profile/'+id">Home</router-link></li>
            <li><router-link to="/login">Login</router-link></li>
            <li><router-link to="/signup">Sign Up</router-link></li>
        </ul>
    </div>
</template>

<script>
export default {
    props: {
        id:String
    },
    mounted(){
        if (!localStorage.getItem('account')){
            this.$router.push('/login')
        } else {
            try{
                let id = JSON.parse(localStorage.getItem('account')).id
                this.$http.get(`/api/user/${id}/profile`).catch(()=>{
                    localStorage.removeItem('account')
                    this.$router.push('/login')
                })
            } catch(e){
                localStorage.removeItem('account')
                this.$router.push('/login')
            }
        }
        
    }
}
</script>