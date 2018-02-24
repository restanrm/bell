<template lang="jade">
  .app
    header
      span Bell
    List
    //- SoundUpload
</template>

<script>
  import List from './components/List'
  import SoundUpload from './components/SoundUpload'
  
  export default {
    name: 'app',
    components: {
      List,
      SoundUpload
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
  
  main {
    text-align: center;
    margin-top: 40px;
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
