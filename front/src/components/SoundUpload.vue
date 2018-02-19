<template lang="jade">
  .jumbotron.soundUpload
    .upload-files
      .form-group
        .container
          .row
            .col-6
              label Filename: 
              input(type="text",id="inputField")
            .col-3
              file-upload.btn.btn-primary(
                ref="upload",
                v-model="file",
                :post-action='uploadPath',
                @input-file="inputFile",
                :size=1024 * 1024 * 1,
                :data="{filename:'filename'}",
                :name="uploadFile",
                :multiple=false,
                )
                i.fa.fa-file
                span  Upload your sound
            .col-3
              button.btn.btn-success(
                v-show="!$refs.upload || !$refs.upload.active", 
                @click.prevent="$refs.upload.active = true",
                type="button",
                )
                i.fa.fa-upload
                span  Start uploading
              button.btn.btn-success(
                v-show="$refs.upload && $refs.upload.active", 
                @click.prevent="$refs.upload.active = false",
                type="button",
                )
                i.fa.fa-times
                span  Stop upload
</template>

<script>
  import FileUpload from 'vue-upload-component'
  export default {
    name: 'SoundUpload',
    components: {
      FileUpload
    },
    data () {
      var basepath = ''
      if (process.env.NODE_ENV === 'development') {
        basepath = 'http://localhost:10101'
      };
      return {
        uploadPath: basepath + '/api/v1/sounds',
        filename: ''
      }
    },
    methods: {
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

  button.play-btn {
    margin: 15px 0 !important;
    padding: 15px 0 !important;
    width: 100%;
		background: #063F63;
  }

  .player {
    background: #243447;
  }
</style>
