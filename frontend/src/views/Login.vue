<template>
    <div class="login-card">
        <h2>Login</h2>
        <input type="email" placeholder="Email" v-model="email">
        <input type="password" placeholder="Password" v-model="password">
        <div id="pass_error" v-if="error">{{ error }}</div>
        <input type="button" value="Sign In"  @click="login()">
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
            error: ""
        }
    },
    methods: {
        login(){
            if (!this.email || !this.password){
                this.error = "All fields are required!"
                return
            }
            let credentials = JSON.stringify({
                email: this.email,
                password: this.password
            })
            this.$http.post('api/user/login', credentials).then(user =>{
                localStorage.setItem("account", JSON.stringify(user.data))
                this.$router.push('/profile')
            }).catch(() => {this.error = "Invalid email or password!"})
        }
    }
}
</script>