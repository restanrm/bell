<template>
  <div class="list">
    <div class="jumbotron">
      <div class="container">
        <div class="row">
          <div class="col" v-for="sound in sounds"><button class="btn btn-primary" v-on:click="play(sound.name)">{{sound.name}}</button></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'list',
    data () {
      return {
        sounds: [],
        playURL: '/api/v1/play/',
        soundsURL: '/api/v1/sounds'
      }
    },
    methods: {
      play: function (sound) {
        this.$http.get(this.playURL + sound)
      },
      updateSounds: function () {
        this.$http.get(this.soundsURL).then(response => {
          this.sounds = response.data
        })
      }
    },
    created: function () {
      this.updateSounds()
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
  h1,
  h2 {
    font-weight: normal;
  }
  
  ul {
    list-style-type: none;
    padding: 0;
  }
  
  li {
    display: inline-block;
    margin: 0 10px;
  }
  
  a {
    color: #35495E;
  }
</style>
