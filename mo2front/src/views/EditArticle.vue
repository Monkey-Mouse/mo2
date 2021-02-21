<template>
  <div>
    <v-row justify="center">
      <v-col cols="12" lg="7" class="mo2editor">
        <v-skeleton-loader
          class="mb-10 mt-10"
          v-if="loading"
          type="heading"
        ></v-skeleton-loader>
        <v-skeleton-loader
          v-if="loading"
          type="paragraph@9"
        ></v-skeleton-loader>
      </v-col>
    </v-row>
    <editor
      v-show="!loading"
      :uploadImgs="uploadImgs"
      :content="content"
      @loaded="editorLoad"
    />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Editor from "../components/MO2Editor.vue";
import { globaldic, UploadImgToQiniu } from "@/utils";
@Component({
  components: {
    Editor,
  },
})
export default class EditArticle extends Vue {
  uploadImgs = UploadImgToQiniu;
  loading = true;
  editorLoaded = false;
  content = "";
  created() {
    if (globaldic.article) {
      this.content = globaldic.article;
      delete globaldic.article;
    }
  }
  editorLoad() {
    this.editorLoaded = true;
    if (!this.$route.params["id"] && this.content !== "") {
      this.loading = false;
    }
  }
}
</script>
