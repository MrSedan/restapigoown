<template>
    <div class="login-card">
        <h2>Login</h2>
        <input type="email" placeholder="Email">
        <input type="password" placeholder="Password">
        <input type="button" value="Sign In">
        <router-link to="/signup" class="signup">Sign Up?</router-link>
    </div>
</template>

<style src="../static/login.scss" lang="scss" scoped>