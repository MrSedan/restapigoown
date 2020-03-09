<template>
    <div class="flex-container">
        <div class="content">
            <h1>Chat with {{ this.$route.params.id }}</h1>
            <div id="chat">
                <div id="msg-list">
                <div class="message" v-for="i in 30" :key="i">
                    <p>{{i}} Test</p>
                </div>  
                </div>
                <div id="send_msg">
                        <input type="text" placeholder="Message">
                        <input type="button" value="Send">
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: "Chat",
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


<style lang="scss" src="@/static/chat.scss" scoped></style>