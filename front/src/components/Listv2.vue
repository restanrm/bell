<template lang="jade">
  div.player
    audio(
      :src="soundPath",
      ref="player",
      autoplay,
    )

    // list options (enable tags or select sound)
    .options
      .searcher
        .toggler
          bootstrap-toggle(
            v-model="tag",
            :options="{off:'<i class=\"fa fa-play\"></i> Sounds', on:'<i class=\"fa fa-pause\"></i> Tags'}"
          )
      input(type="text",v-model="search",placeholder="search sound")

    .list
      div(v-for="elem in elements")
        button.btn.btn-primary.play-btn(
          v-on:click="play(elem)",
          ){{elem}}
</template>

<script>
  import BootstrapToggle from 'vue-bootstrap-toggle'
  export default {
    name: 'list',
    components: {BootstrapToggle},
    props: {
      destination: String,
      // tag: Boolean,
      soundToPlay: String
    },
    data () {
      var basepath = ''
      if (process.env.NODE_ENV === 'development') {
        basepath = 'http://localhost:10101'
      };
      return {
        tag: false,
        search: '',
        sounds: [],
        playOnServer: true,
        soundPath: '',
        playLocallyURL: basepath + '/api/v1/sounds/',
        playURL: basepath + '/api/v1/play/',
        soundsURL: basepath + '/api/v1/sounds'
      }
    },
    computed: {
      tagNames: function () {
        var t = []
        // push all tags in 't'
        this.sounds.forEach(sound => {
          if (sound.tags) {
            sound.tags.forEach(tag => { t.push(tag) })
          }
        })
        // deduplicates tags in t
        var res = t.filter(function (value, index, self) {
          return self.indexOf(value) === index
        }).filter(tag => {
          // filter out with filter
          return tag.toLowerCase().includes(this.search.toLowerCase())
        })
        return res
      },
      soundNames: function () {
        var list = []
        // filter sounds and then put all sound.name in a list.
        this.sounds.filter(sound => {
          return sound.name.toLowerCase().includes(this.search.toLowerCase())
        }).forEach(sound => {
          list.push(sound.name)
        })
        return list
      },
      elements: function () {
        if (this.tag) {
          return this.tagNames
        } else {
          return this.soundNames
        }
      }
    },
    methods: {
      play: function (sound) {
        var url = this.playURL + sound
        if (this.destination !== '') {
          url += '?destination=' + this.destination
        }
        this.$http.get(url)
      },
      updateSounds: function () {
        this.$http.get(this.soundsURL).then(response => {
          this.sounds = response.data
        })
      }
    },
    watch: {
      soundToPlay: function () {
        console.log('Watcher function: ' + this.soundToPlay)
        // add event listener on ended event of player to reset soundToPlay value.
        // else the watch trigger is never redone on update of same button
        var self = this
        this.$refs.player.addEventListener('ended', function () {
          self.$emit('update:soundToPlay', '')
        })

        this.soundPath = this.playLocallyURL + this.soundToPlay
        this.$refs.player.play()
      }
    },
    created: function () {
      this.player = document.getElementById('player')
      this.updateSounds()
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss">
  @import "../styles/settings.scss";

  .player {
    background: $primary-darker;
    padding: 2rem 1rem;
    margin: 0;
  }

  .searcher {
    margin: 15px 0;
    display: flex;
    align-items: center;
    .toggler {
      height: 34px;
      padding: 0;
      display: block;
      margin: 15px 0;

      .btn {
        max-height: 34px;
        border: none;
      }

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

  .options {
    display: grid;
    // grid-template-columns: 200px auto;
    grid-template-columns: repeat(2, minmax(100px, 1fr));
  }

  div div.list {
    display: grid;
    grid-gap: 20px;
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));

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
  }
</style>
