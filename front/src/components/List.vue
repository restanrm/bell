<template>
  <div class="list">
    <div class="jumbotron player">
      <div class="container-fluid">
        <div class="row">
          <div class="col-sm-6 col-md-4 col-lg-2" v-for="sound in sounds"><button class="btn btn-primary play-btn" v-on:click="play(sound.name)">{{sound.name}}</button></div>
        </div>
      </div>
    </div>
  </div>
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
