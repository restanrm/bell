
<template lang="jade">
  .app

    nav.navbar.navbar-expand-lg

      a.navbar-brand(
        href="#"
        v-on:click="currentTab = ''"
      ) Bell

      div.collapse.navbar-collapse#navbarSupportedContent
        ul.navbar-nav.mr-auto
          li.nav-item(
            v-for="tab in tabs"
            v-on:click="currentTab = tab"
            v-bind:key="tab"
          )
            a(href="#") {{tab}}

      button.navbar-toggler.custom-toggler(type="button",
        data-toggle="collapse",
        data-target="#navbarSupportedContent",
        aria-controls="navbarSupportedContent",
        aria-expanded="true",
        aria-label="Toggle navigation",
        )
        span.custom-toggler
          i.fa.fa-bars.white

    Destination(
      v-bind:class="{hidden: selectedComponent(currentTab, 'Destination')}"
      v-bind:destination="dest"
      v-on:update:destination="dest = $event"
    )
    Register(
      v-bind:class="{hidden: selectedComponent(currentTab, 'Register')}"
      v-bind:soundToPlay="soundToPlay"
      v-on:register:play="soundToPlay = $event"
    )
    SoundUpload(
      v-bind:class="{hidden: selectedComponent(currentTab, 'SoundUpload')}"
    )

    List(
      v-bind:soundToPlay="soundToPlay"
      v-on:update:soundToPlay="soundToPlay = $event"
      v-bind:destination="dest"
    )

</template>

<script>
  import Destination from '@/components/Destination'
  import List from '@/components/List'
  import SoundUpload from '@/components/SoundUpload'
  import Register from '@/components/Register'
  export default {
    name: 'app',
    components: {Destination, List, SoundUpload, Register},
    data: function () {
      return {
        dest: 'test',
        soundToPlay: '',
        currentTab: '',
        tabs: ['Register', 'Destination', 'SoundUpload']
      }
    },
    methods: {
      showDest: function () {
        console.log(this.dest)
        console.log(this.soundToPlay)
      },
      selectedComponent: function (component, target) {
        if (component === target) {
          return ''
        }
        return 'hidden'
      }
    }
  }
</script>

<style lang="scss">
  @import "styles/settings.scss";
  body {
    margin: 0;
    background: #141d26;
  }

  #app {
    font-family: "Avenir", Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    color: #2c3e50;
  }

  .navbar, .navbar-brand {
    background: $secondary-light !important;
    color: white;
    a, .navbar-toggler-icon {
      color: white;
    }
  }

  .custom-toggler{
    border-color: white;
    // background: white;
    color: white;
  }

  main {
    text-align: center;
    margin-top: 40px;
  }

  .hidden {
    display: none;
  }

  header {
    margin: 0;
    height: 56px;
    padding: 0 16px 0 24px;
    background-color: $secondary-light;
    color: #ffffff;
    span {
      display: block;
      position: relative;
      font-size: 20px;
      line-height: 1;
      letter-spacing: 0.02em;
      font-weight: 400;
      box-sizing: border-box;
      padding-top: 16px;
    }
  }
</style>
