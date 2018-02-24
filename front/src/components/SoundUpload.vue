<template lang="jade">
  .jumbotron.soundUpload
    .upload-files
      .form-group
        .container
          .row
            .col-6
              label Name for your sound: 
              input(type="text",id="inputField",v-model="soundName")
            .col-3
              file-upload.btn.btn-primary(
                ref="upload",
                v-model="files",
                :post-action='uploadPath',
                :size=1024 * 1024 * 1,
                :data="{name:soundName}",
                name="uploadFile",
                :multiple=false,
                @input="updateValue",
                )
                i.fa.fa-file
                span  Choose sound
            .col-3
              button.btn.btn-success(
                v-show="!$refs.upload || !$refs.upload.active && soundName!=''", 
                @click.prevent="$refs.upload.active = true",
                type="button",
                )
                i.fa.fa-upload
                span  Start uploading
              button.btn.btn-danger(
                v-show="$refs.upload && $refs.upload.active", 
                @click.prevent="$refs.upload.active = false",
                type="button",
                )
                i.fa.fa-times
                span  Stop upload
      ul
        li(v-for="(f, index) in files",:key="f.id")
          span(v-if="f.error") {{f.error}}
          span(v-else-if="f.success").alert.alert-success Successfully uploaded content
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
        filename: '',
        soundName: ''
      }
    },
    methods: {
      updateValue (value) {
        let refresh = true
        if (this.files) {
          for (let i = 0; i < this.files.length; i++) {
            if (this.files[i].success === false) {
              refresh = false
            }
          }
          if (refresh === true) {
            window.location.reload()
          }
        }
      }
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss">
  @import "../styles/settings.scss";

  .jumbotron.soundUpload {
    background: $secondary-lighter;
    padding: 2rem;
    margin: 0px;
  }

</style>
