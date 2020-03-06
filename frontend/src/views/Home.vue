<template>
    <div id="greetings">
        <h1>Welcome to HackerGroup!</h1>
        <h2>This site now in development!!!</h2>
        <router-link to="/login" class="butt">Login</router-link>
        <router-link to="/signup" class="butt">Sign up</router-link>
    </div>
</template>

<style scoped>
    #greetings {
        margin-top: 60px;
        text-align: center;
        box-shadow: 5px 5px 10px rgba(0,0,0,0.2);
        width: 40%;
        transform: translateX(70%);
        padding: 15px;
    }
    #greetings .butt{
        display: block;
        margin: 10px;
        border: 2px solid rgba(28, 56, 179, 0.623);
        outline: none;
        text-decoration: none;
        color: black;
        padding: 10px;
        width: 10%;
        transform: translateX(300%);
        border-radius: 24px;
        transition: .3s ease;
    }
    #greetings .butt:hover{
        background: rgb(28, 56, 179);
        color: white;
    }
</style>