<template>
    <div class="flex-container">
        <LeftMenu></LeftMenu>
        <div class="container">
            <h1>Chat</h1>
            <div class="message" v-for="i in messages" :key="i">
                {{i.id}}:{{i.msg}}   
            </div>
            <input type="text" v-model="message" placeholder="Message">
            <input type="submit" value="Send" @click="sendMsg()">
        </div>
    </div>
</template>

<script>
import LeftMenu from "@/components/LeftMenu.vue"
export default {
    name: "Messages",
    components: {
        LeftMenu
    },
    data() {
        return {
            id: this.$route.params.id,
            socket: null,
            messages: [],
            message: ""
        }
    },
    methods:{
        sendMsg(){
            if (this.message.length > 0) {
                let id = JSON.parse(localStorage.getItem('account')).id
                let msg =  JSON.stringify({id: id, msg: this.message})
                this.socket.send(msg)
            }
        }
    },
    mounted(){
        try{
            let id = JSON.parse(localStorage.getItem('account')).id
            this.socket = new WebSocket(`ws://localhost:8080/api/chat/${id}.${this.id}`)
            this.socket.onopen = () => {
                console.log("Socket connected")
                let msg = JSON.stringify({id: id, msg: "connected"})
                this.socket.send(msg)
            }

            this.socket.onclose = (event) => {
                console.log("Socket closed", event)
            }

            this.socket.onmessage = (msg) => {
                this.messages.push(JSON.parse(msg.data))
                console.log(msg)
            }

            this.socket.onerror = (event) => {
                console.log("Socket error: ", event)
            }
        } catch(e) {
            this.$router.push("/")
        }
        
    }
}
</script>