<template lang="jade">
  .destination
    button(
      v-on:click="updateClients()"
    ) Refresh
    div.list
      div(v-for="c in clients")
        input(
          type="radio"
          v-model="client"
          v-bind:value="c"
          v-bind:id="c"
        )
        label(v-bind:for="c") {{c}}


</template>

<script>
  export default {
    name: 'destination',
    data () {
      var basepath = ''
      if (process.env.NODE_ENV === 'development') {
        basepath = 'http://localhost:10101'
      };
      return {
        clients: [],
        client: '',
        clientsURL: basepath + '/api/v1/clients'
      }
    },
    props: ['destination'],
    watch: {
      client: function () {
        this.updateDest(this.client)
      }
    },
    methods: {
      updateDest: function (destination) {
        this.$emit('update:destination', destination)
      },
      updateClients: function () {
        this.$http.get(this.clientsURL).then(response => {
          this.clients = response.data.clients
        })
      }
    },
    created: function () {
      this.updateClients()
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss">
  @import "../styles/settings.scss";

  .destination {
    padding: 4rem 2rem;
    background: $primary-darker;
  }

  div {
    color: white;
    align-items: center;

    div.list {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
    }
  }

</style>
