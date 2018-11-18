<template lang="jade">
  .register

    .list

      button.btn.btn-primary(
        v-on:click="register()"
      ) {{ share ? "disconnect" : "connect" }}

      input(
        type="text"
        placeHolder="Your name"
        value=""
        v-model="name"
      )


</template>

<script>
  import BootstrapToggle from 'vue-bootstrap-toggle'
  export default {
    name: 'register',
    components: {BootstrapToggle},
    data () {
      var loc = window.location
      var uri = ''
      if (loc.protocol === 'http:') {
        uri = 'ws:'
      } else {
        uri = 'wss:'
      }
      uri += '//' + loc.host
      if (process.env.NODE_ENV === 'development') {
        uri = 'ws://localhost:10101'
      }
      uri += '/api/v1/clients/register'
      return {
        ws: {},
        name: '',
        share: false,
        registerURL: uri
      }
    },
    props: {
      soundToPlay: String
    },
    methods: {
      updateDest: function (destination) {
        this.$emit('update:destination', destination)
      },
      register: function () {
        var vm = this
        // toggle the connect state
        this.share = !this.share
        if (!this.share) {
          this.ws.close()
          return
        }

        // connect the clients to the websocket
        this.ws = new WebSocket(vm.registerURL)
        this.ws.onopen = function (event) {
          vm.ws.send(JSON.stringify({name: vm.name}))
        }
        this.ws.onerror = function (event) {
          console.log('An error happened: ' + event)
        }
        this.ws.onclose = function (event) {
          console.log('close function called' + event)
        }
        this.ws.onmessage = function (event) {
          var msg = JSON.parse(event.data)
          if (msg.name && msg.name !== '') {
            vm.name = msg.name
            return
          }
          if (msg.type && msg.type === 'sound') {
            vm.$emit('register:play', msg.data)
          } else {
            console.log('Received unhandled message: ' + msg)
          }
        }
      }
    },
    created: function () {
      this.register()
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss">
  @import "../styles/settings.scss";

  .register {
    padding: 2rem 1rem;
    background: $primary-darker;

    .list {
      display: grid;
      grid-gap: 20px;
      grid-template-columns: 120px minmax(400px, 1fr);
    }

  }

</style>
