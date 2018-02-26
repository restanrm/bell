<template lang="jade">
  .app
    nav.navbar.navbar-toggleable-md
      button.navbar-toggler.custom-toggler(type="button",
        data-toggle="collapse",
        data-target="#navbarSupportedContent",
        aria-controls="navbarSupportedContent",
        aria-expanded="false",
        aria-label="Toggle navigation",
        )
        span.custom-toggler
          i.fa.fa-bars.white
      a.navbar-brand(href="#") Bell
      .collapse.navbar-collapse#navbarSupportedContent
        router-link.nav-item.nav-link.active(to="/") List
        router-link.nav-item.nav-link(to="/upload") Upload
    router-view
</template>

<script>
  
  export default {
    name: 'app'
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
