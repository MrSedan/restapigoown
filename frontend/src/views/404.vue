<template>
    <div id="NotFound">
        <vue-title title="404"></vue-title>
        <h4>Requested page not found!</h4>
        <h1>404</h1>
        <router-link to="/" id="home">Go home</router-link>
    </div>
</template>

<style scoped>
#NotFound{
    text-align: center;
    margin-top: 60px;
    left: 50%;
    top:50%;
}
#NotFound h1{
    font-size: 60pt;
    font-weight: bolder;
}
#NotFound #home{
    padding: 10px;
    border:2px solid rgba(28, 56, 179, 0.623);
    border-radius: 24px;
    outline: none;
    text-decoration: none;
    color: white;
    background: rgb(28, 56, 179);
    transition: .3s ease;
}
#NotFound #home:hover{
    color: black;
    background: none;
}
</style>