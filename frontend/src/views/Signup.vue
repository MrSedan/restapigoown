<template>
    <div class="signup-card">
        <h2>Sign Up</h2>
        <input type="text" placeholder="First name">
        <input type="text" placeholder="Last Name">
        <input type="email" placeholder="E-mail">
        <input type="password" placeholder="Password">
        <input type="password" placeholder="Repeat password" autocomplete="off">
        <input type="button" value="Sign Up">
        <router-link to="/login" class="login">Login?</router-link>
    </div>
</template>

<style src="../static/signup.scss" lang="scss" scoped>