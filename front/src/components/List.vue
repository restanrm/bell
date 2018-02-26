<template lang="jade">
  .list
    .jumbotron.player
      .container-fluid
        .row
          .col-4
          input.col-4(type="text",v-model="search",placeholder="search sound")
        .row
          .col-sm-6.col-md-4.col-lg-2(v-for="sound in filteredSounds")
            button.btn.btn-primary.play-btn(v-on:click="play(sound.name)"){{sound.name}}
</template>

<script>
  export default {
    name: 'list',
    data () {
      var basepath = ''
      if (process.env.NODE_ENV === 'development') {
        basepath = 'http://localhost:10101'
      };
      return {
        search: '',
        sounds: [],
        playURL: basepath + '/api/v1/play/',
        soundsURL: basepath + '/api/v1/sounds'
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
    computed: {
      filteredSounds: function () {
        return this.sounds.filter(sound => {
          return sound.name.toLowerCase().includes(this.search.toLowerCase())
        })
      }
    },
    created: function () {
      this.updateSounds()
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss">
  @import "../styles/settings.scss";

  .player {
    background: $primary-darker;
    margin: 0;
  }
  
  ul {
    list-style-type: none;
    padding: 0;
    li {
      display: inline-block;
      margin: 0 10px;
    }
  }
  
  button.play-btn {
    margin: 15px 0 !important;
    padding: 15px 0 !important;
    width: 100%;
		background: $primary-dark;
    border-color: $primary-light;
    &:hover{
      background: $secondary-light;
    }
  }

</style>
