<template>
  <div class="list">
    <h1>{{ msg }}</h1>
  </div>
</template>

<script>
export default {
  name: 'list',
  data () {
    this.$http.get('/api/v1/sounds').then(response => {
      return {
        msg: response.body;
      }
    })
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
h1, h2 {
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
