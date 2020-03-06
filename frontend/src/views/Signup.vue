<template>
    <div class="signup-card">
        <h2>Sign Up</h2>
        <input type="text" placeholder="First name" v-model="firstname">
        <input type="text" placeholder="Last Name" v-model="lastname">
        <input type="email" placeholder="E-mail" v-model="email">
        <input type="password" placeholder="Password" v-model="password">
        <input type="password" placeholder="Repeat password" autocomplete="off" v-model="re_password">
        <div id="pass_error" v-if="error">{{ error }}</div>
        <input type="button" value="Sign Up" @click="CreateUser()">
        <router-link to="/login" class="login">Login?</router-link>
    </div>
</template>

<style src="../static/signup.scss" lang="scss" scoped></style>
<script>
export default {
    name: "SignUp",
    data(){
        return {
            firstname: "",
            lastname: "",
            email: "",
            password: "",
            re_password: "",
            error: ""
        }
    },
    methods: {
        CreateUser(){
            if (!(this.firstname.length && this.lastname.length && this.email.length && this.password.length && this.re_password.length)) {
                this.error = "All fields are required!"
                return
            }
            if (this.password != this.re_password){
                this.error="Passwords does not match!"
                return
            }
            if (this.password.length < 8 || this.password.length > 64){
                this.error = "Password length must be between 8 and 64!"
            }
            let ac = JSON.stringify({
                first_name: this.firstname,
                last_name: this.lastname,
                email: this.email.toLowerCase().trim(),
                password: this.password
            })
            this.$http.post('/api/user/create', ac).then(() => {this.$router.push('/login')})
            .catch(e => {
                if (e.response.status == 400){
                    this.error = "This email alreay registered!"
                }
            })
        }
    },
}
</script>