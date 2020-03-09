<template>
    <div class="inner-width">
        <router-link to="/" tag="p">Mysite</router-link>
        <div class="nav-hid" :class="{active: isAuth}">
            <i class="menu-toggle-btn fas" @click="changeActive()" :class="[!isActive ? 'fa-bars' : 'fa-times']"></i>
            <nav class="navigation-menu" :class="{active: isActive}">
                <router-link to="/"><i class="fas fa-home"></i>Home</router-link>
                <router-link to="/messages"><i class="far fa-envelope"></i>Messages</router-link>
            </nav>
        </div>
    </div>
</template>

<script>
export default {
    data(){
        return{
            isActive: false,
            isAuth: false
        }
    },
    methods: {
        changeActive(){
            this.isActive = !this.isActive
        }
    },
    watch:{
        '$route'(){
                try{
                    this.isActive = false
                    let id = JSON.parse(localStorage.getItem('account')).id
                    if (id > 0 && !this.isAuth) {
                        this.isAuth = true
                    }
                } catch(e){
                    this.isAuth = false
                    localStorage.removeItem('account')
                }
        }
    }
}
</script>

<style lang="scss" src="@/static/topmenu.scss" scoped></style>