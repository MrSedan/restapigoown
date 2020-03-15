<template>
    <div id="edit-profile-card" @keyup.enter="sendProfile()">
        <vue-title title="Edit Profile"></vue-title>
        <h1 class="card-name">Edit your profile</h1>
        <input type="text" placeholder="First Name" class="inputText" v-model="firstname">
        <input type="text" placeholder="Last name" class="inputText" v-model="lastname">
        <textarea class="inputText" placeholder="About" v-model="about"></textarea>
        <ul id="pass_error" v-if="errors.length > 0">
            <li v-for="error in errors" :key="error">{{ error }}</li>
        </ul>
        <input type="button" value="Change" class="btn" @click="sendProfile()">
    </div>
</template>


<script>
export default {
    name: "EditProfile",
    data(){
        return{
            real_firstname: "",
            real_lastname: "",
            firstname: "",
            lastname: "",
            about: "",
            real_about: "",
            id: null,
            errors: []
        }
    },
    methods:{
        sendProfile(){
            this.errors = []
            if (!this.firstname.match(/^[a-zа-яё]{0,20}$/i)){
                this.errors.push("First Name must contain only letters")
            }
            if (!this.lastname.match(/^[a-zа-яё]{0,30}$/i)){
                this.errors.push("Last Name must contain only letters")
            }
            if (!this.about.match(/^[a-z\d\s.:;\-!'@#%$,+=^()*а-яё]{0,2000}$/i)){
                this.errors.push("Your about contain bad symbols or very large(>120 symbols)")
            }
            if (this.errors.length > 0){
                this.obnull()
                return
            }
            let token = JSON.parse(localStorage.getItem('account')).token
            if (token.length == 0 || token == null){
                localStorage.removeItem('account')
                this.$router.push("/login")
                return
            }
            this.$http.post(`/api/user/${this.id}/edit/profile`, {
                token: token,
                first_name: this.firstname.capitalize(true),
                last_name: this.lastname.capitalize(true),
                about: this.about
            }).then(()=>{
                this.$router.push(`/profile/${this.id}`)
            }).catch(() => {
                this.errors.push("An error occurred!")
            })
        },
        obnull(){
            this.firstname = this.real_firstname
            this.lastname = this.real_lastname
            this.about = this.real_about
        }
    },
    mounted(){
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
        this.$http.get(`/api/user/${this.id}/profile`)
        .then(r => {
            if (r.data.about != undefined){
                this.real_about = r.data.about
            }
            if (r.data.first_name != undefined){
                this.real_firstname = r.data.first_name
            }
            if (r.data.last_name != undefined){
                this.real_lastname = r.data.last_name
            }
            this.obnull()
        })
        
    }
}
</script>

<style lang="scss" scoped>
#pass_error {
    background: rgba($color: #b3022065, $alpha: 0.3);
    font-size: 13px;
    margin: 10px;
    padding: 8px;
    width: 350px;
    margin-left: auto;
    margin-right: auto;
    text-align: left;
    border-radius: 8px;
    color: rgba($color: #5c5c5c, $alpha: 1.0);
    list-style: square;
}

#pass_error li{
    margin-left: 20px;
}
#edit-profile-card{
    text-align: center;
    position: absolute;
    min-width: 300px;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    box-shadow: 5px 5px 10px rgba(0,0,0,0.2);
    border-radius: 8px;
    padding: 30px;
}

.inputText{
    display: block;
    margin-left: auto;
    margin-right: auto;
    margin-bottom: 10px;
}
.card-name{
    margin: 10px auto 20px;
}
input{
    background: none;
    outline: none;
    border: 2px solid rgba(28, 56, 179, 0.623);
    text-align: center;
    padding: 8px 14px;
    border-radius: 24px;
    width: 200px;
    transition: .3s linear;
}
input:focus{
    border: 2px solid rgba(179, 28, 28, 0.623);
    width: 250px;
}
textarea{
    padding: 10px;
    border-radius: 8px;
    resize: none;
    outline: none;
    background: none;
    width: 200px;
    height: 1.2em;
    border: 2px solid rgba(28, 56, 179, 0.623);
    transition: .3s linear;
}
textarea:focus{
    border: 2px solid rgba(179, 28, 28, 0.623);
    height: 4em;
    width: 220px;
}
#edit-profile-card .btn{ 
    min-width: 80px;
    width: 20%;
    border: 2px solid rgba(28, 56, 179, 0.623);
    outline: none;
    color: black;
    padding: 10px;
    width: 20%;
    margin-bottom: 10px;
    margin-left: auto;
    margin-right: auto;
    border-radius: 24px;
    transition: .3s ease;
}
#edit-profile-card .btn:hover{
    background: rgb(28, 56, 179);
    color: white;
}
</style>