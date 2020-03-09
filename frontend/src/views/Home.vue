<template>
    <div id="greetings">
        <h1>Welcome to HackerGroup!</h1>
        <h2>This site now in development!!!</h2>
        <router-link to="/login" class="butt">Login</router-link>
        <router-link to="/signup" class="butt">Sign up</router-link>
    </div>
</template>

<script>
export default {
    name: "Home",
    mounted(){
        if(localStorage.getItem('account')){
            try{
                let id = JSON.parse(localStorage.getItem('account')).id
                this.$router.push("/profile/"+id)
            } catch(e){
                localStorage.removeItem('account')
            }
        }
    }
}
</script>

<style scoped>
    #greetings {
        margin-top: 60px;
        text-align: center;
        position: absolute;
        box-shadow: 5px 5px 10px rgba(0,0,0,0.2);
        padding: 15px;
        min-width: 300px;
        left: 50%;
        transform: translateX(-50%);
    }
    #greetings .butt{
        display: block;
        border: 2px solid rgba(28, 56, 179, 0.623);
        outline: none;
        text-decoration: none;
        min-width: 40px;
        color: black;
        padding: 10px;
        width: 20%;
        margin-bottom: 10px;
        margin-left: auto;
        margin-right: auto;
        border-radius: 24px;
        transition: .3s ease;
    }
    #greetings .butt:hover{
        background: rgb(28, 56, 179);
        color: white;
    }
</style>