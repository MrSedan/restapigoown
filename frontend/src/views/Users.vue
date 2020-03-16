<template>
    <div class="flex-container">
        <vue-title title="Users"></vue-title>
        <div class="container">
            <h1 v-if="users.length != 0">Users</h1>
            <h1 v-else>Users not found :-(</h1>
            <router-link v-for="i in users" :key="i.id" :to="'/profile/'+i.id" class="userurl" tag="div">
            <img :src="'/api/user/'+i.id+'/avatar'">
            <h3>{{i.name}}</h3>
            </router-link>
            
        </div>
    </div>
</template>

<script>
export default {
    name: "Users",
    data() {
        return{
            users: [],
            id: null
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
        this.$http.post(`/api/user/getalluser`, this.$qs.stringify({id: this.id, token: u.token}), {
                'Content-Type': 'application/x-www-form-urlencoded'
            }).then(r => {
                for (let i=0;i<r.data.length;i++){
                    if (this.id != r.data[i].id){

                        this.users.push({id: r.data[i].id, name: r.data[i].user_name})
                    }
                }
            })
    }
}
</script>

<style lang="scss" scoped>
.container{
    margin: 0 40px;
}

.userurl{
    min-width: 300px;
    height: 80px;
    margin-bottom: 20px;
    border-radius: 15px;
    box-shadow: 5px 5px 10px rgba(0,0,0,0.2);
    padding: 0 10px;
    text-decoration: none;
    color: black;
    outline: none;
    display: flex;
    transition: .3s linear;
    align-items: center;
    vertical-align: middle;
}

.userurl h3{
    display: inline-block;
    line-height: 80px;
    margin: auto 10px;
}

.userurl img{
    max-width: 74px;
    max-height: 74px;
    border-radius: 50%;
    display: inline-block;
    margin-right: 10px;
   
}

.userurl:hover{
    background: rgb(173, 173, 173);
    cursor: pointer;
}
</style>