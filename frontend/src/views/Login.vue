<template>
    <div class="login-card">
        <h2>Login</h2>
        <input type="email" placeholder="Email" v-model="email">
        <input type="password" placeholder="Password" v-model="password">
        <input type="button" value="Sign In"  @click="login()">
        {{info}}
        <router-link to="/signup" class="signup">Sign Up?</router-link>
    </div>
</template>

<style src="../static/login.scss" lang="scss" scoped></style>

<script>
export default {
    data(){
        return {
            email: "",
            password: "",
            info: ""
        }
    },
    methods: {
        login(){
            if (!this.email || !this.password){
                alert('err')
                return
            }
            let credentials = JSON.stringify({
                email: this.email,
                password: this.password
            })
            this.$http.post('api/user/login', credentials).then(user =>{
                localStorage.setItem("account", JSON.stringify(user.data))
                this.$router.push('/profile')
            }).catch(e => {alert(e)})
        }
    }
}
</script>