<template>
    <div class="flex-container">
        <vue-title title="Chat"></vue-title>
        <div class="content">
            <h1>Chat with <router-link :to="'/profile/'+id" id="name">@{{ getName(id) }}</router-link></h1>
            <div id="chat">
                <div id="msg-list">
                <div class="message" v-for="i in messages" :key="i.id">
                    <router-link :to="'/profile/'+i.from"><img :src="'/api/user/'+i.from+'/avatar'" alt="Error photo"></router-link>
                    <div class="msg-body">
                        <router-link :to="'/profile/'+i.from" class="name">@{{getName(i.from)}}</router-link>
                        <p class="msg_text">{{i.body}}</p>
                        <p class="date">{{i.timestamp | formatDate}}</p>
                    </div>
                </div>  
                </div>
                <div id="send_msg" @keyup.enter="sendMsg()">
                        <input type="text" placeholder="Message" v-model="body">
                        <button id="send-btn" @click="sendMsg()"><i class="fas fa-paper-plane"></i></button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: "Chat",
    data(){
        return {
            id:Number(this.$route.params.id),
            myid: null,
            messages:[],
            socket: null,
            body: "",
            name: "",
            myname: "",
            token:""
        }
    },
    methods: {
        sendMsg(){
            if (this.body.trim().length == 0){
                this.body = ""
                return
            }
            this.socket.send(JSON.stringify({
                token: this.token,
                from: this.myid,
                to: this.id,
                body: this.body,
            }))
            this.body = ""
        },
        getName(id){
            if (id==this.id){
                return this.name
            }
            else{
                return this.myname
            }
        },
        socketCon(){
            try{
                this.socket = new WebSocket(`wss://hackcergroup.tk/api/chat/${this.myid}.${this.id}?token=${this.token}&id=${this.myid}`)
                this.socket.onopen = () => {
                    console.log("Socket connected")
                }

                this.socket.onclose = (event) => {
                    if (event.code == 1006){
                        this.socketCon()
                        return
                    }
                    console.log("Socket closed", event)
                }

                this.socket.onmessage = (msg) => {
                    let mes = JSON.parse(msg.data)
                    mes.timestamp = mes.time
                    this.messages.push(mes)
                    console.log(mes)
                }

                this.socket.onerror = (event) => {
                    console.log("Socket error: ", event)
                }
            } catch(e) {
                this.$router.push("/")
            }
        }
    },
    beforeDestroy(){
        this.socket.close()
    },
    mounted(){
        if(localStorage.getItem('account')){
            var u = JSON.parse(localStorage.getItem('account'))
            if (u.id == 0 || u.token == 0){
                localStorage.removeItem('account')
                this.$router.push("/login")
                return
            }
            this.token = u.token
            this.$http.post("/api/checkauth", this.$qs.stringify({id: u.id, token: u.token}), {
                'Content-Type': 'application/x-www-form-urlencoded'
            }).catch(()=>{
                localStorage.removeItem('account')
                this.$router.push("/login")
                return
            })
            this.myid = u.id
        } else {
            this.$router.push("/login")
            return
        }
        this.$http.get(`/api/user/${this.myid}/profile`)
        .then(r => {
            this.myname = r.data.user_name
        })
        this.$http.get(`/api/user/${this.id}/profile`)
        .then(r => {
            this.name = r.data.user_name
        })
        .catch(()=>{
            this.$router.push("/404")
            return
        })
        this.socketCon()
        this.$http.post(`/api/chat/${this.myid}.${this.id}/gethistory`, this.$qs.stringify({id: u.id, token: u.token}), {
                'Content-Type': 'application/x-www-form-urlencoded'
            }).then(r => {
                let mess = []
                for (let i=0;i<r.data.length; i++){
                    var mes = {}
                    var date = (new Date(r.data[i].time*1000)).toLocaleString()
                    mes = r.data[i]
                    mes.timestamp = r.data[i].time
                    mes.time = date
                    mess.push(mes)
                }
                mess.sort((a,b)=>{
                    return a.timestamp > b.timestamp
                })
                this.messages = mess
            })
        
    }
}
</script>


<style lang="scss" src="@/static/chat.scss" scoped></style>