<template>
  <div id="app">
    <TopMenu></TopMenu>
    <router-view/>
  </div>
</template>

<style src="./static/main.scss" lang="scss"></style>

<script>
import TopMenu from "@/components/TopMenu.vue"
export default {
  components: {
    TopMenu
  }
}
</script>