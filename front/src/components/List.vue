<template lang="jade">
  .list
    .jumbotron.player
      .container-fluid
        audio(
          :src="soundPath",
          id="player",
          autoplay,
        )
        .row.searcher
          .toggler.col-sm-12.col-md-2.col-lg-2
            bootstrap-toggle(
              v-model="playOnServer",
              :options="{on: '<i class=\"fa fa-play\"></i> Server', off:'<i class=\"fa fa-play\"></i> Locally'}",
            )
          input.col-sm-12.offset-sm-0.col-md-6.offset-md-4.col-lg-4.offset-lg-6(type="text",v-model="search",placeholder="search sound")
        .row
          .col-sm-6.col-md-4.col-lg-2(v-for="sound in filteredSounds")
            button.btn.btn-primary.play-btn(
              v-on:click="play(sound.name)",
              ){{sound.name}}
</template>

<script>
  import BootstrapToggle from 'vue-bootstrap-toggle'
  export default {
    name: 'list',
    components: {BootstrapToggle},
    data () {
      var basepath = ''
      if (process.env.NODE_ENV === 'development') {
        basepath = 'http://localhost:10101'
      };
      return {
        search: '',
        sounds: [],
        playOnServer: true,
        soundPath: '',
        playLocallyURL: basepath + '/api/v1/sounds/',
        playURL: basepath + '/api/v1/play/',
        soundsURL: basepath + '/api/v1/sounds'
      }
    },
    methods: {
      play: function (sound) {
        console.log(this.playOnServer)
        if (this.playOnServer) {
          this.$http.get(this.playURL + sound)
        } else {
          this.soundPath = this.playLocallyURL + sound
        }
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

  .searcher {
    margin: 15px 0;
    .toggler {
      padding: 0;

      .toggle-handle {
        background: white;
        border-radius: 0px;
      }
      .toggle-on {
        background: $primary;
      }
      .toggle-off {
        background: $secondary-light;
        color: white;
      }
    }
  }

  .btn-primary {
    border-color: $primary-light;
    &:hover{
      background-color: $primary;
      border-color: $primary-light;
    }
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
      border-color: $primary-light;
    }
  }

</style>
