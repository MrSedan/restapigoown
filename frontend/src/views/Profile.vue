<template>
    <div class="flex-container">
        <LeftMenu></LeftMenu>
        <TestProfile></TestProfile>
    </div>
</template>

<script>
import LeftMenu from "@/components/LeftMenu.vue"
import TestProfile from "@/components/TestProfile.vue"
export default {
    name: 'Profie',
    components: {
        LeftMenu,
        TestProfile
    }
}
</script>