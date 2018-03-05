<template lang="jade">
  .jumbotron.soundUpload.SoundUpload
    .upload-files
      .form-group.row
        label.col-md-2.col-form-label.vert-spacer(for="inputField") Sound name
        input.col-md-5.vert-spacer(type="text",id="inputField",v-model="soundName")
        file-upload.btn.btn-primary.col-md-4.spacer.soundUploadBtn(
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
      .row(v-show="soundName!=''")
        .col-md-4
        button.btn.btn-success.col-md-4.soundUploadBtn(
          v-if="!$refs.upload || !$refs.upload.active && soundName!=''", 
          @click.prevent="$refs.upload.active = true",
          type="button",
          )
          i.fa.fa-upload
          span  Start uploading
        button.btn.btn-danger.col-md-4.soundUploadBtn(
          v-else="$refs.upload && $refs.upload.active", 
          @click.prevent="$refs.upload.active = false",
          type="button",
          )
          i.fa.fa-times
          span  Stop upload
      .row.vert-spacer
        .col-md-4
        div.col-md-4(v-for="(f, index) in files",:key="f.id")
          span(v-if="f.error") {{f.error}}
          span(v-else-if="f.success").alert.alert-success Successfully uploaded sound
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
    background: $primary-darker;
    padding: 2rem;
    margin: 0px;
  }

  .spacer {
    margin: 1.5em;
  }

  .vert-spacer {
    margin: 1.5em 0;
  }

  .soundUploadBtn {
    margin: 15px;
    padding: 15px;
  }

  .btn-success {
    background: $secondary;
    border-color: $secondary-dark;
    }

  .btn-primary {
    background: $primary;
    border-color: $primary-dark;
    &:hover {
      background: $primary-light;
    }
  }

  label {
    color: white;
  }

</style>
