<template>
    <div class="signup-card" @keyup.enter="CreateUser()">
        <vue-title title="Sign up"></vue-title>
        <h2>Sign Up</h2>
        <input type="text" placeholder="Nickname" v-model="nickname">
        <input type="text" placeholder="E-mail" v-model="email">
        <input type="password" placeholder="Password" v-model="password">
        <input type="password" placeholder="Repeat password" autocomplete="off" v-model="re_password">
        <ul id="pass_error" v-if="errors.length > 0">
            <li v-for="error in errors" :key="error">{{ error }}</li>
        </ul>
        <input type="button" value="Sign Up" @click="CreateUser()">
        <router-link to="/login" class="login">Login?</router-link>
    </div>
</template>

<style src="../static/signup.scss" lang="scss" scoped></style>
<script>
String.prototype.capitalize = function(lower) {
    return (lower ? this.toLowerCase() : this).replace(/(?:^|\s)\S/g, function(a) { return a.toUpperCase(); });
};
export default {
    name: "SignUp",
    data(){
        return {
            nickname: "",
            email: "",
            password: "",
            re_password: "",
            errors: []
        }
    },
    methods: {
        CreateUser(){
            this.errors = []
            if(!this.nickname.match(/^(?!\d)(?=.*[a-zA-Z\d])(?=\S+$).{2,15}$/)){
                this.errors.push("Nickname can contain only letters and numbers")
            }
            if(!this.email.match(/^[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z]{2,4}$/)){
                this.errors.push("This is not an email!")
            }
            if(!this.password.match(/^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9]).{8,64}$/)){
                this.errors.push("Password must contain lowercase letters, numbers and capital letters and its length must be between 8 and 64")
            }
            if (!(this.email.length && this.password.length && this.re_password.length)) {
                this.errors.push("All fields are required!")
            }
            if (this.password != this.re_password){
                this.error="Passwords does not match!"
            }
            if (this.errors.length > 0){
                return
            }
            let ac = JSON.stringify({
                user_name: this.nickname.toLowerCase().trim(),
                email: this.email.toLowerCase().trim(),
                password: this.password
            })
            this.$http.post('/api/user/create', ac).then(() => {this.$router.push('/login')})
            .catch(e => {
                if (e.response.status == 400){
                    this.errors.push("This email is already registered!")
                }
            })
        }
    },
}
</script>